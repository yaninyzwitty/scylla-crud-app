package middleware

import "net/http"

type Middleware func(next http.Handler) http.Handler

func ChainMiddlewares(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
