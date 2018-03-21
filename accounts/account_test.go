package accounts_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joyent/triton-service-groups/accounts"
	"github.com/joyent/triton-service-groups/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	if os.Getenv("TSG_TEST") == "" {
		t.Skip("Acceptance tests skipped unless env 'TRITON_TEST=1' set")
		return
	}

	db, err := testutils.NewTestDB()
	if err != nil {
		t.Error(err)
	}
	db.Clear(t)

	store := accounts.NewStore(db.Conn)
	require.NotNil(t, store)

	account := accounts.New(store)
	require.NotNil(t, account)
}

func TestInsert(t *testing.T) {
	if os.Getenv("TSG_TEST") == "" {
		t.Skip("Acceptance tests skipped unless env 'TRITON_TEST=1' set")
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

	assert.NotZero(t, account.ID)
	assert.Equal(t, account.AccountName, "baconuser")
	assert.Equal(t, account.TritonUUID, "f5435e8b-70b8-4e4d-8c59-1dbe5d100b5b")
	assert.NotZero(t, account.CreatedAt)
	assert.NotZero(t, account.UpdatedAt)
	assert.Equal(t, account.CreatedAt, account.UpdatedAt)
}

func TestSave(t *testing.T) {
	if os.Getenv("TSG_TEST") == "" {
		t.Skip("Acceptance tests skipped unless env 'TRITON_TEST=1' set")
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

	account.AccountName = "demouser"
	account.TritonUUID = "f5435e8b-70b8-4e4d-8c59-1dbe5d100b5b"

	err = account.Insert(context.Background())
	require.NoError(t, err)

	assert.Equal(t, account.CreatedAt, account.UpdatedAt)

	account.AccountName = "hackerman"

	time.Sleep(2 * time.Second)

	err = account.Save(context.Background())
	require.NoError(t, err)

	assert.NotZero(t, account.ID)
	assert.Equal(t, account.AccountName, "hackerman")
	assert.Equal(t, account.TritonUUID, "f5435e8b-70b8-4e4d-8c59-1dbe5d100b5b")
	assert.NotZero(t, account.CreatedAt)
	assert.NotZero(t, account.UpdatedAt)
	assert.NotEqual(t, account.CreatedAt, account.UpdatedAt)
}

func TestExists(t *testing.T) {
	if os.Getenv("TSG_TEST") == "" {
		t.Skip("Acceptance tests skipped unless env 'TRITON_TEST=1' set")
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

	created := accounts.New(store)
	require.NotNil(t, created)
	created.AccountName = "firstcreate"
	created.Insert(context.Background())

	newAccount := accounts.New(store)
	require.NotNil(t, newAccount)

	newAccount.ID = created.ID

	// true if account shares ID with a previously created account
	{
		exists, err := newAccount.Exists(context.Background())
		require.NoError(t, err)
		assert.True(t, exists)
	}

	account := accounts.New(store)
	require.NotNil(t, account)

	account.AccountName = "notexist"

	// false if account with name has not been created before
	{
		exists, err := account.Exists(context.Background())
		require.NoError(t, err)
		assert.False(t, exists)
	}

	err = account.Insert(context.Background())
	require.NoError(t, err)

	// true if account has already been created
	{
		exists, err := account.Exists(context.Background())
		require.NoError(t, err)
		assert.True(t, exists)
	}

	account.AccountName = ""
	account.ID = 0

	// false and error if account does not include any fields
	{
		exists, err := account.Exists(context.Background())
		require.Error(t, err)
		assert.False(t, exists)
	}
}
