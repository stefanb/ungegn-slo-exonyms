package main

import (
	"errors"
	"fmt"
	"math"
	"regexp"

	"strconv"
	"strings"
)

//"125° 8′ 6″ Z"
var dmsRegex = regexp.MustCompile(`^\s*([0-1]?[0-9]?[0-9])°\s*([0-5]?[0-9])′\s*([0-5]?[0-9])″\s*([SsJjVvZz])\s*$`)

// ParseDMS takes the input DMS notation string (eg "125° 8′ 6″ Z") and converts it to decimal degrees
func ParseDMS(dms string) (float64, error) {
	if !hasValue(dms) {
		return 0, errors.New("No dms value to parse: " + strconv.Quote(dms))
	}

	matches := dmsRegex.FindStringSubmatch(dms)
	if len(matches) != 5 {
		return 0, fmt.Errorf("Error parsing %s", strconv.Quote(dms))
	}

	q := 0
	var dmax int

	switch strings.ToUpper(matches[4]) {
	case "J":
		q = -1
		dmax = 90
	case "Z":
		q = -1
		dmax = 180
	case "S":
		q = 1
		dmax = 90
	case "V":
		q = 1
		dmax = 180
	default:
		return 0, fmt.Errorf("Invalid direction %s", strconv.Quote(matches[4]))
	}

	// errors disallowed by regex :)
	d, _ := strconv.Atoi(matches[1])
	m, _ := strconv.Atoi(matches[2])
	s, _ := strconv.Atoi(matches[3])

	deg := (float64(d) + float64(m)/60 + float64(s)/60/60)

	if d > dmax {
		return 0, fmt.Errorf("%d deegrees to large for direction %s", d, strconv.Quote(matches[4]))
	}
	return float64(q) * math.Round(deg*10000) / 10000, nil
}
