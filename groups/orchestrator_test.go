package groups_v1

import (
	"testing"
)

func TestBase64Encode(t *testing.T) {
	tests := []struct {
		value   string
		expeted string
	}{
		// No bytes to encode, thus no output.
		{
			"",
			"",
		},
		{
			"tsg",
			"dHNn",
		},
	}

	for _, tt := range tests {
		actual := base64Encode(tt.value)
		if tt.expeted != actual {
			t.Errorf("expected value %#v, got %#v", tt.expeted, actual)
		}
	}

}

func TestEscapeNewlines(t *testing.T) {
	tests := []struct {
		value   string
		expeted string
	}{
		{
			"",
			"",
		},
		{
			"\n",
			"\\n",
		},
		// Only escape LF (\n), as CR+LF (\r\n) is not supported.
		{
			"\r\n",
			"\r\\n",
		},
		{
			"\t",
			"\t",
		},
	}

	for _, tt := range tests {
		actual := escapeNewlines(tt.value)
		if tt.expeted != actual {
			t.Errorf("expected value %#v, got %#v", tt.expeted, actual)
		}
	}
}
