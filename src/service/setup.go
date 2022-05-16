package service

import (
	"dota2_fantasy/src/repo"
	"dota2_fantasy/src/util"
)

type Services struct {
	AuthnService AuthnService
}

func SetupServices(config util.Config, repos repo.Repos) Services {
	services := Services{
		AuthnService: NewAuthnService(config, repos.AuthnRepo),
	}

	return services
}
