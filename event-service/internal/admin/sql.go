package admin

func sqlDropTableEvent() string {
	return "DROP TABLE event"
}

func sqlTruncateSchemaMigrations() string {
	return "TRUNCATE TABLE schema_migrations;"
}
