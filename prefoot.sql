DROP TABLE IF EXISTS "coach_careers" CASCADE;
DROP TABLE IF EXISTS "player_statistics" CASCADE;
DROP TABLE IF EXISTS "players" CASCADE;
DROP TABLE IF EXISTS "player_statistics_season" CASCADE;
DROP TABLE IF EXISTS "formations" CASCADE;
DROP TABLE IF EXISTS "events" CASCADE;
DROP TABLE IF EXISTS "venues" CASCADE;
DROP TABLE IF EXISTS "teams" CASCADE;
DROP TABLE IF EXISTS "fixtures" CASCADE;
DROP TABLE IF EXISTS "leagues" CASCADE;
DROP TABLE IF EXISTS "coaches" CASCADE;

DROP TABLE IF EXISTS "team_statistics" CASCADE;
DROP TABLE IF EXISTS "results" CASCADE;
DROP TABLE IF EXISTS "seasons" CASCADE;
DROP TABLE IF EXISTS "rounds" CASCADE;

CREATE TABLE "leagues" (
  "id" integer PRIMARY KEY,
  "name" varchar,
  "country" varchar
);

CREATE TABLE "teams" (
  "id" integer PRIMARY KEY,
  "name" varchar,
  "country" varchar,
  "code" varchar
);

CREATE TABLE "seasons" (
  "league" int,
  "season" int,
  "team" int,
  PRIMARY KEY ("league", "season", "team")
);

CREATE TABLE "results"(
  "team" integer,
  "league" integer,
  "fixture" integer,
  "round" integer,
  "season" integer,
  "points" integer,
  "goals_for" integer,
  "goals_against" integer,
  "modus" integer,
  "elapsed" integer,
  PRIMARY KEY ("team", "round" ,"season")
);

CREATE TABLE "fixtures" (
  "id" integer PRIMARY KEY,
  "league" integer,
  "round" integer,
  "referee" varchar,
  "timezone" varchar,
  "timestamp" integer,
  "venue" integer,
  "season" integer default 0,
  "home_team" integer,
  "away_team" integer,
  "home_goals" integer,
  "away_goals" integer,
  "home_goals_half" integer,
  "away_goals_half" integer,
  CONSTRAINT fk_fixtures_leagues FOREIGN KEY ("league") REFERENCES "leagues" ("id") DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE "rounds" (
  "start" integer,
  "end" integer,
  "round" integer,
  "season" integer,
  "league" integer,
  PRIMARY KEY ("round", "season", "league")
);

CREATE INDEX rounds_timestamp ON rounds("timestamp");

CREATE TABLE "venues" (
  "id" integer PRIMARY KEY,
  "name" varchar,
  "city" varchar
);

CREATE TABLE "events" (
  "id" integer PRIMARY KEY,
  "fixture" integer,
  "player" integer,
  "assist" integer,
  "minute" integer,
  "team" integer,
  "type" varchar
);

CREATE TABLE "formations" (
  "fixture" integer,
  "team" integer,
  "formation" varchar,
  "player1" integer,
  "player2" integer,
  "player3" integer,
  "player4" integer,
  "player5" integer,
  "player6" integer,
  "player7" integer,
  "player8" integer,
  "player9" integer,
  "player10" integer,
  "player11" integer,
  "sub1" integer,
  "sub2" integer,
  "sub3" integer,
  "sub4" integer,
  "sub5" integer,
  "coach" integer,
  PRIMARY KEY("fixture", "team")
);

CREATE TABLE "players" (
  "id" integer,
  "team" integer,
  "season" integer,
  "firstname" varchar,
  "lastname" varchar,
  "birthplace" varchar,
  "birthcountry" varchar,
  "birthdate" varchar
  PRIMARY KEY ("id", "team", "season")
);

CREATE TABLE "player_statistics" (
  "player" integer,
  "fixture" integer,
  "team" integer,
  "season" integer,
  "league" integer,
  "minutes" integer,
  "position" varchar,
  "rating" float,
  "captain" boolean,
  "substitute" boolean,
  "shots_total" integer,
  "shots_on" integer,
  "goals_scored" integer,
  "goals_assisted" integer,
  "passes_total" integer,
  "passes_key" integer,
  "accuracy" integer,
  "tackles" integer,
  "block" integer,
  "interceptions" integer,
  "duels_total" integer,
  "duels_won" integer,
  "dribbles_total" integer,
  "dribbles_won" integer,
  "yellow" integer,
  "red" integer,
  "penalty_won" integer,
  "penalty_committed" integer,
  "penalty_scored" integer,
  "penalty_missed" integer,
  "penalty_saved" integer,
  "saves" integer,
  PRIMARY KEY("player", "fixture")
);

CREATE TABLE "player_statistics_season" (
  "player" integer,
  "season" integer,
	"team" integer,
  "minutes" integer,
  "position" varchar,
  "rating" float,
  "captain" boolean,
  "games" integer,
  "lineups" integer,
  "shots_total" integer,
  "shots_on" integer,
  "goals_scored" integer,
  "goals_assisted" integer,
  "passes_total" integer,
  "passes_key" integer,
  "accuracy" integer,
  "tackles" integer,
  "block" integer,
  "interceptions" integer,
  "duels_total" integer,
  "duels_won" integer,
  "dribbles_total" integer,
  "dribbles_won" integer,
  "yellow" integer,
  "red" integer,
  "penalty_won" integer,
  "penalty_committed" integer,
  "penalty_scored" integer,
  "penalty_missed" integer,
  "penalty_saved" integer,
  "saves" integer,
  PRIMARY KEY("player", "season", "team")
);


CREATE TABLE "team_statistics" (
  "team" integer,
  "fixture" integer,
  "shots_total" integer,
  "shots_on" integer,
  "shots_off" integer,
  "shots_blocked" integer,
  "shots_box" integer,
  "shots_outside" integer,
  "offsides" integer,
  "fouls" integer,
  "corners" integer,
  "possession" integer,
  "yellow" integer,
  "red" integer,
  "gk_saves" integer,
  "passes_total" integer,
  "passes_accurate" integer,
  "passes_percent" integer,
  "expected_goals" float,
  PRIMARY KEY ("team", "fixture")
);

CREATE TABLE "coaches" (
  "id" integer PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "coach_careers" (
  "coach" integer ,
  "team" integer,
  "start" timestamp,
  "end" timestamp NULL,
  PRIMARY KEY ("coach", "team", "start")
);

CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

ALTER TABLE "fixtures" ADD FOREIGN KEY ("home_team") REFERENCES "teams" ("id") DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "fixtures" ADD FOREIGN KEY ("away_team") REFERENCES "teams" ("id") DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "events" ADD FOREIGN KEY ("fixture") REFERENCES "fixtures" ("id") DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "formations" ADD FOREIGN KEY ("fixture") REFERENCES "fixtures" ("id") DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "player_statistics" ADD FOREIGN KEY ("player") REFERENCES "players" ("id") DEFERRABLE INITIALLY DEFERRED;