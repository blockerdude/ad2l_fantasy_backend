package router

import (
	"dota2_fantasy/src/service"
	"dota2_fantasy/src/util"

	"github.com/gorilla/mux"
)

type Routers struct {
	AuthnRouter AuthnRouter
}

func SetupRouters(config util.Config, mw Middleware, services service.Services, baseRouter *mux.Router) Routers {
	routers := Routers{
		AuthnRouter: *NewAuthnRouter(config, mw, services.AuthnService),
	}

	routers.AuthnRouter.SetupRoutes(baseRouter)

	return routers
}
