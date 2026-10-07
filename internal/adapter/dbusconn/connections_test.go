// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package dbusconn

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/gostafa/linuxdesktop/internal/port"
	"github.com/gostafa/linuxdesktop/internal/schema"
)

type countedTransport struct {
	closes atomic.Int32
	closed chan struct{}
	once   sync.Once
}

func (*countedTransport) Read([]byte) (int, error)       { return 0, io.EOF }
func (*countedTransport) Write(data []byte) (int, error) { return len(data), nil }
func (transport *countedTransport) Close() error {
	transport.closes.Add(1)
	transport.once.Do(func() { close(transport.closed) })
	return nil
}

func fixtureConnection(t testing.TB) (*dbus.Conn, *countedTransport) {
	t.Helper()
	transport := &countedTransport{closed: make(chan struct{})}
	conn, err := dbus.NewConn(transport)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	return conn, transport
}

func fastPolicy(attempts uint) schema.RetryPolicy {
	return schema.RetryPolicy{
		MaxAttempts:     attempts,
		InitialInterval: time.Microsecond,
		MaxInterval:     time.Microsecond,
	}
}

func TestInitializationSharingAndCallerCancellation(t *testing.T) {
	t.Parallel()
	conn, _ := fixtureConnection(t)
	started, release := make(chan struct{}), make(chan struct{})
	var calls atomic.Int32
	cache := newConnections(
		t.Context(),
		fastPolicy(3),
		time.Second,
		func(context.Context, port.BusKind) (*dbus.Conn, error) {
			calls.Add(1)
			close(started)
			<-release
			return conn, nil
		},
	)
	t.Cleanup(func() { _ = closeConnections(cache) })
	ctx, cancel := context.WithCancelCause(t.Context())
	abandoned := make(chan error, 1)
	go func() { _, err := connect(ctx, cache, port.SessionBus); abandoned <- err }()
	<-started
	cause := errors.New("caller stopped")
	cancel(cause)
	if err := <-abandoned; !errors.Is(err, cause) {
		t.Fatal(err)
	}
	results := make(chan error, 32)
	for range cap(results) {
		go func() {
			got, err := connect(t.Context(), cache, port.SessionBus)
			if got != conn && err == nil {
				err = errors.New("different connection")
			}
			results <- err
		}()
	}
	close(release)
	for range cap(results) {
		if err := <-results; err != nil {
			t.Fatal(err)
		}
	}
	if calls.Load() != 1 {
		t.Fatal(calls.Load())
	}
}

func TestConnectionRetriesAndCachedFailures(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name            string
		failure         error
		limit           uint
		failCount, want int32
	}{
		{"recovery", syscall.ECONNREFUSED, 3, 2, 3},
		{"exhaustion", syscall.ECONNRESET, 3, 4, 3},
		{"single attempt", syscall.ECONNREFUSED, 1, 2, 1},
		{"missing", os.ErrNotExist, 3, 2, 1},
		{"permission", os.ErrPermission, 3, 2, 1},
		{"authentication", errors.New("authentication denied"), 3, 2, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			conn, _ := fixtureConnection(t)
			var calls atomic.Int32
			var lifetime context.Context
			cache := newConnections(
				t.Context(),
				fastPolicy(tc.limit),
				time.Second,
				func(ctx context.Context, _ port.BusKind) (*dbus.Conn, error) {
					if calls.Add(1) <= tc.failCount {
						return nil, fmt.Errorf("open: %w", tc.failure)
					}
					lifetime = ctx
					return conn, nil
				},
			)
			t.Cleanup(func() { _ = closeConnections(cache) })
			for range 2 {
				got, err := connect(t.Context(), cache, port.SessionBus)
				if tc.want > tc.failCount {
					if err != nil || got != conn || lifetime.Err() != nil {
						t.Fatal(got, err, lifetime)
					}
				} else if !errors.Is(err, tc.failure) || got != nil {
					t.Fatal(got, err)
				}
			}
			if calls.Load() != tc.want {
				t.Fatal(calls.Load(), tc.want)
			}
		})
	}
}

func TestInitializationDeadline(t *testing.T) {
	t.Parallel()
	cache := newConnections(
		t.Context(),
		fastPolicy(3),
		10*time.Millisecond,
		func(ctx context.Context, _ port.BusKind) (*dbus.Conn, error) {
			<-ctx.Done()
			return nil, context.Cause(ctx)
		},
	)
	t.Cleanup(func() { _ = closeConnections(cache) })
	if conn, err := connect(
		t.Context(),
		cache,
		port.SessionBus,
	); conn != nil ||
		!errors.Is(err, context.DeadlineExceeded) {
		t.Fatal(conn, err)
	}
}

func TestCloseDuringInitialization(t *testing.T) {
	t.Parallel()
	conn, transport := fixtureConnection(t)
	started, release := make(chan struct{}), make(chan struct{})
	cache := newConnections(
		t.Context(),
		fastPolicy(3),
		0,
		func(context.Context, port.BusKind) (*dbus.Conn, error) {
			close(started)
			<-release // Simulate an opener that cannot be interrupted.
			return conn, nil
		},
	)
	finished := make(chan error, 1)
	go func() { _, err := connect(t.Context(), cache, port.SessionBus); finished <- err }()
	<-started
	if err := closeConnections(cache); err != nil {
		t.Fatal(err)
	}
	if _, err := connect(t.Context(), cache, port.SystemBus); !errors.Is(err, ErrClosed) {
		t.Fatal(err)
	}
	if err := <-finished; !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	close(release)
	<-transport.closed
	if transport.closes.Load() != 1 {
		t.Fatal(transport.closes.Load())
	}
	if err := closeConnections(cache); err != nil {
		t.Fatal(err)
	}
}

