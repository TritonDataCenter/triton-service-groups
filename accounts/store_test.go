package accounts_test

import (
	"context"
	"os"
	"testing"

	"github.com/joyent/triton-service-groups/accounts"
	"github.com/joyent/triton-service-groups/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindByID(t *testing.T) {
	if os.Getenv("TSG_TEST") == "" {
		t.Skip("Acceptance tests skipped unless env 'TSG_TEST=1' set")
		return
	}

	db, err := testutils.NewTestDB()
	if err != nil {
		t.Error(err)
	}
	db.Clear(t)
	defer db.Clear(t)

	store := accounts.NewStore(db.Conn)
	require.NotNil(t, store)

	account := accounts.New(store)
	require.NotNil(t, account)

	account.AccountName = "baconuser"
	account.TritonUUID = "f5435e8b-70b8-4e4d-8c59-1dbe5d100b5b"

	err = account.Insert(context.Background())
	require.NoError(t, err)
	require.NotZero(t, account.ID)

	found, err := store.FindByID(context.Background(), account.ID)
	require.NoError(t, err)

	assert.Equal(t, account.ID, found.ID)
	assert.Equal(t, account.AccountName, found.AccountName)
	assert.Equal(t, account.TritonUUID, found.TritonUUID)
	assert.Equal(t, account.CreatedAt, found.CreatedAt)
	assert.Equal(t, account.UpdatedAt, found.UpdatedAt)
}

func TestFindByName(t *testing.T) {
	if os.Getenv("TSG_TEST") == "" {
		t.Skip("Acceptance tests skipped unless env 'TSG_TEST=1' set")
		return
	}

	db, err := testutils.NewTestDB()
	if err != nil {
		t.Error(err)
	}
	db.Clear(t)
	defer db.Clear(t)

	store := accounts.NewStore(db.Conn)
	require.NotNil(t, store)

	account := accounts.New(store)
	require.NotNil(t, account)

	account.AccountName = "baconuser"
	account.TritonUUID = "f5435e8b-70b8-4e4d-8c59-1dbe5d100b5b"

	err = account.Insert(context.Background())
	require.NoError(t, err)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.AccountName)

	found, err := store.FindByName(context.Background(), account.AccountName)
	require.NoError(t, err)

	assert.Equal(t, account.ID, found.ID)
	assert.Equal(t, account.AccountName, found.AccountName)
	assert.Equal(t, account.TritonUUID, found.TritonUUID)
	assert.Equal(t, account.CreatedAt, found.CreatedAt)
	assert.Equal(t, account.UpdatedAt, found.UpdatedAt)
}
