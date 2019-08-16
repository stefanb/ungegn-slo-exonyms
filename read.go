package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
)

func main() {
	// http://ungegn.zrc-sazu.si/Portals/7/VELIKA%20PREGLEDNICA_slo.xlsx
	excelFileName := "VELIKA PREGLEDNICA_slo.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		log.Fatal(err)
	}
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
			fmt.Println(ex)
		}
	}
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
