package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	//	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/tealeg/xlsx"
	//"github.com/serjvanilla/go-overpass"
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

	saveJSON(getGeoJSON(xlFile.Sheets[0]), "eksonimi.geojson")
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

func hasValue(in string) bool {
	return in != "" && in != "–" && in != "0" && in != "ni" && in != "???" && strings.TrimSpace(in) != ""
}
