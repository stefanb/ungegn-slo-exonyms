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
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				text := cell.String()
				fmt.Printf("%s\n", text)
			}
		}
	}
}
