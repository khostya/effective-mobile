package sl

import "log/slog"

func Code(code uint64) slog.Attr {
	return slog.Attr{
		Key:   "code",
		Value: slog.Uint64Value(code),
	}
}
