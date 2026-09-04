CREATE SCHEMA tgbot;

CREATE TABLE IF NOT EXISTS tgbot.users
(
    id          SERIAL              PRIMARY KEY,
    username    TEXT    NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS tgbot.pages
(
    id          SERIAL              PRIMARY KEY,
    url TEXT             NOT NULL,
    user_id     INTEGER  NOT NULL   REFERENCES tgbot.users(id)
);