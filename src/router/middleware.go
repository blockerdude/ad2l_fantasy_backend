package router

import (
	fantasycontext "dota2_fantasy/src/fantasyContext"
	"dota2_fantasy/src/service"
	"dota2_fantasy/src/util"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Middleware interface {
	WithBaseMiddleware(handler http.HandlerFunc, middlewares ...func(http.Handler) http.HandlerFunc) http.HandlerFunc
	AssignDBPointer(next http.Handler) http.HandlerFunc
	HandlePanic(next http.Handler) http.HandlerFunc
	RequireLogin(next http.Handler) http.HandlerFunc
}

func NewMiddleware(dbConn util.DBConnection, authnSvc service.AuthnService) Middleware {
	return middleware{
		dbConn:   dbConn,
		authnSvc: authnSvc,
	}
}

type middleware struct {
	dbConn   util.DBConnection
	authnSvc service.AuthnService
}

func (m middleware) WithBaseMiddleware(handler http.HandlerFunc, middlewares ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	middlewares = append(middlewares, m.AssignDBPointer)
	middlewares = append(middlewares, m.HandlePanic)

	for _, next := range middlewares {
		handler = next(handler)
	}
	return handler
}

func (m middleware) AssignDBPointer(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Here in test assign db pointer")

		ctx := fantasycontext.WithDBPool(r.Context(), m.dbConn.GetPool())
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func (m middleware) RequireLogin(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		sessionCookie, err := req.Cookie(SESSION_COOKIE)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := req.Context()
		authn, err := m.authnSvc.GetAuthnByToken(ctx, sessionCookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Has to have been active/logged in within 30 minutes
		timeDif := time.Now().Sub(authn.LastAction)
		if timeDif > (30 * time.Minute) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx = fantasycontext.WithAuthn(ctx, authn)
		req = req.WithContext(ctx)

		if err := m.authnSvc.UpdateLastActionTime(req.Context()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

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
