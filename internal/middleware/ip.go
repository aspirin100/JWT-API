package middleware

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/ogen-go/ogen/middleware"
)

var (
	clientIPHeaders = []string{"X-Forwarded-For"}
	ctxKey          = ctxIP{}

	ErrEmptyIP   = errors.New("empty ip")
	ErrInvalidIP = errors.New("invalid ip")
)

type ctxIP struct{}

func DetectIP(req middleware.Request, next middleware.Next) (middleware.Response, error) { //nolint:gocritic
	userIP, err := parseClientIP(req.Raw)
	if err != nil {
		return next(req)
	}

	req.SetContext(context.WithValue(req.Context, ctxKey, userIP))

	return next(req)
}

func GetClientIP(ctx context.Context) (net.IP, bool) {
	v, ok := ctx.Value(ctxKey).(net.IP)

	return v, ok
}

// SetClientIP используется в тестах.
func SetClientIP(ctx context.Context, ip net.IP) context.Context {
	return context.WithValue(ctx, ctxKey, ip)
}

func parseClientIP(r *http.Request) (net.IP, error) {
	for _, header := range clientIPHeaders {
		v := r.Header.Get(header)
		if v != "" {
			return net.ParseIP(v), nil
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, errors.Join(ErrInvalidIP, fmt.Errorf("failed to split host from port: %w", err))
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return nil, ErrEmptyIP
	}

	return ip, nil
}
