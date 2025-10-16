package db

type SchemaMigrations struct {
	Version int  `db:"version"`
	Dirty   bool `db:"dirty"`
}
