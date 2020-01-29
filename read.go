package main

import (
	"encoding/json"
	"errors"
	"math"

	//	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/tealeg/xlsx"

	//"github.com/serjvanilla/go-overpass"
	geojson "github.com/paulmach/go.geojson"
)

func main() {
	// http://ungegn.zrc-sazu.si/Portals/7/VELIKA%20PREGLEDNICA_slo.xlsx
	excelFileName := "VELIKA PREGLEDNICA_slo.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		log.Fatal(err)
	}

	exonyms := make([]Exonym, 0)
	featureCollection := geojson.NewFeatureCollection()
	for _, sheet := range xlFile.Sheets {
		for i, row := range sheet.Rows {
			if i == 0 || row.Cells[0].Value == "" {
				continue
			}
			ex := &Exonym{}
			err := row.ReadStruct(ex)
			if err != nil {
				log.Fatal("Error reading to struct:", err)
			}
			errMsg := ""
			ex.Lat, err = ParseDMS(ex.LatDMS)
			if err != nil {
				errMsg = errMsg + err.Error()
			}
			ex.Lon, err = ParseDMS(ex.LonDMS)
			if err != nil {
				errMsg = errMsg + err.Error()
			}
			//fmt.Println(ex)
			// for _, lang := range strings.Split(ex.LangOrig,"/") {
			//	fmt.Println(lang+",")
			//}
			f := geojson.NewPointFeature([]float64{ex.Lon, ex.Lat})
			f.SetProperty("name", ex.NameOrig)
			f.SetProperty("name:sl", ex.NameSl)
			f.SetProperty("source:name:sl", "GIAM")

			setFeatureType(f, ex.FeatureType)

			if hasValue(ex.NameSlAlt) {
				ex.NameSlAlt = strings.ReplaceAll(ex.NameSlAlt, "/", ";")
				ex.NameSlAlt = strings.ReplaceAll(ex.NameSlAlt, ", ", ";")
				f.SetProperty("alt_name:sl", ex.NameSlAlt)
			}

			setOptionalProperty(f, "name:en", ex.NameEn)
			setOptionalProperty(f, "name:fr", ex.NameFr)
			setOptionalProperty(f, "name:de", ex.NameDe)
			setOptionalProperty(f, "name:es", ex.NameEs)
			setOptionalProperty(f, "name:ru", ex.NameRu)
			setOptionalProperty(f, "name:it", ex.NameIt)
			setOptionalProperty(f, "name:hr", ex.NameHr)
			setOptionalProperty(f, "name:hu", ex.NameHu)

			setOptionalProperty(f, "name:etymology", ex.Etymology)
			setOptionalProperty(f, "note", ex.Note)

			if errMsg != "" {
				log.Println("ERROR:", errMsg)
				f.SetProperty("error", errMsg)
				f.SetProperty("marker-color", "#ff0000")
			}
			featureCollection.AddFeature(f)
			exonyms = append(exonyms, *ex)
		}
	}

	log.Printf("Read %d features.", len(featureCollection.Features))

	saveJSON(featureCollection, "eksonimi.geojson")
	saveJSON(exonyms, "eksonimi-giam.json")
}

func saveJSON(obj interface{}, jsonFilename string) {
	rawJSON, err := json.MarshalIndent(obj, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(jsonFilename, rawJSON, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Saved %d Bytes to %s.", len(rawJSON), jsonFilename)
}

func setFeatureType(f *geojson.Feature, featureType string) {
	f.SetProperty("semantic-type", featureType)

	if strings.Contains(featureType, "upravn") {
		f.SetProperty("marker-symbol", "circle-stroked")
	}

	if strings.Contains(featureType, "naselj") {
		f.SetProperty("marker-symbol", "city")
	}

	if strings.Contains(featureType, "zgodovinsk") {
		f.SetProperty("marker-color", "#D2B48C")
	}

	if strings.Contains(featureType, "država") {
		f.SetProperty("marker-color", "#D400D4")
		f.SetProperty("marker-size", "large")
		f.SetProperty("marker-symbol", "embassy")
	}
}

func setOptionalProperty(f *geojson.Feature, key string, value string) {
	if !hasValue(value) {
		return
	}

	f.SetProperty(key, strings.TrimSpace(value))
}

func hasValue(in string) bool {
	return in != "" && in != "–" && in != "0" && in != "ni" && in != "???" && strings.TrimSpace(in) != ""
}

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
