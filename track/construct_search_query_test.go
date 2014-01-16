package track

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

func TestConstructSearchURLAllArgumentsIncluded(t *testing.T) {
	title := "asdf"
	artist := "qwer"
	album := "ty"

	expected := []string{
		url.QueryEscape(fmt.Sprintf("track:\"%s\" artist:\"%s\" album:\"%s\"", title, artist, album)),
		url.QueryEscape(fmt.Sprintf("track:\"%s\" artist:\"%s\"", title, artist)),
		url.QueryEscape(fmt.Sprintf("track:\"%s\" album:\"%s\"", title, album)),
	}
	actual, err := constructSearchQuery(title, artist, album)

	if err != nil {
		t.Errorf("Expected error to be nil. Got: %s", err.Error())
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Incorrect spotify search query.\nExpected: %s\nActual: %s", expected, actual)
	}
}

func TestConstructSearchURLEmptyTitleReturnsError(t *testing.T) {
	title := ""
	artist := "qwer"
	album := "ty"

	_, err := constructSearchQuery(title, artist, album)

	if err == nil {
		t.Error("Passing empty title should return error.")
	}
}

func TestConstructSearchURLEmptyTitleCheckErrorMessage(t *testing.T) {
	title := ""
	artist := "qwer"
	album := "ty"

	_, err := constructSearchQuery(title, artist, album)

	expected := "github.com/joarleth/spotify/track: A title and at least one of article and album must be passed as arguments."
	actual := err.Error()

	if expected != actual {
		t.Errorf("Unecpected error message.\nExpected: %v\nActual: %v", expected, actual)
	}
}

func TestConstructSearchURLEmptyTitleReturnsTrackErrorArgumentError(t *testing.T) {
	title := ""
	artist := "qwer"
	album := "ty"

	_, err := constructSearchQuery(title, artist, album)

	terr, isTrackError := err.(TrackError)

	if isTrackError == false {
		t.Errorf("Expected error to be of type TrackError")
	}

	if terr.ErrorType != ArgumentError {
		t.Errorf("Expected ErrorType to be ArgumentError.")
	}
}

func TestConstructSearchURLEmptyTitleReturnsNil(t *testing.T) {
	title := ""
	artist := "qwer"
	album := "ty"

	actual, _ := constructSearchQuery(title, artist, album)

	if actual != nil {
		t.Errorf("Expected constructSearchQuery to return nil.\nActual return value: %v", actual)
	}
}

func TestConstructSearchURLOnlyTitleNotEmptyReturnsError(t *testing.T) {
	title := "asdf"
	artist := ""
	album := ""

	_, err := constructSearchQuery(title, artist, album)

	if err == nil {
		t.Error("Passing empty artist and album should return error.")
	}
}

func TestConstructSearchURLOnlyTitleCheckErrorMessage(t *testing.T) {
	title := "asdf"
	artist := ""
	album := ""

	_, err := constructSearchQuery(title, artist, album)

	expected := "github.com/joarleth/spotify/track: A title and at least one of article and album must be passed as arguments."
	actual := err.Error()

	if expected != actual {
		t.Errorf("Unecpected error message.\nExpected: %v\nActual: %v", expected, actual)
	}
}

func TestConstructSearchURLOnlyTitleReturnsNil(t *testing.T) {
	title := "asdf"
	artist := ""
	album := ""

	actual, _ := constructSearchQuery(title, artist, album)

	if actual != nil {
		t.Errorf("Expected constructSearchQuery to return nil.\nActual return value: %v", actual)
	}
}

func TestConstructSearchURLEmptyArtistOmitsArtistFromSearchString(t *testing.T) {
	title := "asdf"
	artist := ""
	album := "ty"

	expected := []string{url.QueryEscape(fmt.Sprintf("track:\"%s\" album:\"%s\"", title, album))}
	actual, err := constructSearchQuery(title, artist, album)

	if err != nil {
		t.Errorf("Expected error to be nil. Got: %s", err.Error())
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Incorrect spotify search query.\nExpected: %s\nActual: %s", expected, actual)
	}
}

func TestConstructSearchURLEmptyAlbumOmitsAlbumFromSearchString(t *testing.T) {
	title := "asdf"
	artist := "qwer"
	album := ""

	expected := []string{url.QueryEscape(fmt.Sprintf("track:\"%s\" artist:\"%s\"", title, artist))}
	actual, err := constructSearchQuery(title, artist, album)

	if err != nil {
		t.Errorf("Expected error to be nil. Got: %s", err.Error())
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Incorrect spotify search query.\nExpected: %s\nActual: %s", expected, actual)
	}
}

func TestConstructSearchURLStringsAreTrimmed(t *testing.T) {
	title := " asdf  "
	artist := "		 "
	album := "  ty"

	expected := []string{url.QueryEscape(fmt.Sprintf("track:\"%s\" album:\"%s\"", "asdf", "ty"))}
	actual, err := constructSearchQuery(title, artist, album)

	if err != nil {
		t.Errorf("Expected error to be nil. Got: %s", err.Error())
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Incorrect spotify search query.\nExpected: %s\nActual: %s", expected, actual)
	}
}
