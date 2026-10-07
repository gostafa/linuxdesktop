// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package dbusconn

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"syscall"
	"time"

	"github.com/cenkalti/backoff/v7"
	"github.com/godbus/dbus/v5"
	"github.com/gostafa/linuxdesktop/internal/port"
	"github.com/gostafa/linuxdesktop/internal/probe"
	"github.com/gostafa/linuxdesktop/internal/schema"
	"github.com/gostafa/singleton"
)

//nolint:revive // Construction explicitly binds owner, policy, budget, and injected opener.
func newConnections(
	base context.Context,
	policy schema.RetryPolicy,
	timeout time.Duration,
	open func(context.Context, port.BusKind) (*dbus.Conn, error),
) *connections {
	owner, cancel := context.WithCancel(base)
	cache := new(connections)

	//nolint:fatcontext // The cache owns the context for the entire lifetime of its private connections.
	cache.base, cache.cancel = owner, cancel
	cache.policy, cache.timeout, cache.open = policy, timeout, open

	return cache
}

// connect shares one initialization while each caller bounds only its own wait.
func connect(ctx context.Context, cache *connections, kind port.BusKind) (*dbus.Conn, error) {
	if int(kind) >= len(cache.providers) {
		return nil, ErrBusKind
	}

	conn, closed := connectionState(cache, kind)
	if closed {
		return nil, ErrClosed
	}

	if conn != nil {
		return conn, nil
	}

	return connectionValues[*dbus.Conn](awaitedConnection(ctx, cache, kind))
}

func connectionState(cache *connections, kind port.BusKind) (*dbus.Conn, bool) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	return cache.conns[kind], cache.closed
}

func awaitedConnection(
	ctx context.Context,
	cache *connections,
	kind port.BusKind,
) (*dbus.Conn, error) {
	provider, err := connectionProvider(cache, kind)
	if err != nil {
		return nil, fmt.Errorf("dbusconn: select provider: %w", err)
	}

	waitCtx, cancel := connectionWait(ctx, cache.base)

	defer cancel()

	conn, err := provider.Get(waitCtx)
	if err != nil {
		return nil, fmt.Errorf("dbusconn: cached connection: %w", err)
	}

	return connectionValues[*dbus.Conn](acceptedConnection(cache, conn))
}

//nolint:ireturn // Each caller's wait also ends when the owning run shuts down.
func connectionWait(ctx, owner context.Context) (waitContext context.Context, cleanup func()) {
	waitCtx, cancel := context.WithCancelCause(probe.Default(ctx, context.Background))
	stop := context.AfterFunc(owner, func() { cancel(context.Cause(owner)) })

	return waitCtx, func() { stop(); cancel(nil) }
}

func acceptedConnection(cache *connections, conn *dbus.Conn) (*dbus.Conn, error) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	if cache.closed {
		return nil, ErrClosed
	}

	return conn, nil
}

//nolint:contextcheck // A singleton initializes under the bus owner's context, independently of callers.
func connectionProvider(
	cache *connections,
	kind port.BusKind,
) (*singleton.Provider[*dbus.Conn], error) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	if cache.closed {
		return nil, ErrClosed
	}

	if cache.providers[kind] == nil {
		cache.providers[kind] = singleton.MustNew(func(context.Context) (*dbus.Conn, error) {
			return initializeConnection(cache, kind)
		}, singleton.WithMaxAttempts(singleAttempt), singleton.WithInitializationTimeout(zero))
	}

	return cache.providers[kind], nil
}

func initializeConnection(cache *connections, kind port.BusKind) (*dbus.Conn, error) {
	ctx, cancel := initializationContext(cache.base, cache.timeout)
	defer cancel()

	conn, err := retryConnection(ctx, cache, kind)
	if err != nil {
		return nil, fmt.Errorf("dbusconn: initialize connection: %w", err)
	}

	return connectionValues[*dbus.Conn](publishConnection(ctx, cache, kind, conn))
}

//nolint:ireturn // Return the standard context with its cleanup function.
func initializationContext(
	ctx context.Context,
	timeout time.Duration,
) (context.Context, context.CancelFunc) {
	if timeout > zero {
		return context.WithTimeout(ctx, timeout)
	}

	return context.WithCancel(ctx)
}

func retryConnection(
	ctx context.Context,
	cache *connections,
	kind port.BusKind,
) (*dbus.Conn, error) {
	conn, err := backoff.Retry(ctx, func() (*dbus.Conn, error) {
		return connectionAttempt(ctx, cache, kind)
	},
		backoff.WithBackOff(connectionPolicy(cache.policy)),
		backoff.WithMaxTries(cache.policy.MaxAttempts),
		backoff.WithMaxElapsedTime(zero),
	)

	return connectionValues[*dbus.Conn](conn, err)
}

