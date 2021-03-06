package main

import (
	"strconv"

	"log"
	"strings"

	geojson "github.com/paulmach/go.geojson"
	"github.com/tealeg/xlsx"
)

// getGeoJSON takes the excel sheet and converts it to geoJSON struct
func getGeoJSON(xlSheet *xlsx.Sheet) *geojson.FeatureCollection {

	featureCollection := geojson.NewFeatureCollection()
	for i, row := range xlSheet.Rows {
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
			errMsg = errMsg + "; " + err.Error()
		}
		ex.Lon, err = ParseDMS(ex.LonDMS)
		if err != nil {
			errMsg = errMsg + "; " + err.Error()
		}
		//fmt.Println(ex)
		// for _, lang := range strings.Split(ex.LangOrig,"/") {
		//	fmt.Println(lang+",")
		//}
		f := geojson.NewPointFeature([]float64{ex.Lon, ex.Lat})
		f.SetProperty("name", ex.NameOrig)
		var nameSlTag string

		switch ex.Status {
		case "standardiziran", "standardized":
			switch strings.ToUpper(ex.RecommendedUse) {
			case "A", "B", "C":
				nameSlTag = "name:sl"
				f.SetProperty(nameSlTag, ex.NameSl)
			default:
				errMsg = errMsg + "; Invalid recommendedUse " + strconv.Quote(ex.RecommendedUse) + " of standardized exonym"
			}

		case "nestandardiziran", "non-standardized":
			switch strings.ToUpper(ex.RecommendedUse) {
			case "A", "B", "C":
				nameSlTag = "name:sl"
				f.SetProperty(nameSlTag, ex.NameSl)
			case "D", "E":
				nameSlTag = "alt_name:sl"
				if hasValue(ex.NameSlAlt) {
					ex.NameSlAlt = ex.NameSl + ";" + ex.NameSlAlt
				} else {
					ex.NameSlAlt = ex.NameSl
				}
				f.SetProperty("marker-size", "small")
			default:
				errMsg = errMsg + "; Unknown recommendedUse " + strconv.Quote(ex.RecommendedUse)
			}
		default:
			errMsg = errMsg + "; Unknown status " + strconv.Quote(ex.Status)
		}

		f.SetProperty("source:"+nameSlTag, "GIAM")

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

		setOptionalProperty(f, "name:etymology:sl", ex.Etymology)
		setOptionalProperty(f, "note:sl", ex.Note)

		if errMsg != "" {
			errMsg = strings.TrimPrefix(errMsg, "; ")
			log.Println("ERROR:", ex.NameSl, "-", errMsg)
			f.SetProperty("error", errMsg)
			f.SetProperty("marker-color", "#ff0000")
		}
		featureCollection.AddFeature(f)
	}

	log.Printf("Read %d features.", len(featureCollection.Features))

	return featureCollection
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
