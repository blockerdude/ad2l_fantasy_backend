CREATE TABLE conference (
  id SERIAL PRIMARY KEY,
  super_admin BOOLEAN NOT NULL DEFAULT FALSE,
  email VARCHAR(128) NOT NULL,
  display_name VARCHAR(128) NOT NULL,
  last_action TIMESTAMPTZ NOT NULL,
  session_token VARCHAR(128) NOT NULL,
);

CREATE TABLE league_participant (
    league_id INT NOT NULL,
    user_id INT NOT NULL,
    league_admin BOOLEAN NOT NULL DEFAULT FALSE,
    paid BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT fk_leagueparticipant_league_id FOREIGN KEY(league_id) REFERENCES league(id),
    CONSTRAINT fk_leagueparticipant_user_id FOREIGN KEY(user_id) REFERENCES user(id),
    PRIMARY KEY (league_id, user_id) 
);

-- Note: At the moment player_ids will not have a foreign key relationship check
-- it will be up to the business logic to enforce that valid player_ids are stored
-- in the future this could change with a change to create a new table to store picks
CREATE TABLE user_pick (
    league_participant_id INT NOT NULL,
    timeframe_id INT NOT NULL,
    player_ids INTEGER ARRAY, 
    CONSTRAINT fk_userpick_leagueparticipant_id FOREIGN KEY(league_participant_id) REFERENCES league_participant(id),
    CONSTRAINT fk_userpick_timeframe_id FOREIGN KEY(timeframe_id) REFERENCES timeframe(id),
    PRIMARY KEY (league_participant_id, timeframe_id)
);

CREATE TABLE player_score (
    timeframe_id INT NOT NULL,
    player_id INT NOT NULL,
    match_id VARCHAR(20) NOT NULL, -- If the match is not reported this is expcted to be something like '1' or '2'
    substitute_id INT, -- Can be nullable (if no substitute exists)
    total DOUBLE PRECISION, -- This is an override, if no data can be provided this will still give a total score which can be used. If total is set, no matches or real data are valid
    kills INT,
    deaths INT,
    last_hits INT,
    denies INT,
    teamfight_participation DECIMAL(3,2), -- 0.00 to 1.00
    gpm INT,
    tower_kills INT,
    rosh_kills INT,
    obs_placed INT,
    camps_stacked INT,
    runes_taken INT,
    first_blood BOOLEAN,
    stun_seconds DOUBLE PRECISION,
    CONSTRAINT fk_playerscore_timeframe_id FOREIGN KEY(timeframe_id) REFERENCES timeframe(id),
    CONSTRAINT fk_playerscore_player_id FOREIGN KEY(player_id) REFERENCES player(id),
    CONSTRAINT fk_playerscore_substitute_id FOREIGN KEY(substitute_id) REFERENCES player(id),
    PRIMARY KEY (timeframe_id, player_id, match_id)
);

CREATE TABLE user_score (
    league_participant_id INT NOT NULL,
    timeframe_id INT NOT NULL,
    score DOUBLE PRECISION NOT NULL,
);