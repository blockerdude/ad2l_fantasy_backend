package repo

type Repos struct {
	AuthnRepo      AuthnRepo
	ConferenceRepo ConferenceRepo
}

func SetupRepos() Repos {

	repos := Repos{
		AuthnRepo:      NewAuthnRepo(),
		ConferenceRepo: NewConferenceRepo(),
	}

	return repos

}
