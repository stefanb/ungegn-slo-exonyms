package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	geojson "github.com/paulmach/go.geojson"
)

func setWikidataOverpassLink(f *geojson.Feature, key string, wikidataID string) {
	if !looksLikeWikidataID(wikidataID) {
		return
	}

	wikidataOverpassLink := getWikidataOverpassLink(wikidataID)

	if wikidataOverpassLink != "" {
		f.SetProperty(key, wikidataOverpassLink)
	}
}

func setWikidataLink(f *geojson.Feature, key string, wikidataID string) {
	if !looksLikeWikidataID(wikidataID) {
		return
	}

	f.SetProperty(key, "https://www.wikidata.org/wiki/"+wikidataID)
}

var wikidataRegex = regexp.MustCompile("Q[1-9][0-9]*")

func looksLikeWikidataID(value string) bool {
	return strings.HasPrefix(value, "Q") && wikidataRegex.MatchString(value)
}

// Generates a link to OSM element via wikidata, as per https://www.wikidata.org/wiki/Template:Overpasslink
func getWikidataOverpassLink(wikidataID string) string {

	//https://overpass-api.de/api/interpreter?data=%5Bout%3Acustom%5D%3Bnode%5Bwikidata%3DQ641%5D%3Bif%28count%28nodes%29%3D%3D0%29%7Bway%5Bwikidata%3DQ641%5D%3B%7D%3Bif%28count%28ways%29%3D%3D0%29%7Brel%5Bwikidata%3DQ641%5D%3B%7D%3Bout%201%3B

	const linkPrefix = "https://overpass-api.de/api/interpreter?data="
	const dataParamTemplate = "[out:custom];node[wikidata=%s];if(count(nodes)==0){way[wikidata=%s];};if(count(ways)==0){rel[wikidata=%s];};out 1;"

	return linkPrefix + url.QueryEscape(fmt.Sprintf(dataParamTemplate, wikidataID, wikidataID, wikidataID))
}
