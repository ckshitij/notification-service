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

	id := getRequestIDFromContext(ctx)
	if id != "" {
		fields = append(fields, String("request_id", id))
	}

	return fields
}

func getRequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if requestID, ok := ctx.Value("req-id").(string); ok {
		return requestID
	}

	return ""
}
