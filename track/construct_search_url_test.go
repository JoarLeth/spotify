package track

import (
	"fmt"
	"testing"
)

func TestConstructSearchURLAllArgumentsIncluded(t *testing.T) {
	title := "asdf"
	artist := "qwer"
	album := "ty"

	expected := fmt.Sprintf("http://ws.spotify.com/search/1/track?q=track:%s artist:%s album:%s", title, artist, album)
	actual := ConstructSearchURL(title, artist, album)

	if expected != actual {
		t.Errorf("Incorrect spotify search url.\nExpected: %s\nActual: %s", expected, actual)
	}
}

func TestConstructSearchURLEmptyTitleReturnsEmptyString(t *testing.T) {
	title := ""
	artist := "qwer"
	album := "ty"

	expected := ""
	actual := ConstructSearchURL(title, artist, album)

	if expected != actual {
		t.Errorf("Empty title should return empty string. Got: %s", actual)
	}
}

func TestConstructSearchURLEmptyArtistOmitsArtistFromSearchString(t *testing.T) {
	title := "asdf"
	artist := ""
	album := "ty"

	expected := fmt.Sprintf("http://ws.spotify.com/search/1/track?q=track:%s album:%s", title, album)
	actual := ConstructSearchURL(title, artist, album)

	if expected != actual {
		t.Errorf("Incorrect spotify search url.\nExpected: %s\nActual: %s", expected, actual)
	}
}
