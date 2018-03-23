package convert

import "fmt"

var EmptyUUID [16]byte

func BytesToUUID(b [16]byte) string {
	if b != EmptyUUID {
		return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	}

	return ""
}