func connectionPolicy(settings schema.RetryPolicy) *backoff.ExponentialBackOff {
	policy := backoff.NewExponentialBackOff()

	policy.InitialInterval, policy.MaxInterval = settings.InitialInterval, settings.MaxInterval
	policy.Multiplier, policy.RandomizationFactor = retryMultiplier, retryJitter

	return policy
}

// Only the pending attempt observes the initialization budget. After success,
// its lifetime remains attached to the owning run rather than that budget.
//
//nolint:contextcheck // Pending attempts observe the budget but successful connections retain the owner's lifetime.
func connectionAttempt(
	ctx context.Context,
	cache *connections,
	kind port.BusKind,
) (*dbus.Conn, error) {
	err := context.Cause(ctx)
	if err != nil {
		return connectionValues[*dbus.Conn](nil, permanentConnectionError(err))
	}

	attempt, cancel := context.WithCancel(cache.base)
	stop := context.AfterFunc(ctx, cancel)
	conn, err := cache.open(attempt, kind)

	stop()

	return connectionValues[*dbus.Conn](finishAttempt(ctx, cancel, conn, err))
}

//nolint:revive // The result and cancellation function must be checked against the shared budget together.
func finishAttempt(
	ctx context.Context,
	cancel context.CancelFunc,
	conn *dbus.Conn,
	err error,
) (*dbus.Conn, error) {
	cause := context.Cause(ctx)
	if cause != nil {
		cancel()

		return connectionValues[*dbus.Conn](
			nil,
			permanentConnectionError(errors.Join(cause, err, closeConnection(conn))),
		)
	}

	if err != nil || conn == nil {
		cancel()
	}

	return connectionValues[*dbus.Conn](attemptValue(conn, err))
}

func attemptValue(conn *dbus.Conn, err error) (*dbus.Conn, error) {
	if err != nil {
		return connectionValues[*dbus.Conn](
			nil,
			retryableConnectionError(err, closeConnection(conn)),
		)
	}

	if conn == nil {
		return connectionValues[*dbus.Conn](nil, permanentConnectionError(errNoConnection))
	}

	return conn, nil
}

func retryableConnectionError(err, cleanup error) error {
	failure := errors.Join(err, cleanup)
	if !transientConnectionError(err) {
		return errors.Join(permanentConnectionError(failure))
	}

	return failure
}

func permanentConnectionError(err error) error {
	return fmt.Errorf("dbusconn: stop initialization: %w", backoff.Permanent(err))
}

//nolint:revive // Publication needs the shared budget, owner, bus kind, and opened resource.
func publishConnection(
	ctx context.Context,
	cache *connections,
	kind port.BusKind,
	conn *dbus.Conn,
) (*dbus.Conn, error) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	if cache.closed {
		return nil, errors.Join(ErrClosed, closeConnection(conn))
	}

	err := context.Cause(ctx)
	if err != nil {
		return nil, errors.Join(err, closeConnection(conn))
	}

	cache.conns[kind] = conn

	return conn, nil
}

func transientConnectionError(err error) bool {
	causes := []error{
		syscall.ECONNREFUSED,
		syscall.ECONNRESET,
		syscall.EINTR,
		syscall.EAGAIN,
		io.EOF,
		io.ErrUnexpectedEOF,
	}
	for i := range causes {
		if errors.Is(err, causes[i]) {
			return true
		}
	}

	var network net.Error

	return errors.As(err, &network) && network.Timeout()
}

// closeConnections never waits for an opener; late results are discarded by
// connectionAttempt or publishConnection, which also close their resources.
func closeConnections(cache *connections) error {
	conns := detachConnections(cache)
	defer cache.cancel()

	failures := make([]error, zero, len(conns))
	for i := range conns {
		failures = append(failures, closeConnection(conns[i]))
	}

	return errors.Join(failures...)
}

func detachConnections(cache *connections) [busCount]*dbus.Conn {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	if cache.closed {
		return [busCount]*dbus.Conn{}
	}

	cache.closed = true

	conns := cache.conns

	cache.conns = [busCount]*dbus.Conn{}

	return conns
}

func closeConnection(conn *dbus.Conn) error {
	if conn == nil {
		return nil
	}

	return errors.Join(conn.Close())
}

// connectionValues preserves a result and its already classified error chain.
func connectionValues[Value any](value Value, err error) (Value, error) {
	return value, errors.Join(err)
}
