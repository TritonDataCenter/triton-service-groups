package accounts_test

import (
	"context"
	"os"
	"testing"

	"github.com/joyent/triton-service-groups/accounts"
	"github.com/joyent/triton-service-groups/keys"
	"github.com/joyent/triton-service-groups/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	if os.Getenv("TSG_TEST") == "" {
		t.Skip("Acceptance tests skipped unless env 'TSG_TEST=1' set")
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

	accountName := "johndoe"
	tritonUUID := "f5435e8b-70b8-4e4d-8c59-1dbe5d100b5b"

	account.AccountName = accountName
	account.TritonUUID = tritonUUID

	err = account.Insert(context.Background())
	require.NoError(t, err)

	assert.NotZero(t, account.ID)
	assert.Equal(t, account.AccountName, accountName)
	assert.Equal(t, account.TritonUUID, tritonUUID)
	assert.NotZero(t, account.CreatedAt)
	assert.NotZero(t, account.UpdatedAt)
	assert.Equal(t, account.CreatedAt, account.UpdatedAt)
}

func TestSave(t *testing.T) {
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

	// with an empty KeyID
	{
		account := accounts.New(store)
		require.NotNil(t, account)

		accountName := "secondname"
		tritonUUID := "f5435e8b-70b8-4e4d-8c59-1dbe5d100b5b"

		account.AccountName = "firstname"
		account.TritonUUID = tritonUUID

		err = account.Insert(context.Background())
		require.NoError(t, err)

		assert.Empty(t, account.KeyID)
		assert.Equal(t, account.CreatedAt, account.UpdatedAt)

		account.AccountName = accountName

		err = account.Save(context.Background())
		require.NoError(t, err)

		acct, err := store.FindByID(context.Background(), account.ID)
		if err != nil {
			t.Error(err)
		}

		assert.NotZero(t, acct.ID)
		assert.Equal(t, accountName, acct.AccountName)
		assert.Equal(t, tritonUUID, acct.TritonUUID)
		assert.Equal(t, "", acct.KeyID)
		assert.NotZero(t, acct.CreatedAt)
		assert.NotZero(t, acct.UpdatedAt)
		assert.NotEqual(t, acct.CreatedAt, acct.UpdatedAt)
	}

	// with a valid KeyID
	{
		account := accounts.New(store)
		require.NotNil(t, account)

		accountName := "seconduser"
		tritonUUID := "f5435e8b-70b8-4e4d-8c59-1dbe5d100b5b"

		account.AccountName = "demouser"
		account.TritonUUID = tritonUUID

		err = account.Insert(context.Background())
		require.NoError(t, err)

		assert.Empty(t, account.KeyID)
		assert.Equal(t, account.CreatedAt, account.UpdatedAt)

		keyStore := keys.NewStore(db.Conn)
		require.NotNil(t, keyStore)

		key := keys.New(keyStore)
		require.NotNil(t, key)

		key.Name = "testkey"
		key.Fingerprint = "blahblahblah"
		key.Material = "blahblahblah"

		err = key.Insert(context.Background())
		require.NoError(t, err)

		account.AccountName = accountName
		account.KeyID = key.ID

		err = account.Save(context.Background())
		require.NoError(t, err)

		acct, err := store.FindByID(context.Background(), account.ID)
		if err != nil {
			t.Error(err)
		}

		assert.NotZero(t, acct.ID)
		assert.Equal(t, accountName, acct.AccountName)
		assert.Equal(t, tritonUUID, acct.TritonUUID)
		assert.Equal(t, key.ID, acct.KeyID)
		assert.NotZero(t, acct.CreatedAt)
		assert.NotZero(t, acct.UpdatedAt)
		assert.NotEqual(t, acct.CreatedAt, acct.UpdatedAt)
	}
}

func TestExists(t *testing.T) {
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
	account.ID = ""

	// false and error if account does not include any fields
	{
		exists, err := account.Exists(context.Background())
		require.Error(t, err)
		assert.False(t, exists)
	}
}
