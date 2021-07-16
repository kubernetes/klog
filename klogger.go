package klog

// clone from github.com/go-logr/logr v0.4.0
type Logger interface {
	// Enabled tests whether this Logger is enabled.  For example, commandline
	// flags might be used to set the logging verbosity and disable some info
	// logs.
	Enabled() bool

	// Info logs a non-error message with the given key/value pairs as context.
	//
	// The msg argument should be used to add some constant description to
	// the log line.  The key/value pairs can then be used to add additional
	// variable information.  The key/value pairs should alternate string
	// keys and arbitrary values.
	Info(msg string, keysAndValues ...interface{})

	// Error logs an error, with the given message and key/value pairs as context.
	// It functions similarly to calling Info with the "error" named value, but may
	// have unique behavior, and should be preferred for logging errors (see the
	// package documentations for more information).
	//
	// The msg field should be used to add context to any underlying error,
	// while the err field should be used to attach the actual error that
	// triggered this log line, if present.
	Error(err error, msg string, keysAndValues ...interface{})

	// V returns an Logger value for a specific verbosity level, relative to
	// this Logger.  In other words, V values are additive.  V higher verbosity
	// level means a log message is less important.  It's illegal to pass a log
	// level less than zero.
	V(level int) Logger

	// WithValues adds some key-value pairs of context to a logger.
	// See Info for documentation on how key/value pairs work.
	WithValues(keysAndValues ...interface{}) Logger

	// WithName adds a new element to the logger's name.
	// Successive calls with WithName continue to append
	// suffixes to the logger's name.  It's strongly recommended
	// that name segments contain only letters, digits, and hyphens
	// (see the package documentation for more information).
	WithName(name string) Logger
}

type CallDepthLogger interface {
	Logger

	// WithCallDepth returns a Logger that will offset the call stack by the
	// specified number of frames when logging call site information.  If depth
	// is 0 the attribution should be to the direct caller of this method.  If
	// depth is 1 the attribution should skip 1 call frame, and so on.
	// Successive calls to this are additive.
	WithCallDepth(depth int) Logger
}

// WithCallDepth returns a Logger that will offset the call stack by the
// specified number of frames when logging call site information, if possible.
// This is useful for users who have helper functions between the "real" call
// site and the actual calls to Logger methods.  If depth is 0 the attribution
// should be to the direct caller of this function.  If depth is 1 the
// attribution should skip 1 call frame, and so on.  Successive calls to this
// are additive.
//
// If the underlying log implementation supports the CallDepthLogger interface,
// the WithCallDepth method will be called and the result returned.  If the
// implementation does not support CallDepthLogger, the original Logger will be
// returned.
//
// Callers which care about whether this was supported or not should test for
// CallDepthLogger support themselves.
func WithCallDepth(logger Logger, depth int) Logger {
	if decorator, ok := logger.(CallDepthLogger); ok {
		return decorator.WithCallDepth(depth)
	}
	return logger
}
