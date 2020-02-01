package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	//	"fmt"
	"io/ioutil"
	"log"
	"strings"

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

	wholeExcelFile, err := xlFile.ToSlice()
	if err != nil {
		log.Fatal(err)
	}
	saveJSON(wholeExcelFile, "eksonimi-giam.json")
	for i, sheet := range wholeExcelFile {
		if len(sheet) > 0 {
			// sheet has some content
			saveCSV(wholeExcelFile[i], fmt.Sprintf("eksonimi-giam-%d.csv", i))
		}

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
				case "A", "B":
					nameSlTag = "name:sl"
					f.SetProperty(nameSlTag, ex.NameSl)
				default:
					errMsg = errMsg + "; Invalid recommendedUse " + strconv.Quote(ex.RecommendedUse) + " of standardized exonym"
				}

			case "nestandardiziran", "non-standardized":
				switch strings.ToUpper(ex.RecommendedUse) {
				case "A", "B":
					nameSlTag = "name:sl"
					f.SetProperty(nameSlTag, ex.NameSl)
				case "C", "D", "E":
					nameSlTag = "alt_name:sl"
					if hasValue(ex.NameSlAlt) {
						ex.NameSlAlt = ex.NameSl + ";" + ex.NameSlAlt
					} else {
						ex.NameSlAlt = ex.NameSl
					}
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

			setOptionalProperty(f, "name:etymology", ex.Etymology)
			setOptionalProperty(f, "note", ex.Note)

			if errMsg != "" {
				errMsg = strings.TrimPrefix(errMsg, "; ")
				log.Println("ERROR:", ex.NameSl, "-", errMsg)
				f.SetProperty("error", errMsg)
				f.SetProperty("marker-color", "#ff0000")
			}
			featureCollection.AddFeature(f)
		}
	}

	log.Printf("Read %d features.", len(featureCollection.Features))

	saveJSON(featureCollection, "eksonimi.geojson")
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

func saveCSV(obj [][]string, csvFilename string) {
	csvfile, err := os.Create(csvFilename)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer csvfile.Close()
	csvwriter := csv.NewWriter(csvfile)
	defer csvwriter.Flush()

	err = csvwriter.WriteAll(obj)
	if err != nil {
		log.Fatalf("failed writing CSV file: %s", err)
	}

	stat, err := csvfile.Stat()
	if err != nil {
		log.Fatalf("Failed to stat CSV file: %s", err)
	}
	log.Printf("Saved %d Bytes to %s.", stat.Size(), csvFilename)

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
