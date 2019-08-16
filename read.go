package main

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"
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
			ex.Lat = ParseDMS(ex.LatDMS)
			ex.Lon = ParseDMS(ex.LonDMS)
			//fmt.Println(ex)
			// for _, lang := range strings.Split(ex.LangOrig,"/") {
			//	fmt.Println(lang+",")
			//}
			f := geojson.NewPointFeature([]float64{ex.Lon, ex.Lat})
			f.SetProperty("name:sl", ex.NameSl)
			f.SetProperty("source:name:sl", "ungegn")
			featureCollection.AddFeature(f)
		}
	}

	rawJSON, err := json.MarshalIndent(featureCollection, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	geoJsonFilename := "eksonimi.geojson"

	err = ioutil.WriteFile(geoJsonFilename, rawJSON, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Saved %d addresses to %s.", len(featureCollection.Features), geoJsonFilename)

}

func ParseDMS(dms string) float64 {
	d, m, s, q := 0, 0, 0, 1
	for _, part := range strings.Fields(dms) {
		if len(part) == 1 {
			if part == "J" || part == "Z" {
				q = -1
			}
			continue
		}

		unit, ulen := utf8.DecodeLastRuneInString(part)
		if len(part) == ulen {
			log.Println("No value to parse from:", part)
			continue
		}
		val, err := strconv.Atoi(part[:len(part)-ulen])
		if err != nil {
			log.Println("Error parsing:", strconv.Quote(part), err)
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
			fmt.Println("Unknown angle unit:", unit)
		}

	}
	return float64(q) * (float64(d) + float64(m)/60 + float64(s)/60/60)
}

type Exonym struct {
	ID          int     `xlsx:"0"`
	NameSl      string  `xlsx:"1"`
	NameOrig    string  `xlsx:"4"`
	LangOrig    string  `xlsx:"5"`
	FeatureType string  `xlsx:"8"`
	LatDMS      string  `xlsx:"9"`
	Lat         float64 `xlsx:"-"`
	LonDMS      string  `xlsx:"10"`
	Lon         float64 `xlsx:"-"`
	//BoolVal bool `xlsx:"4"`
}
