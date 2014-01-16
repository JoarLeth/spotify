package track

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const (
	errorPrefix = "github.com/joarleth/spotify/track: "
)

func TestFetchTracksXMLReturnsErrorOnInternalServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", http.StatusInternalServerError)
	}))

	_, err := fetchTracksXML(ts.URL)

	if err == nil {
		t.Error("Expected FetchTracksXML to return an error.")
	}
}

func TestFetchTracksXMLReturnsCorrectErrorMessageOnNonOkStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", http.StatusInternalServerError)
	}))

	_, err := fetchTracksXML(ts.URL)

	expected := fmt.Sprintf("%sGET request in FetchTracksXML returned status %d rather than %d or %d", errorPrefix, http.StatusInternalServerError, http.StatusOK, http.StatusNotModified)
	actual := err.Error()

	if expected != actual {
		t.Errorf("Unecpected error message.\nExpected: %v\nActual: %v", expected, actual)
	}

	terr, isTrackError := err.(TrackError)

	if isTrackError == false {
		t.Errorf("Expected error to be of type TrackError")
	}

	if terr.ErrorType != ExternalServiceError {
		t.Errorf("Expected ErrorType to be ExternalServceError.")
	}
}

func TestFetchTracksXMLReturnsXMLData(t *testing.T) {
	xml_data := getTextFileData(t, "tracks.xml")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Write(xml_data)
	}))

	expected := xml_data
	actual, _ := fetchTracksXML(ts.URL)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Unexpected xml data.\nExpected: %v\nActual: %v", string(expected), string(actual))
	}
}

func TestFetchTracksXML403ReturnsTrackErrorRateLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", http.StatusForbidden)
	}))

	_, err := fetchTracksXML(ts.URL)

	expectedMessage := errorPrefix + "Rate limit exceeded at Spotify Metadata API."
	acutalMessage := err.Error()

	if expectedMessage != acutalMessage {
		t.Errorf("Unexpected error message.\nExpected: %v\nActual: %v", expectedMessage, acutalMessage)
	}

	terr, isTrackError := err.(TrackError)

	if isTrackError == false {
		t.Errorf("Expected error to be of type TrackError")
	}

	if terr.ErrorType != RateLimitError {
		t.Errorf("Expected ErrorType to be RateLimitError.")
	}
}
