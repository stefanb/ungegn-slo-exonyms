package main

type Exonym struct {
	ID          int     `xlsx:"0"`
	NameSl      string  `xlsx:"1"`
	NameOrig    string  `xlsx:"4"`
	LangOrig    string  `xlsx:"5"`
	FeatureType string  `xlsx:"8"`
	LatDMS      string  `xlsx:"9"`
	Lat         float64 `xlsx:"-" json:"-"`
	LonDMS      string  `xlsx:"10"`
	Lon         float64 `xlsx:"-" json:"-"`
	//BoolVal bool `xlsx:"4"`
	NameSlAlt string `xlsx:"14"`
	NameEn    string `xlsx:"25"`
	NameFr    string `xlsx:"26"`
	NameDe    string `xlsx:"27"`
	NameEs    string `xlsx:"28"`
	NameRu    string `xlsx:"29"`
	NameIt    string `xlsx:"30"`
	NameHr    string `xlsx:"31"`
	NameHu    string `xlsx:"32"`

	Etymology string `xlsx:"33"`
	Note      string `xlsx:"34"`
	Wikidata  string `xlsx:"35"`
}
