package router

import (
	"dota2_fantasy/src/service"
	"dota2_fantasy/src/util"

	"github.com/gorilla/mux"
)

type Routers struct {
	AuthnRouter AuthnRouter
}

func SetupRouters(config util.Config, services service.Services, baseRouter *mux.Router) Routers {
	routers := Routers{
		AuthnRouter: *NewAuthnRouter(config, services.AuthnService),
	}

	routers.AuthnRouter.SetupRoutes(baseRouter)

	return routers
}
