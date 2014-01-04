package track

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestExtractTracksFromJSONCorrectNumberOfResults(t *testing.T) {
	xml_data := getTracksFile(t)

	track_list, _ := extractTracksFromXML(xml_data)

	expected := 30
	actual := len(track_list)

	if expected != actual {
		t.Errorf("Unexpected number of tracks. Expected: %v, got: %v", expected, actual)
	}
}

func TestExtractTracksFromJSONFirstTrackCorrect(t *testing.T) {
	xml_data := getTracksFile(t)

	track_list, _ := extractTracksFromXML(xml_data)

	expected := Track{
		Name:    "True Affection",
		Artists: []string{"The Blow"},
		Album:   "Paper Television",
		Href:    "spotify:track:0tO8FKgGQzzuf8KGkHGeIw",
	}
	actual := track_list[0]

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Resulting track not matching expected.\nExpected: %v\nActual: %v", expected, actual)
	}
}

func TestExtractTracksFromJSONTrackWithMultipleArtists(t *testing.T) {
	xml_data := getTracksFile(t)

	track_list, _ := extractTracksFromXML(xml_data)

	expected := Track{
		Name:    "Affection - True 2 Life Remix",
		Artists: []string{"Pat Bedeau", "Steve Gurley", "Shishani"},
		Album:   "Affection",
		Href:    "spotify:track:1tJD2Pk0o3TBcMjDSOxMKp",
	}
	actual := track_list[8]

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Resulting track not matching expected.\nExpected: %v\nActual: %v", expected, actual)
	}
}

func TestExtractTracksFromJSONMalformadXMLReturnsError(t *testing.T) {
	xml_data := "<tracks>"

	_, err := extractTracksFromXML([]byte(xml_data))

	if err == nil {
		t.Error("Expexted error caused by malformed xml.")
	}
}

func TestExtractTracksFromJSONMalformadXMLCorrectErrorMessage(t *testing.T) {
	xml_data := "<tracks>"

	_, err := extractTracksFromXML([]byte(xml_data))

	expected := "spotify/track: unable to unmarshal xml_data in extractTracksFromXML"
	actual := err.Error()

	if expected != actual {
		t.Errorf("Unexpected error message.\nExpected: %v\nActual: %v", expected, actual)
	}
}

func TestExtractTracksFromJSONMalformadXMLErrorIsNil(t *testing.T) {
	xml_data := getTracksFile(t)

	_, err := extractTracksFromXML(xml_data)

	if err != nil {
		t.Errorf("Expexted error to be nil. Got: %v", err.Error())
	}
}

func getTracksFile(t *testing.T) []byte {
	xml_file, open_file_error := os.Open("tracks.xml")
	defer xml_file.Close()

	if open_file_error != nil {
		t.Fatalf("Failed to open file. Error: %v", open_file_error.Error())
	}

	xml_data, read_file_error := ioutil.ReadAll(xml_file)

	if read_file_error != nil {
		t.Fatalf("Failed to read file. Error: %v", read_file_error.Error())
	}

	return xml_data
}
