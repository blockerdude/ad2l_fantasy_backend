package repo

type Repos struct {
	AuthnRepo             AuthnRepo
	ConferenceRepo        ConferenceRepo
	LeagueParticipantRepo LeagueParticipantRepo
	LeaugeRepo            LeagueRepo
	PlayerRepo            PlayerRepo
	PlayerScoreRepo       PlayerScoreRepo
	RosterRepo            RosterRepo
	SeasonRepo            SeasonRepo
	TeamRepo              TeamRepo
	TimeframeRepo         TimeframeRepo
}

func SetupRepos() Repos {

	repos := Repos{
		AuthnRepo:             NewAuthnRepo(),
		ConferenceRepo:        NewConferenceRepo(),
		LeagueParticipantRepo: NewLeagueParticipantRepo(),
		LeaugeRepo:            NewLeagueRepo(),
		PlayerRepo:            NewPlayerRepo(),
		PlayerScoreRepo:       NewPlayerScoreRepo(),
		RosterRepo:            NewRosterRepo(),
		TeamRepo:              NewTeamRepo(),
		TimeframeRepo:         NewTimeframeRepo(),
	}

	return repos

}
