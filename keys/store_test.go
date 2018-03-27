package keys_test

import (
	"context"
	"os"
	"testing"

	"github.com/joyent/triton-service-groups/keys"
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

	store := keys.NewStore(db.Conn)
	require.NotNil(t, store)

	key := keys.New(store)
	require.NotNil(t, key)

	key.Name = "TSG_Management"
	key.Fingerprint = "12:23:34:45:56:67:78:89:90:0A:AB:BC:CD:DE:AD:01"
	key.Material = "this is key material"

	err = key.Insert(context.Background())
	require.NoError(t, err)

	require.NotZero(t, key.ID)

	found, err := store.FindByID(context.Background(), key.ID)
	require.NoError(t, err)

	assert.Equal(t, key.ID, found.ID)
	assert.Equal(t, key.Name, found.Name)
	assert.Equal(t, key.Fingerprint, found.Fingerprint)
	assert.Equal(t, key.Material, found.Material)
	assert.Equal(t, key.CreatedAt, found.CreatedAt)
	assert.Equal(t, key.UpdatedAt, found.UpdatedAt)
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

	store := keys.NewStore(db.Conn)
	require.NotNil(t, store)

	key := keys.New(store)
	require.NotNil(t, key)

	key.Name = "TSG_Management"
	key.Fingerprint = "12:23:34:45:56:67:78:89:90:0A:AB:BC:CD:DE:AD:01"
	key.Material = "this is key material"

	err = key.Insert(context.Background())
	require.NoError(t, err)

	require.NotZero(t, key.ID)

	found, err := store.FindByName(context.Background(), key.Name)
	require.NoError(t, err)

	assert.Equal(t, key.ID, found.ID)
	assert.Equal(t, key.Name, found.Name)
	assert.Equal(t, key.Fingerprint, found.Fingerprint)
	assert.Equal(t, key.Material, found.Material)
	assert.Equal(t, key.CreatedAt, found.CreatedAt)
	assert.Equal(t, key.UpdatedAt, found.UpdatedAt)
}
