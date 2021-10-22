package middleware

import "net/http"

type HttpMiddleware func(http.HandlerFunc) http.HandlerFunc

type ChainMiddleware interface {
	BuildChain(handlerFunc http.HandlerFunc) http.HandlerFunc
}

type chainMiddleware struct {
	m []HttpMiddleware
}

func NewChainMiddleware(m ...HttpMiddleware) ChainMiddleware {
	return chainMiddleware{m: m}
}

func (c chainMiddleware) BuildChain(handlerFunc http.HandlerFunc) http.HandlerFunc {
	if len(c.m) == 0 {
		return handlerFunc
	}

	return c.m[0](NewChainMiddleware(c.m[1:cap(c.m)]...).BuildChain(handlerFunc))
}
