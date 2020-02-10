package dsmr_test

import (
	"github.com/roaldnefs/go-dsmr"
	"testing"
)

var testTelegram = `/ISK5\2M550T-1012

0-0:1.0.0(190718204947S)
1-0:1.7.0(00.056*kW)
1-0:2.7.0(00.100*kW)
0-0:96.14.0(0001)
0-1:24.2.1(191118114002W)(00000.003*m3)
0-2:24.2.1(200208141004W)(00417.143*m3)
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

func TestActualElectricityPowerDelivered(t *testing.T) {
	telegram, err := dsmr.ParseTelegram(testTelegram)
	if err != nil {
		t.Error("parsing telegram returned an error")
	}

	delivered, ok := telegram.ActualElectricityPowerDelivered()
	if !ok {
		t.Error("retrieving actual electricity power delivered failed")
	}

	do := telegram.DataObjects["1-0:1.7.0"]
	if delivered != "00.056" || do.Value != delivered {
		t.Error("wrong value retrieved for actual electricity power delivered expected 00.056 got " + delivered)
	}

	if do.Unit != "kW" {
		t.Error("wrong unit retrieved for actual electricity power delivered expected kW got " + do.Unit)
	}
}

func TestActualElectricityPowerReceived(t *testing.T) {
	telegram, err := dsmr.ParseTelegram(testTelegram)
	if err != nil {
		t.Error("parsing telegram returned an error")
	}

	received, ok := telegram.ActualElectricityPowerReceived()
	if !ok {
		t.Error("retrieving actual electricity power received failed")
	}

	do := telegram.DataObjects["1-0:2.7.0"]
	if received != "00.100" || do.Value != received {
		t.Error("wrong value retrieved for actual electricity power received expected 00.100 got " + received)
	}

	if do.Unit != "kW" {
		t.Error("wrong unit retrieved for actual electricity power received expected kW got " + do.Unit)
	}
}

func TestTarifIndicatorElectricity(t *testing.T) {
	telegram, err := dsmr.ParseTelegram(testTelegram)
	if err != nil {
		t.Error("parsing telegram returned an error")
	}

	tarif, ok := telegram.TariffIndicatorElectricity()
	if !ok {
		t.Error("retrieving tarif indicator for electricity failed")
	}

	do := telegram.DataObjects["0-0:96.14.0"]
	if tarif != "0001" || do.Value != tarif {
		t.Error("wrong value retrieved for tarif indicator electricity expected 0001 got " + tarif)
	}
}

func TestMeterReadingGasDeliveredToClient(t *testing.T) {
	telegram, err := dsmr.ParseTelegram(testTelegram)
	if err != nil {
		t.Error("parsing telegram returned an error")
	}

	delivered, ok := telegram.MeterReadingGasDeliveredToClient(2)
	if !ok {
		t.Error("retrieving meter reading gas delivered to client (channel 2) failed")
	}

	do := telegram.DataObjects["0-2:24.2.1"]
	if delivered != "00417.143" || do.Value != delivered {
		t.Error("wrong value retrieved for meter reading gas delivered to client (channel 2) expected 00417.143 got " + delivered)
	}

	if do.Unit != "m3" {
		t.Error("Wrong unit retrieved for meter reading gas delivered to client (channel 2). Expected m3 got " + do.Unit)
	}
}