package track

import (
	"reflect"
	"testing"
)

func TestExtractSingleTrackFromXMLCorrectTrack(t *testing.T) {
	xml_data := getTracksFile(t)

	expected := Track{
		Name:    "True Affection",
		Artists: []string{"The Blow"},
		Album:   "Paper Television",
		Href:    "spotify:track:0tO8FKgGQzzuf8KGkHGeIw",
	}
	actual, _ := ExtractSingleTrackFromXML(xml_data)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Resulting track not matching expected.\nExpected: %v\nActual: %v", expected, actual)
	}
}

func TestExtractSingleTrackFromXMLReturnsError(t *testing.T) {
	xml_data := "<tracks>"

	_, err := ExtractSingleTrackFromXML([]byte(xml_data))

	if err == nil {
		t.Error("Expexted error caused by malformed xml.")
	}
}

func TestExtractSingleTrackFromXMLCorrectErrorMessage(t *testing.T) {
	xml_data := "<tracks>"

	_, err := ExtractSingleTrackFromXML([]byte(xml_data))

	expected := "spotify/track: unable to unmarshal xml_data in extractTracksFromXML"
	actual := err.Error()

	if expected != actual {
		t.Errorf("Unexpected error message.\nExpected: %v\nActual: %v", expected, actual)
	}
}
