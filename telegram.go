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
