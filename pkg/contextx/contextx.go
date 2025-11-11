package contextx

import "context"

type ctxKey string

const (
	UserIDKey   ctxKey = "user_id"
	UsernameKey ctxKey = "username"
)

// GetUserID — безопасно достаёт user_id
func GetUserID(ctx context.Context) string {
	if val, ok := ctx.Value(UserIDKey).(string); ok {
		return val
	}
	return ""
}

// GetUsername — безопасно достаёт username
func GetUsername(ctx context.Context) string {
	if val, ok := ctx.Value(UsernameKey).(string); ok {
		return val
	}
	return ""
}
