// Copyright (c) 2017, Bas van der Lei
// Modified work Copyright (c) 2019 Roald Nefs

package dsmr

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// DateTimeFormat used in a telegram in the YYMMDDhhmmssX format. The ASCII
// presentation of timestamp with Year, Month, Day, Hour, Minute, Second,
// and an indication whether DST is active (X=S) or DST is not active (X=W).
const DateTimeFormat = "060102150405"

var (
	// OBIS Reduces ID-code and value
	objectRegexp = regexp.MustCompile("([0-9]+-[0-9]+:[0-9]+.[0-9]+.[0-9]+)\\((.*)\\)")
	// Value and unit
	// In the DSMR the unit is optional
	valueRegexp = regexp.MustCompile("([^*]+)\\*?(.*)")
)

// Telegram represent a single DSMR P1 message.
type Telegram struct {
	Header   string
	Version  string    // Version information for P1 output.
	DateTime time.Time // Date-time stamp of the P1 message.

	DataObjects map[string]DataObject
}

// MeterReadingElectricityDeliveredToClientTariff1 returns the meter reading
// electricity delivered to client (Tariff 1) in 0,001 kWh.
func (t Telegram) MeterReadingElectricityDeliveredToClientTariff1() (string, bool) {
	if do, ok := t.DataObjects["1-0:1.8.1"]; ok {
		return do.Value, true
	}
	return "", false
}

// MeterReadingElectricityDeliveredToClientTariff2 returns the meter reading
// electricity delivered to client (Tariff 2) in 0,001 kWh.
func (t Telegram) MeterReadingElectricityDeliveredToClientTariff2() (string, bool) {
	if do, ok := t.DataObjects["1-0:1.8.2"]; ok {
		return do.Value, true
	}
	return "", false
}

// MeterReadingElectricityDeliveredByClientTariff1 returns the meter reading
// electricity delivered by client (Tariff 1) in 0,001 kWh.
func (t Telegram) MeterReadingElectricityDeliveredByClientTariff1() (string, bool) {
	if do, ok := t.DataObjects["1-0:2.8.1"]; ok {
		return do.Value, true
	}
	return "", false
}

// MeterReadingElectricityDeliveredByClientTariff2 returns the meter reading
// electricity delivered by client (Tariff 2) in 0,001 kWh.
func (t Telegram) MeterReadingElectricityDeliveredByClientTariff2() (string, bool) {
	if do, ok := t.DataObjects["1-0:2.8.2"]; ok {
		return do.Value, true
	}
	return "", false
}

// TariffIndicatorElectricity returns the tariff indicator that can be used to
// switch tariff dependent loads e.g. boilers.
func (t Telegram) TariffIndicatorElectricity() (string, bool) {
	if do, ok := t.DataObjects["0-0:96.14.0"]; ok {
		return do.Value, true
	}
	return "", false
}

// ActualElectricityPowerDelivered returns the actual electricity power
// delivered (+P) in 1 Watt resolution.
func (t Telegram) ActualElectricityPowerDelivered() (string, bool) {
	if do, ok := t.DataObjects["1-0:1.7.0"]; ok {
		return do.Value, true
	}
	return "", false
}

// ActualElectricityPowerReceived returns the actual electricity power
// received (-P) in 1 Watt resolution.
func (t Telegram) ActualElectricityPowerReceived() (string, bool) {
	if do, ok := t.DataObjects["1-0:2.7.0"]; ok {
		return do.Value, true
	}
	return "", false
}

// NumberOfPowerFailuresInAnyPhase returns the number of power failures in any
// phase.
func (t Telegram) NumberOfPowerFailuresInAnyPhase() (string, bool) {
	if do, ok := t.DataObjects["0-0:96.7.21"]; ok {
		return do.Value, true
	}
	return "", false
}

// NumberOfLongPowerFailuresInAnyPhase returns the number of power failures in any
// phase.
func (t Telegram) NumberOfLongPowerFailuresInAnyPhase() (string, bool) {
	if do, ok := t.DataObjects["0-0:96.7.9"]; ok {
		return do.Value, true
	}
	return "", false
}

// TODO update regex so it will match:
//   None: 1-0:99.97.0()
//   Two: 1-0:99.97.0(2)(0-0:96.7.19)(101208152415W)(0000000240*s)(101208151004W)(0000000301*s)
// PowerFailureEventLog returns the power failure event log (long power
// failures).
func (t Telegram) PowerFailureEventLog() (string, bool) {
	if do, ok := t.DataObjects["1-0:99.97.0"]; ok {
		return do.Value, true
	}
	return "", false
}

// NumberOfVoltageSagsInPhaseL1 returns the number of voltage sags in phase L1.
func (t Telegram) NumberOfVoltageSagsInPhaseL1() (string, bool) {
	if do, ok := t.DataObjects["1-0:32.32.0"]; ok {
		return do.Value, true
	}
	return "", false
}

// NumberOfVoltageSagsInPhaseL2 returns the number of voltage sags in phase L2.
func (t Telegram) NumberOfVoltageSagsInPhaseL2() (string, bool) {
	if do, ok := t.DataObjects["1-0:52.32.0"]; ok {
		return do.Value, true
	}
	return "", false
}

// NumberOfVoltageSagsInPhaseL3 returns the number of voltage sags in phase L3.
func (t Telegram) NumberOfVoltageSagsInPhaseL3() (string, bool) {
	if do, ok := t.DataObjects["1-0:72.32.0"]; ok {
		return do.Value, true
	}
	return "", false
}

