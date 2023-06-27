package lib

import (
	"testing"

	"github.com/google/uuid"
)

func TestToHex(t *testing.T) {
	guuid := uuid.NewString()
	hex := decodeUuid(guuid)
	huuid := encodeUuid(hex)

	if guuid != huuid {
		t.Errorf("Uuids are different %s and %s", guuid, huuid)
	}
}
