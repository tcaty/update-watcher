CREATE TABLE IF NOT EXISTS versionRecords(
    id serial PRIMARY KEY,
    target VARCHAR (255) UNIQUE NOT NULL,
    version VARCHAR (255) NOT NULL
);