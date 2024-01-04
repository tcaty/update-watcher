package repository

const (
	createVersionsTableQuery = `
		CREATE TABLE IF NOT EXISTS versions (
			id serial PRIMARY KEY,
			target VARCHAR (255) UNIQUE NOT NULL,
			version VARCHAR (255) UNIQUE NOT NULL
		)
	`
)
