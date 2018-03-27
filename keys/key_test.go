package keys_test

import (
	"context"
	"os"
	"testing"
	"time"

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

	store := keys.NewStore(db.Conn)
	require.NotNil(t, store)

	key := keys.New(store)
	require.NotNil(t, key)
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

	store := keys.NewStore(db.Conn)
	require.NotNil(t, store)

	key := keys.New(store)
	require.NotNil(t, key)

	key.Name = "TSG_Management"
	key.Fingerprint = "12:23:34:45:56:67:78:89:90:0A:AB:BC:CD:DE:AD:01"
	key.Material = "this is key material"

	err = key.Insert(context.Background())
	require.NoError(t, err)

	assert.NotZero(t, key.ID)
	assert.Equal(t, key.Name, "TSG_Management")
	assert.Equal(t, key.Fingerprint, "12:23:34:45:56:67:78:89:90:0A:AB:BC:CD:DE:AD:01")
	assert.Equal(t, key.Material, "this is key material")
	assert.False(t, key.Archived)
	assert.NotZero(t, key.CreatedAt)
	assert.NotZero(t, key.UpdatedAt)
	assert.Equal(t, key.CreatedAt, key.UpdatedAt)
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

	store := keys.NewStore(db.Conn)
	require.NotNil(t, store)

	key := keys.New(store)
	require.NotNil(t, key)

	key.Name = "TSG_Management"
	key.Fingerprint = "12:23:34:45:56:67:78:89:90:0A:AB:BC:CD:DE:AD:01"
	key.Material = "this is key material"

	err = key.Insert(context.Background())
	require.NoError(t, err)

	assert.Equal(t, key.CreatedAt, key.UpdatedAt)

	key.Name = "hackerman"

	time.Sleep(1 * time.Second)

	err = key.Save(context.Background())
	require.NoError(t, err)

	assert.NotZero(t, key.ID)
	assert.Equal(t, key.Name, "hackerman")
	assert.Equal(t, key.Fingerprint, "12:23:34:45:56:67:78:89:90:0A:AB:BC:CD:DE:AD:01")
	assert.Equal(t, key.Material, "this is key material")
	assert.False(t, key.Archived)
	assert.NotZero(t, key.CreatedAt)
	assert.NotZero(t, key.UpdatedAt)
	assert.NotEqual(t, key.CreatedAt, key.UpdatedAt)
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

	store := keys.NewStore(db.Conn)
	require.NotNil(t, store)

	created := keys.New(store)
	require.NotNil(t, created)
	created.Name = "firstcreate"
	created.Insert(context.Background())

	newKey := keys.New(store)
	require.NotNil(t, newKey)

	newKey.ID = created.ID

	// true if key shares ID with a previously created key
	{
		exists, err := newKey.Exists(context.Background())
		require.NoError(t, err)
		assert.True(t, exists)
	}

	key := keys.New(store)
	require.NotNil(t, key)

	key.Name = "notexist"

	// false if key with name has not been created before
	{
		exists, err := key.Exists(context.Background())
		require.NoError(t, err)
		assert.False(t, exists)
	}

	err = key.Insert(context.Background())
	require.NoError(t, err)

	// true if key has already been created
	{
		exists, err := key.Exists(context.Background())
		require.NoError(t, err)
		assert.True(t, exists)
	}

	key.Name = ""
	key.ID = ""

	// false and error if key does not include any fields
	{
		exists, err := key.Exists(context.Background())
		require.Error(t, err)
		assert.False(t, exists)
	}
}