// NumberOfVoltageSwellsInPhaseL1 returns the number of voltage swells in phase L1.
func (t Telegram) NumberOfVoltageSwellsInPhaseL1() (string, bool) {
	if do, ok := t.DataObjects["1-0:32.36.0"]; ok {
		return do.Value, true
	}
	return "", false
}

// NumberOfVoltageSwellsInPhaseL2 returns the number of voltage swells in phase L2.
func (t Telegram) NumberOfVoltageSwellsInPhaseL2() (string, bool) {
	if do, ok := t.DataObjects["1-0:52.36.0"]; ok {
		return do.Value, true
	}
	return "", false
}

// NumberOfVoltageSwellsInPhaseL3 returns the number of voltage swells in phase L3.
func (t Telegram) NumberOfVoltageSwellsInPhaseL3() (string, bool) {
	if do, ok := t.DataObjects["1-0:72.36.0"]; ok {
		return do.Value, true
	}
	return "", false
}

// TextMessage returns text message of max 1024 characters.
func (t Telegram) TextMessage() (string, bool) {
	if do, ok := t.DataObjects["0-0:96.13.0"]; ok {
		return do.Value, true
	}
	return "", false
}

// InstantaneousVoltageL1 returns the instantaneous voltage L1 in V resolution.
func (t Telegram) InstantaneousVoltageL1() (string, bool) {
	if do, ok := t.DataObjects["1-0:32.7.0"]; ok {
		return do.Value, true
	}
	return "", false
}

// InstantaneousVoltageL2 returns the instantaneous voltage L3 in V resolution.
func (t Telegram) InstantaneousVoltageL2() (string, bool) {
	if do, ok := t.DataObjects["1-0:52.7.0"]; ok {
		return do.Value, true
	}
	return "", false
}

// InstantaneousVoltageL3 returns the instantaneous voltage L3 in V resolution.
func (t Telegram) InstantaneousVoltageL3() (string, bool) {
	if do, ok := t.DataObjects["1-0:72.7.0"]; ok {
		return do.Value, true
	}
	return "", false
}

// InstantaneousCurrentL1 returns the instantaneous current L1 in A resolution.
func (t Telegram) InstantaneousCurrentL1() (string, bool) {
	if do, ok := t.DataObjects["1-0:31.7.0"]; ok {
		return do.Value, true
	}
	return "", false
}

// InstantaneousCurrentL2 returns the instantaneous current L2 in A resolution.
func (t Telegram) InstantaneousCurrentL2() (string, bool) {
	if do, ok := t.DataObjects["1-0:51.7.0"]; ok {
		return do.Value, true
	}
	return "", false
}

// InstantaneousCurrentL3 returns the instantaneous current L3 in A resolution.
func (t Telegram) InstantaneousCurrentL3() (string, bool) {
	if do, ok := t.DataObjects["1-0:71.7.0"]; ok {
		return do.Value, true
	}
	return "", false
}

// DataObject represent data object and it's reference to the OBIS as defined
// in EN-EN-IEC 62056-61:2002 Electricity metering – Data exchange for meter
// reading, tariff and load control – Part 61: OBIS Object Identification
// System.
type DataObject struct {
	OBIS  string // OBIS reduced ID-code
	Value string
	Unit  string
}

// ParseTelegram will parse the DSMR telegram.
func ParseTelegram(telegram string) (t Telegram, err error) {
	t.DataObjects = make(map[string]DataObject)

	for _, line := range strings.Split(telegram, "\n") {
		line = strings.TrimSpace(line)

		// Skip empty line or the footer of the telegram
		if line == "" || line[0] == '!' {
			continue
		}

		// The header of the telegram
		if line[0] == '/' {
			t.Header = line
			continue
		}

		do, err := ParseDataObject(line)
		if err != nil {
			// TODO logging
			continue
		}

		switch do.OBIS {
		// Version information for P1 output.
		case "1-3:0.2.8":
			t.Version = do.Value
		// Date time of P1 output.
		case "0-0:1.0.0":
			if len(do.Value) > 2 {
				// Remove the DST indicator from the timestamp
				rawDateTime := do.Value[:len(do.Value)-1]
				location, err := time.LoadLocation("Europe/Amsterdam")
				if err != nil {
					// TODO logging
					continue
				}
				dateTime, err := time.ParseInLocation(DateTimeFormat, rawDateTime, location)
				if err != nil {
					// TODO logging
					continue
				}
				t.DateTime = dateTime
			}
		default:
			t.DataObjects[do.OBIS] = do
		}
	}
	return t, nil
}

// ParseDataObject will parse a single line into a DataObject.
func ParseDataObject(do string) (DataObject, error) {
	// Extract the OBIS reduced ID-code and the corresponding value, e.g:
	// 1-3:0.2.8(50) --> 1-3:0.2.8 and (50)
	match := objectRegexp.FindStringSubmatch(strings.TrimSpace(do))
	if match == nil || len(match) < 3 {
		return DataObject{}, fmt.Errorf("no valid DSMR object found")
	}
	obis := match[1]
	rawValue := match[2]

	// Extract the value and the unit from the raw value, e.g:
	// (000099.999*kWh) -> 000099.999 and kWh
	match = valueRegexp.FindStringSubmatch(rawValue)
	if match == nil {
		return DataObject{
			OBIS:  obis,
			Value: rawValue,
		}, nil
	}
	if len(match) > 2 {
		return DataObject{
			OBIS:  obis,
			Value: match[1],
			Unit:  match[2],
		}, nil
	}
	return DataObject{
		OBIS:  obis,
		Value: match[1],
	}, nil
}
