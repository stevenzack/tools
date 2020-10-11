package ctxToolkit

import (
	"context"
	"errors"
)

func NotEmptyString(ctx context.Context, key string) (string, error) {
	v := ctx.Value(key)
	if v == nil {
		return "", errors.New("value is nil")
	}
	str, ok := v.(string)
	if !ok {
		return "", errors.New("value is not string")
	}
	if len(str) == 0 {
		return "", errors.New("value is an empty string")
	}
	return str, nil
}

func Int(ctx context.Context, key string) (int, error) {
	v := ctx.Value(key)
	if v == nil {
		return 0, errors.New(key + " is nil")
	}
	i, ok := v.(int)
	if !ok {
		return 0, errors.New(key + " is not int")
	}
	return i, nil
}

func WithMap(ctx context.Context, m map[string]interface{}) context.Context {
	for k, v := range m {
		ctx = context.WithValue(ctx, k, v)
	}
	return ctx
}
