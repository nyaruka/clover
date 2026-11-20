package runtime_test

import (
	"context"
	"testing"

	"github.com/nyaruka/clover/v26/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadOnlyConnectionDiscarded(t *testing.T) {
	ctx := context.Background()

	cfg := runtime.NewDefaultConfig()
	cfg.DB = "postgres://clover_test:temba@postgres:5432/clover_test?sslmode=disable"
	rt, err := runtime.NewRuntime(cfg)
	require.NoError(t, err)
	defer rt.DB.Close()

	db := rt.DB
	db.SetMaxOpenConns(1)

	backendPID := func() int {
		var pid int
		require.NoError(t, db.GetContext(ctx, &pid, `SELECT pg_backend_pid()`))
		return pid
	}

	pid1 := backendPID()

	// make the connection behave like one to a primary that has since become a read-only replica
	db.MustExecContext(ctx, `SET SESSION CHARACTERISTICS AS TRANSACTION READ ONLY`)
	_, err = db.ExecContext(ctx, `DELETE FROM urn_mappings`)
	assert.ErrorContains(t, err, "read-only transaction")

	assert.NotEqual(t, pid1, backendPID())
}
