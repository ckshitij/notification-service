package logger

import "context"

type ctxKey string

const (
	RequestIDKey ctxKey = "request_id"
)

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func extractContextFields(ctx context.Context) []Field {
	var fields []Field

	if v := ctx.Value(RequestIDKey); v != nil {
		fields = append(fields, String("request_id", v.(string)))
	}

	return fields
}
