package migrations_test

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nyaruka/clover/v26/migrations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTablesHavePrimaryKeys(t *testing.T) {
	db, err := sqlx.Open("postgres", "postgres://clover_test:temba@postgres:5432/clover_test?sslmode=disable")
	require.NoError(t, err)
	defer db.Close()

	db.Exec("drop table urn_mappings cascade;")
	db.Exec("drop table interchanges cascade;")
	db.Exec("drop table channels cascade;")
	db.Exec("drop table migrations;")
	require.NoError(t, migrations.Migrate(context.Background(), db))

	// logical replication can't replicate UPDATE or DELETE on a table without a primary key
	var missing []string
	err = db.Select(&missing, `
		SELECT c.relname FROM pg_class c
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname = current_schema() AND c.relkind = 'r'
		AND NOT EXISTS (SELECT 1 FROM pg_constraint p WHERE p.conrelid = c.oid AND p.contype = 'p')
		ORDER BY c.relname`)
	require.NoError(t, err)
	assert.Empty(t, missing)
}
