package postgres

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestGetAllVotes(t *testing.T) {
	dsn := "postgres://postgres:1243maksim@localhost:5432/quality_control?sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	require.NoError(t, err)
	defer db.Close()

	repo := NewReportRepository(db)
	result, err := repo.GetAllVotes(context.Background())

	require.NoError(t, err)
	require.NotNil(t, result)

	for _, v := range *result {
		t.Logf("%+v", v)
	}
}