// Package track fetches tracks using the Spotify Web API
package track

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const trackSearchBaseUrl = "https://api.spotify.com/v1/search/?type=track&q="

type ErrorType int

const (
	ArgumentError ErrorType = iota
	UnexpectedError
	ExternalServiceError
	RateLimitError
)

// Track represent a Spotify track
type Track struct {
	Name    string
	Artists []string
	Album   string
	Uri     string
}

type Searcher struct {
	trackSearchBaseUrl string
}

type TrackError struct {
	Msg           string
	ErrorType     ErrorType
	OriginalError error
}

func (te TrackError) Error() string {
	msg := "github.com/joarleth/spotify/track: " + te.Msg

	if te.OriginalError != nil {
		msg += " Original error: " + te.OriginalError.Error()
	}

	return msg
}

// NewSearcher initializes a default searcher object
func NewSearcher() *Searcher {
	return &Searcher{
		trackSearchBaseUrl: trackSearchBaseUrl,
	}
}

// Find returns a track from Spotify matching title and at least one of artist and album.
// The data is fetched from Spotify's Web API. (https://developer.spotify.com/web-api/)
func (s Searcher) Find(title, artist, album string) (Track, error) {
	searchQueries, err := constructSearchQuery(title, artist, album)

	if err != nil {
		return Track{}, err
	}

	// TODO: Loop through search queries if more than one and no track found
	url := s.trackSearchBaseUrl + searchQueries[0] + "&limit=1"

	println(url)

	data, fetchError := fetchData(url)

	if fetchError != nil {
		return Track{}, fetchError
	}

	track, extractError := s.extractTrackFromJSON(data)

	if extractError != nil {
		return Track{}, extractError
	}

	return track, nil
}

/*

// GetClosestMatch returns a track from spotify matching title and at least one of artist and album.
// The data is fetched from Spotify's metadata API. (https://developer.spotify.com/technologies/web-api/)
// The first track containing the territory string of the Searcher instance, or "worldwide" in it's
// territories string
// is returned.
//
// Please beware of rate limits;
// "The rate limit is currently 10 request per second per ip. This may change."
func (s Searcher) FindClosestMatch(title, artist, album string) (Track, error) {
	search_queries, err := constructSearchQuery(title, artist, album)

	if err != nil {
		return Track{}, err
	}

	// TODO: Loop through search queries if more than one and no track found
	url := s.trackSearchBaseUrl + "/" + search_queries[0]

	xml_data, fetch_error := fetchTracksXML(url)

	if fetch_error != nil {
		return Track{}, fetch_error
	}

	track, extract_err := s.extractSingleTrackFromXML(xml_data)

	if extract_err != nil {
		return Track{}, extract_err
	}

	return track, nil
}

*/

func constructSearchQuery(title, artist, album string) ([]string, error) {
	title = strings.TrimSpace(title)
	artist = strings.TrimSpace(artist)
	album = strings.TrimSpace(album)

	if len(title) > 0 {
		// If both artist and album are supplied, return array of three
		// search queries so that these can be tried in order if no tracks
		// are returned.
		if len(artist) > 0 && len(album) > 0 {
			return []string{
				url.QueryEscape(constructSearchQueryFromTitleArtistAndAlbum(title, artist, album)),
				url.QueryEscape(constructSearchQueryFromTitleAndArtist(title, artist)),
				url.QueryEscape(constructSearchQueryFromTitleAndAlbum(title, album)),
			}, nil
		} else if len(artist) > 0 {
			return []string{
				url.QueryEscape(constructSearchQueryFromTitleAndArtist(title, artist)),
			}, nil
		} else if len(album) > 0 {
			return []string{
				url.QueryEscape(constructSearchQueryFromTitleAndAlbum(title, album)),
			}, nil
		}
	}

	return nil, TrackError{Msg: "A title and at least one of article and album must be passed as arguments.", ErrorType: ArgumentError}
}

// TODO Consider skipping the quotes, maybe
func constructSearchQueryFromTitleAndArtist(title, artist string) string {
	return fmt.Sprintf("track:\"%s\" artist:\"%s\"", title, artist)
}

func constructSearchQueryFromTitleAndAlbum(title, album string) string {
	return fmt.Sprintf("track:\"%s\" album:\"%s\"", title, album)
}

func constructSearchQueryFromTitleArtistAndAlbum(title, artist, album string) string {
	return fmt.Sprintf("track:\"%s\" artist:\"%s\" album:\"%s\"", title, artist, album)
}

func fetchData(url string) ([]byte, error) {
	resp, httpErr := http.Get(url)
	defer resp.Body.Close()

	if httpErr != nil {
		return []byte{}, TrackError{Msg: "Get request failed in fetchData.", ErrorType: UnexpectedError, OriginalError: httpErr}
	}

	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotModified) {
		if resp.StatusCode == http.StatusForbidden {
			return nil, TrackError{Msg: "Rate limit exceeded at Spotify Metadata API.", ErrorType: RateLimitError}
		}

		return nil, TrackError{Msg: fmt.Sprintf("GET request in fetchData returned status %d rather than %d or %d", resp.StatusCode, http.StatusOK, http.StatusNotModified), ErrorType: ExternalServiceError}
	}

	body, ioutilErr := ioutil.ReadAll(resp.Body)

	if ioutilErr != nil {
		return []byte{}, TrackError{Msg: "ioutil.ReadAll failed in fetchData.", ErrorType: UnexpectedError, OriginalError: ioutilErr}
	}

	return body, nil
}

func (s Searcher) extractTrackFromJSON(xml_data []byte) (Track, error) {
	trackCollection, err := extractTrackCollectionFromJSON(xml_data)

	if err != nil {
		return Track{}, err
	}

	if len(trackCollection.Tracks.Items) > 0 {
		trackItem := trackCollection.Tracks.Items[0]

		var artists []string

		for _, artist := range trackItem.Artists {
			artists = append(artists, artist.Name)
		}

		track := Track{
			Name:    trackItem.Name,
			Uri:     trackItem.Uri,
			Album:   trackItem.Album.Name,
			Artists: artists,
		}

		return track, nil
	}

	return Track{}, nil
}

// trackCollection, trackItem, item, album and artist are structs
// used for unmarsahlling json data from the Spotify API.
type trackCollection struct {
	Tracks trackItem
}
type trackItem struct {
	Href  string
	Items []item
}
type item struct {
	Uri     string
	Name    string
	Album   album
	Artists []artist
}
type album struct {
	Name string
}
type artist struct {
	Name string
}

func extractTrackCollectionFromJSON(jsonData []byte) (trackCollection, error) {
	var tc trackCollection
	err := json.Unmarshal(jsonData, &tc)

	if err != nil {
		return trackCollection{}, TrackError{Msg: "Unable to unmarshal jsonData in extractTrackCollectionFromJSON.", OriginalError: err, ErrorType: ExternalServiceError}
	}

	return tc, nil
}
