package router

import (
	fantasycontext "dota2_fantasy/src/fantasyContext"
	"dota2_fantasy/src/util"
	"fmt"
	"log"
	"net/http"
)

type Middleware interface {
	WithBaseMiddleware(handler http.HandlerFunc, middlewares ...func(http.Handler) http.HandlerFunc) http.HandlerFunc
	AssignDBPointer(next http.Handler) http.HandlerFunc
	HandlePanic(next http.Handler) http.HandlerFunc
	TestMW(next http.Handler) http.HandlerFunc
}

func NewMiddleware(dbConn util.DBConnection) Middleware {
	return middleware{
		dbConn: dbConn,
	}
}

type middleware struct {
	dbConn util.DBConnection
}

func (m middleware) WithBaseMiddleware(handler http.HandlerFunc, middlewares ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	middlewares = append(middlewares, m.AssignDBPointer)
	middlewares = append(middlewares, m.HandlePanic)

	for _, next := range middlewares {
		handler = next(handler)
	}
	return handler
}

func (m middleware) TestMW(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Here in test mw")
		next.ServeHTTP(w, r)
	})
}

func (m middleware) AssignDBPointer(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Here in test assign db pointer")

		ctx := fantasycontext.WithDBPool(r.Context(), m.dbConn.GetPool())
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func (m middleware) HandlePanic(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Here in handle panic")

		defer func() {
			if err := recover(); err != nil {
				log.Println("panic occurred:", err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
