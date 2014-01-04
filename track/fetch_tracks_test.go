package track

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFetchTracksXMLReturnsErrorOnInternalServerError(t *testing.T) {
	mockserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", http.StatusInternalServerError)
	}))

	_, err := FetchTracksXML(mockserver.URL)

	if err == nil {
		t.Error("Expected FetchTracksXML to return an error.")
	}
}

func TestFetchTracksXMLReturnsCorrectErrorMessage(t *testing.T) {
	mockserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", http.StatusInternalServerError)
	}))

	_, err := FetchTracksXML(mockserver.URL)

	expected := fmt.Sprintf("spotify/track: GET request in FetchTracksXML returned status %d rather than %d or %d", http.StatusInternalServerError, http.StatusOK, http.StatusNotModified)
	actual := err.Error()

	if expected != actual {
		t.Errorf("Unecpected error message.\nExpected: %v\nActual: %v", expected, actual)
	}
}

func TestFetchTracksXMLReturnsXMLData(t *testing.T) {
	xml_data := getTracksFile(t)

	mockserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Write(xml_data)
	}))

	expected := xml_data
	actual, _ := FetchTracksXML(mockserver.URL)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Unexpected xml data.\nExpected: %v\nActual: %v", string(expected), string(actual))
	}
}
