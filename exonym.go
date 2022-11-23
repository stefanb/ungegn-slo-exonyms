package main

type Exonym struct {
	ID             int     `xlsx:"0"`
	NameSl         string  `xlsx:"2"`
	NameOrig       string  `xlsx:"5"`
	LangOrig       string  `xlsx:"6"`
	FeatureType    string  `xlsx:"20"`
	LatDMS         string  `xlsx:"12"`
	Lat            float64 `xlsx:"-" json:"-"`
	LonDMS         string  `xlsx:"13"`
	Lon            float64 `xlsx:"-" json:"-"`
	Status         string  `xlsx:"8"`
	RecommendedUse string  `xlsx:"9"`
	//BoolVal bool `xlsx:"4"`
	NameSlAlt string `xlsx:"21"`
	NameEn    string `xlsx:"32"`
	NameFr    string `xlsx:"33"`
	NameDe    string `xlsx:"34"`
	NameEs    string `xlsx:"35"`
	NameRu    string `xlsx:"36"`
	NameIt    string `xlsx:"37"`
	NameHr    string `xlsx:"38"`
	NameHu    string `xlsx:"39"`

	Etymology string `xlsx:"40"`
	Note      string `xlsx:"41"`
}
