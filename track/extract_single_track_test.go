package track

import (
	"reflect"
	"testing"
)

func TestExtractSingleTrackFromXMLCorrectTrack(t *testing.T) {
	xml_data := getTextFileData(t, "tracks.xml")

	s := NewSearcher("SE")

	expected := Track{
		Name:        "True Affection",
		Artists:     []string{"The Blow"},
		Album:       "Paper Television",
		Href:        "spotify:track:1js3QhuQP3dwk4l2DrPXDC",
		Territories: "AD AT BE BG CH CY CZ DE DK EE ES FI FR GB GR HU IE IS IT LI LT LU LV MC MT NL NO PL PT RO SE SI SK TR",
	}
	actual, _ := s.extractSingleTrackFromXML(xml_data)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Resulting track not matching expected.\nExpected: %v\nActual: %v", expected, actual)
	}
}

func TestExtractSingleTrackFromXMLNoTracksReturnsNilTrack(t *testing.T) {
	xml_data := getTextFileData(t, "no_tracks.xml")

	s := NewSearcher("SE")

	expected := Track{}
	actual, _ := s.extractSingleTrackFromXML(xml_data)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Resulting track not matching expected.\nExpected: %v\nActual: %v", expected, actual)
	}
}

func TestExtractSingleTrackFromXMLReturnsError(t *testing.T) {
	xml_data := "<tracks>"

	s := NewSearcher("SE")

	_, err := s.extractSingleTrackFromXML([]byte(xml_data))

	if err == nil {
		t.Error("Expexted error caused by malformed xml.")
	}
}

func TestExtractSingleTrackFromXMLCorrectErrorMessage(t *testing.T) {
	xml_data := "<tracks>"

	s := NewSearcher("SE")

	_, err := s.extractSingleTrackFromXML([]byte(xml_data))

	expected := "spotify/track: unable to unmarshal xml_data in extractTracksFromXML"
	actual := err.Error()

	if expected != actual {
		t.Errorf("Unexpected error message.\nExpected: %v\nActual: %v", expected, actual)
	}
}
