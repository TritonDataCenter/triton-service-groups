package convert_test

import (
	"testing"

	"github.com/joyent/triton-service-groups/convert"
	"github.com/stretchr/testify/assert"
)

func TestBytesToUUID(t *testing.T) {
	t1 := [16]byte{220, 197, 3, 156, 150, 183, 64, 152, 135, 16, 139, 76, 49, 97, 68, 13}
	assert.Equal(t, "dcc5039c-96b7-4098-8710-8b4c3161440d", convert.BytesToUUID(t1))

	t2 := [16]byte{109, 246, 173, 84, 58, 53, 73, 2, 151, 45, 85, 44, 46, 188, 233, 84}
	assert.Equal(t, "6df6ad54-3a35-4902-972d-552c2ebce954", convert.BytesToUUID(t2))

	t3 := [16]byte{}
	assert.Equal(t, "", convert.BytesToUUID(t3))
}
