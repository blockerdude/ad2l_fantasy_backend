CREATE TABLE conference (
  id SERIAL PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  description VARCHAR NOT NULL
);

CREATE TABLE season (
  id SERIAL PRIMARY KEY,
  conference_id INT NOT NULL,
  name VARCHAR(128) NOT NULL,
  start_date TIMESTAMPTZ NOT NULL,
  end_date TIMESTAMPTZ NOT NULL,
  CONSTRAINT fk_season_conference_id FOREIGN KEY(conference_id) REFERENCES conference(id)
);

CREATE TABLE timeframe (
  id SERIAL PRIMARY KEY,
  season_id INT NOT NULL,
  open BOOLEAN NOT NULL,
  close_date TIMESTAMPTZ NOT NULL,
  name VARCHAR(128) NOT NULL,
  CONSTRAINT fk_timeframe_season_id FOREIGN KEY(season_id) REFERENCES season(id)
);

CREATE TABLE team (
  id SERIAL PRIMARY KEY,
  conference_id INT NOT NULL,
  name VARCHAR(128) NOT NULL,
  flag_url VARCHAR(128),
  CONSTRAINT fk_team_conference_id FOREIGN KEY(conference_id) REFERENCES conference(id)
);

CREATE TABLE roster (
  id SERIAL PRIMARY KEY,
  season_id INT NOT NULL,
  team_id INT NOT NULL,
  CONSTRAINT fk_roster_season_id FOREIGN KEY(season_id) REFERENCES season(id),
  CONSTRAINT fk_roster_team_id FOREIGN KEY(team_id) REFERENCES team(id)
);

CREATE TABLE player (
  id SERIAL PRIMARY KEY,
  steam_id INT NOT NULL UNIQUE,
  name VARCHAR(128) NOT NULL
);

CREATE TABLE roster_player (
  player_id INT NOT NULL,
  roster_id INT NOT NULL,
  CONSTRAINT fk_rosterplayer_player_id FOREIGN KEY(player_id) REFERENCES player(id),
  CONSTRAINT fk_rosterplayer_roster_id FOREIGN KEY(roster_id) REFERENCES roster(id),
  PRIMARY KEY (player_id, roster_id)
);

CREATE TABLE league (
  id SERIAL PRIMARY KEY,
  season_id INT NOT NULL,
  name VARCHAR(128) NOT NULL,
  description VARCHAR NOT NULL,
  CONSTRAINT fk_league_season_id FOREIGN KEY(season_id) REFERENCES season(id)
);