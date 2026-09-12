// Package logind reads the systemd-logind seat session.
//
// The fast path never touches D-Bus. logind mirrors every session's state into
// a world-readable KEY=VALUE file at /run/systemd/sessions/<id>, which is
// exactly what libsystemd's sd_session_* family reads, so one open and one
// read answers the whole question.
//
// Finding the session id follows the same ladder libsystemd uses: $XDG_SESSION_ID
// first, then the session-<id>.scope component of /proc/self/cgroup, then a
// scan of /run/systemd/sessions for a session owned by this uid. Only if all
// three fail does the probe fall back to D-Bus, where logind's magic
// /org/freedesktop/login1/session/auto object resolves to the caller's own
// session in a single GetAll round trip.
//
// SessionInfo.Name carries the session's PAM service ("gdm-password", "sshd",
// "login"). logind's own Session.Name property is the user name, which would
// make the field a duplicate of SessionInfo.User; the service is the more
// useful of the two and is what this library reports.
package logind
