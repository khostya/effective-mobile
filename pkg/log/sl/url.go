package sl

import "log/slog"

func URL(url string) slog.Attr {
	return slog.Attr{
		Key:   "url",
		Value: slog.StringValue(url),
	}
}
