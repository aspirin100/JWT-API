package middleware

import (
	"log"
	"runtime/debug"

	"github.com/ogen-go/ogen/middleware"
)

func Recover(req middleware.Request, next middleware.Next) (middleware.Response, error) { //nolint:gocritic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic recovered: %s\n %s\n %s\n", req.OperationName, r, debug.Stack())
		}
	}()

	return next(req)
}
