package router

import (
	fantasycontext "dota2_fantasy/src/fantasyContext"
	"dota2_fantasy/src/util"
	"log"
	"net/http"
)

type Middleware interface {
	// WithMiddleware()
	AssignDBPointer(next http.Handler) http.Handler
	HandlePanic(next http.Handler) http.Handler
}

func NewAuthnService(dbConn util.DBConnection) Middleware {
	return middleware{
		dbConn: dbConn,
	}
}

type middleware struct {
	dbConn util.DBConnection
}

func (m middleware) AssignDBPointer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := fantasycontext.WithDBPool(r.Context(), m.dbConn.GetPool())
		r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (m middleware) HandlePanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic occurred:", err)
			}
		}()
	})
}
