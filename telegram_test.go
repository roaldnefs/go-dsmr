package dsmr_test

import (
	"github.com/roaldnefs/go-dsmr"
	"testing"
)

var testTelegram = `/ISK5\2M550T-1012

0-0:1.0.0(190718204947S)
!0000`

// TestTelegramHeader tests parsing the telegram header
func TestTelegramHeader(t *testing.T) {
	telegram, err := dsmr.ParseTelegram(testTelegram)
	if err != nil {
		t.Error("parsing telegram returned an error")
	}

	if telegram.Header != `/ISK5\2M550T-1012` {
		t.Error(`expected '/ISK5\2M550T-1012' as header got '` + telegram.Header + "'")
	}
}
