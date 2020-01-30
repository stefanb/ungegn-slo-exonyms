package main

import (
	"errors"
	"fmt"
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
	d, m, s, q := 0, 0, 0, 0
	errorMsgs := make([]string, 0, 0)
	directionAttempted := false
	for _, part := range strings.Fields(dms) {
		if len(part) == 1 {
			directionAttempted = true
			switch strings.ToUpper(part) {
			case "J", "Z":
				q = -1
			case "S", "V":
				q = 1
			default:
				errorMsgs = append(errorMsgs, "Invalid direction "+strconv.Quote(part))
			}

			continue
		}

		unit, ulen := utf8.DecodeLastRuneInString(part)
		if len(part) == ulen {
			errorMsgs = append(errorMsgs, "No value to parse from "+strconv.Quote(part))
			continue
		}
		val, err := strconv.Atoi(part[:len(part)-ulen])
		if err != nil {
			errorMsgs = append(errorMsgs, "Error parsing "+strconv.Quote(part)+": "+err.Error())
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
			errorMsgs = append(errorMsgs, "Unknown angle unit "+strconv.Quote(string(unit)))
		}
	}

	if q == 0 && !directionAttempted {
		errorMsgs = append(errorMsgs, "Missing direction")
	}
	var err error
	if len(errorMsgs) != 0 {
		err = fmt.Errorf("Error parsing %s: %s", strconv.Quote(dmsIn), strings.Join(errorMsgs, ", "))
	}

	deg := float64(q) * (float64(d) + float64(m)/60 + float64(s)/60/60)
	return math.Round(deg*10000) / 10000, err
}
