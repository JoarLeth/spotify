package track

import (
	"reflect"
	"testing"
)

func TestNewSearcher(t *testing.T) {
	expected := &Searcher{
		preferred_territory:   "EE",
		track_search_base_url: trackSearchBaseUrl,
	}

	actual := NewSearcher("EE")

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Searcher not matching expected value.\nExpected: %#v\nActual: %#v", expected, actual)
	}
}
