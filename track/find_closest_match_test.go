package track

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFindClosestMatch(t *testing.T) {
	xml_data := getTextFileData(t, "tracks2.xml")

	mockserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Write(xml_data)
	}))

	s := new_mock_searcher("SE", mockserver.URL)

	expected := Track{
		Name:        "Uncover",
		Artists:     []string{"Zara Larsson"},
		Album:       "Introducing",
		Href:        "spotify:track:131l5GkXPIk81bxihGypPt",
		Territories: "SE",
	}
	actual, _ := s.FindClosestMatch("Uncover", "Zara Larsson", "")

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Resulting track not matching expected.\nExpected: %#v\nActual: %#v", expected, actual)
	}
}

func TestFindClosestMatchDifferentTerritory(t *testing.T) {
	xml_data := getTextFileData(t, "tracks_bjork.xml")

	mockserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Write(xml_data)
	}))

	s := new_mock_searcher("US", mockserver.URL)

	expected := Track{
		Name:        "Human Behaviour",
		Artists:     []string{"Björk"},
		Album:       "Debut",
		Href:        "spotify:track:5OnyZ56HLhrWOXdzeETqLk",
		Territories: "CA US",
	}
	actual, _ := s.FindClosestMatch("Uncover", "Zara Larsson", "")

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Resulting track not matching expected.\nExpected: %#v\nActual: %#v", expected, actual)
	}
}

func new_mock_searcher(territory, search_url string) *Searcher {
	return &Searcher{
		territory:             territory,
		track_search_base_url: search_url,
	}
}