func TestFailedAttemptAndPublicationCleanup(t *testing.T) {
	t.Parallel()
	conn, transport := fixtureConnection(t)
	cache := newConnections(
		t.Context(),
		fastPolicy(1),
		0,
		func(context.Context, port.BusKind) (*dbus.Conn, error) {
			return conn, syscall.ECONNREFUSED
		},
	)
	t.Cleanup(func() { _ = closeConnections(cache) })
	if _, err := connect(
		t.Context(),
		cache,
		port.SessionBus,
	); !errors.Is(
		err,
		syscall.ECONNREFUSED,
	) {
		t.Fatal(err)
	}
	if transport.closes.Load() != 1 {
		t.Fatal("failed attempt leaked")
	}
	for _, closed := range []bool{false, true} {
		late, tracked := fixtureConnection(t)
		owner := newConnections(t.Context(), fastPolicy(1), 0, nil)
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		if closed {
			_ = closeConnections(owner)
		}
		if got, err := publishConnection(
			ctx,
			owner,
			port.SessionBus,
			late,
		); got != nil ||
			err == nil {
			t.Fatal(got, err)
		}
		if tracked.closes.Load() != 1 {
			t.Fatal("late connection leaked")
		}
		_ = closeConnections(owner)
	}
}

func TestOwnerCancellationAndNilConnection(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	var calls atomic.Int32
	open := func(context.Context, port.BusKind) (*dbus.Conn, error) { calls.Add(1); return nil, nil }
	cache := newConnections(ctx, fastPolicy(3), 0, open)
	t.Cleanup(func() { _ = closeConnections(cache) })
	if _, err := connect(
		t.Context(),
		cache,
		port.SessionBus,
	); !errors.Is(err, context.Canceled) ||
		calls.Load() != 0 {
		t.Fatal(err, calls.Load())
	}
	cache = newConnections(t.Context(), fastPolicy(3), 0, open)
	t.Cleanup(func() { _ = closeConnections(cache) })
	if _, err := connect(t.Context(), cache, port.SessionBus); err == nil || calls.Load() != 1 {
		t.Fatal(err, calls.Load())
	}
}

func TestTransientConnectionErrors(t *testing.T) {
	t.Parallel()
	for _, err := range []error{syscall.ECONNREFUSED, syscall.ECONNRESET, syscall.EINTR, syscall.EAGAIN, io.EOF, io.ErrUnexpectedEOF, os.ErrDeadlineExceeded} {
		if !transientConnectionError(fmt.Errorf("wrapped: %w", err)) {
			t.Fatal(err)
		}
	}
	for _, err := range []error{os.ErrNotExist, os.ErrPermission, errors.New("malformed address")} {
		if transientConnectionError(err) {
			t.Fatal(err)
		}
	}
}

func TestAcquisitionRejectedAfterShutdown(t *testing.T) {
	t.Parallel()
	conn, _ := fixtureConnection(t)
	cache := newConnections(t.Context(), fastPolicy(1), 0, nil)
	_ = closeConnections(cache)
	if provider, err := connectionProvider(
		cache,
		port.SessionBus,
	); provider != nil ||
		!errors.Is(err, ErrClosed) {
		t.Fatal(provider, err)
	}
	if got, err := acceptedConnection(cache, conn); got != nil || !errors.Is(err, ErrClosed) {
		t.Fatal(got, err)
	}
	if got, err := awaitedConnection(
		t.Context(),
		cache,
		port.SessionBus,
	); got != nil ||
		!errors.Is(err, ErrClosed) {
		t.Fatal(got, err)
	}
}

func TestPermanentFailureWithTransientCleanup(t *testing.T) {
	t.Parallel()
	conn, err := dbus.NewConn(transport{closed: io.EOF})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	var calls atomic.Int32
	cache := newConnections(
		t.Context(),
		fastPolicy(3),
		0,
		func(context.Context, port.BusKind) (*dbus.Conn, error) {
			calls.Add(1)
			return conn, os.ErrPermission
		},
	)
	t.Cleanup(func() { _ = closeConnections(cache) })
	_, err = connect(t.Context(), cache, port.SessionBus)
	if !errors.Is(err, os.ErrPermission) || !errors.Is(err, io.EOF) || calls.Load() != 1 {
		t.Fatal("cleanup error changed retry classification", err, calls.Load())
	}
}

func BenchmarkCachedConnection(b *testing.B) {
	conn, _ := fixtureConnection(b)
	var calls atomic.Int32
	cache := newConnections(
		b.Context(),
		fastPolicy(1),
		0,
		func(context.Context, port.BusKind) (*dbus.Conn, error) {
			calls.Add(1)
			return conn, nil
		},
	)
	b.Cleanup(func() { _ = closeConnections(cache) })
	if _, err := connect(b.Context(), cache, port.SessionBus); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if _, err := connect(b.Context(), cache, port.SessionBus); err != nil {
			b.Fatal(err)
		}
	}
	b.ReportMetric(float64(calls.Load()), "opens")
}
