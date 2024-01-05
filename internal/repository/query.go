package repository

const (
	createVersionsTableQuery = `
		CREATE TABLE IF NOT EXISTS versionRecords (
			id serial PRIMARY KEY,
			target VARCHAR (255) UNIQUE NOT NULL,
			version VARCHAR (255) UNIQUE NOT NULL
		)
	`
	selectIdByTargetQuery = `
		SELECT id FROM versionRecords WHERE target = $1
	`
	selectVersionByTargetQuery = `
		SELECT version FROM versionRecords WHERE target = $1
	`
	insertVersionRecords = `
		INSERT INTO versionRecords (target, version) VALUES ($1, $2)
	`
	updateVersionRecordVersionByTarget = `
		UPDATE versionRecords
		SET version = $1
		WHERE target = $2
	`
)
