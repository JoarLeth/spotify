package track

import (
	"reflect"
	"testing"
)

func TestFind(t *testing.T) {
	s := NewSearcher()

	expected := Track{Name: "Labyrinth",
		Uri:     "spotify:track:7f7y9A3Spuus0SBsuDMdMa",
		Artists: []string{"Bella Hardy"},
		Album:   "Songs Lost & Stolen",
	}
	actual, _ := s.Find("Labyrinth", "Bella Hardy", "")

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Actual track not matching expected. \nExpected: %v\nActual:   %v\n", expected, actual)
	}
}

func TestExtractTracksFromJSON(t *testing.T) {
	data := getTextFileData(t, "test_data/tracks.json")

	trackCollection, _ := extractTrackCollectionFromJSON(data)

	expectedLength := 20
	actualLength := len(trackCollection.Tracks.Items)

	if expectedLength != actualLength {
		t.Errorf("Unexpected number of tracks. Expected: %v, got: %v", expectedLength, actualLength)
	}

	expectedFirstTrack := item{Uri: "spotify:track:4ry6oqlwdsooYtniYJFkt5",
		Name:    "Human Behaviour",
		Artists: []artist{artist{Name: "Björk"}},
		Album:   album{Name: "Debut (Ecopac)"},
	}

	actualFirstTrack := trackCollection.Tracks.Items[0]

	if !reflect.DeepEqual(expectedFirstTrack, actualFirstTrack) {
		t.Errorf("Actual first track not matching expected. \nExpected: %v\nActual:   %v\n", expectedFirstTrack, actualFirstTrack)
	}
}

/*
func TestExtractTracksFromJSONFirstTrackCorrect(t *testing.T) {
	xml_data := getTextFileData(t, "tracks.xml")

	track_list, _ := extractTracksFromXML(xml_data)

	expected := Track{
		Name:        "True Affection",
		Artists:     []string{"The Blow"},
		Album:       "Paper Television",
		Href:        "spotify:track:0tO8FKgGQzzuf8KGkHGeIw",
		Territories: "US",
	}
	actual := track_list[0]

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Resulting track not matching expected.\nExpected: %v\nActual: %v", expected, actual)
	}
}

func TestExtractTracksFromJSONTrackWithMultipleArtists(t *testing.T) {
	xml_data := getTextFileData(t, "tracks.xml")

	track_list, _ := extractTracksFromXML(xml_data)

	expectedFirst := Track{
		Name:        "Affection - True 2 Life Remix",
		Artists:     []string{"Pat Bedeau", "Steve Gurley", "Shishani"},
		Album:       "Affection",
		Href:        "spotify:track:1tJD2Pk0o3TBcMjDSOxMKp",
		Territories: "AD AR AT AU BE BG BO BR CA CH CL CO CR CY CZ DE DK DO EC EE ES FI FR GB GR GT HK HN HU IE IS IT LI LT LU LV MC MT MX MY NI NL NO NZ PA PE PH PL PT PY RO SE SG SI SK SV TR TW US UY",
	}

	actualFirst := track_list[8]

	if !reflect.DeepEqual(actualFirst, actualFirst) {
		t.Errorf("Resulting track not matching expected.\nExpected: %#v\nActual: %#v", expectedFirst, actualFirst)
	}

	expectedLast := Track{
		Name:        "Affection",
		Artists:     []string{"The True Bypass"},
		Album:       "No Hero Sound",
		Href:        "spotify:track:5CG3XEWCSetrQSKgqhN6NR",
		Territories: "AD AR AT AU BE BG BO BR CA CH CL CO CR CY CZ DE DK DO EC EE ES FI FR GB GR GT HK HN HU IE IS IT LI LT LU LV MC MT MX MY NI NL NO NZ PA PE PH PL PT PY RO SE SG SI SK SV TR TW US UY",
	}

	actualLast := track_list[29]

	if !reflect.DeepEqual(expectedLast, actualLast) {
		t.Errorf("Resulting track not matching expected.\nExpected: %#v\nActual: %#v", expectedLast, actualLast)
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

	expected := errorPrefix + "Unable to unmarshal xml_data in extractTracksFromXML. Original error: EOF"
	actual := err.Error()

	if expected != actual {
		t.Errorf("Unexpected error message.\nExpected: %v\nActual: %v", expected, actual)
	}
}

func TestExtractTracksFromJSONMalformadXMLErrorIsNil(t *testing.T) {
	xml_data := getTextFileData(t, "tracks.xml")

	_, err := extractTracksFromXML(xml_data)

	if err != nil {
		t.Errorf("Expexted error to be nil. Got: %v", err.Error())
	}
}

func TestExtractTracksFromXMLWithNoTracksReturnsNil(t *testing.T) {
	xml_data := getTextFileData(t, "no_tracks.xml")

	track_list, _ := extractTracksFromXML(xml_data)

	if track_list != nil {
		t.Errorf("Expected track_list to be nil.\nActual value: %v", track_list)
	}
}

func getTextFileData(t *testing.T, filename string) []byte {
	file, open_file_error := os.Open(filename)
	defer file.Close()

	if open_file_error != nil {
		t.Fatalf("Failed to open file. Error: %v", open_file_error.Error())
	}

	data, read_file_error := ioutil.ReadAll(file)

	if read_file_error != nil {
		t.Fatalf("Failed to read file. Error: %v", read_file_error.Error())
	}

	return data
}
*/
