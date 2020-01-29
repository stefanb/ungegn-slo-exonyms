package main

import (
	"errors"
	"math"

	"strconv"
	"strings"
	"unicode/utf8"
)

func ParseDMS(dms string) (float64, error) {
	if !hasValue(dms) {
		return 0, errors.New("No dms value to parse: " + strconv.Quote(dms))
	}

	dmsIn := dms
	dms = strings.ReplaceAll(dms, "\u00b0", "\u00b0 ")
	dms = strings.ReplaceAll(dms, "\u2032", "\u2032 ")
	dms = strings.ReplaceAll(dms, "\u2033", "\u2033 ")
	d, m, s, q := 0, 0, 0, 1
	errMsg := ""
	for _, part := range strings.Fields(dms) {
		if len(part) == 1 {
			if part == "J" || part == "Z" {
				q = -1
			}
			continue
		}

		unit, ulen := utf8.DecodeLastRuneInString(part)
		if len(part) == ulen {
			errMsg = errMsg + "No value to parse from: " + part + " in: " + dmsIn
			continue
		}
		val, err := strconv.Atoi(part[:len(part)-ulen])
		if err != nil {
			errMsg = errMsg + "Error parsing: " + strconv.Quote(part) + " from: " + strconv.Quote(dmsIn) + " reason: " + err.Error()
			continue
		}

		switch unit {
		case 176: //'\u00b0':
			d = val
		case 8242: //'\u2032':
			m = val
		case 8243: //'\u2033':
			s = val
		default:
			errMsg = errMsg + "Unknown angle unit:" + string(unit)
		}
	}

	var err error
	if errMsg != "" {
		err = errors.New(errMsg)
	}

	deg := float64(q) * (float64(d) + float64(m)/60 + float64(s)/60/60)
	return math.Round(deg*10000) / 10000, err
}
