// Package track fetches tracks using Spotify's metadata API
package track

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const trackSearchBaseUrl = "http://ws.spotify.com/search/1/track?q="
const preferredTerritory = "se"

type Root struct {
	XMLName   xml.Name `xml:"tracks"`
	TrackList []Track  `xml:"track"`
}

type Track struct {
	Name        string   `xml:"name"`
	Artists     []string `xml:"artist>name"`
	Album       string   `xml:"album>name"`
	Href        string   `xml:"href,attr"`
	Territories string   `xml:"album>availability>territories"`
}

type Searcher struct {
	territory             string
	track_search_base_url string
}

type TrackError struct {
	Msg           string
	ErrorType     ErrorType
	OriginalError error
}

type ErrorType int

const (
	ArgumentError        ErrorType = 3
	UnexpectedError      ErrorType = 4
	ExternalServiceError ErrorType = 5
	RateLimitError       ErrorType = 6
)

func (te TrackError) Error() string {
	msg := "github.com/joarleth/spotify/track: " + te.Msg

	if te.OriginalError != nil {
		msg += " Original error: " + te.OriginalError.Error()
	}

	return msg
}

func NewSearcher(territory string) *Searcher {
	return &Searcher{
		territory:             territory,
		track_search_base_url: trackSearchBaseUrl,
	}
}

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
	url := s.track_search_base_url + "/" + search_queries[0]

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

func constructSearchQueryFromTitleAndArtist(title, artist string) string {
	return fmt.Sprintf("track:\"%s\" artist:\"%s\"", title, artist)
}

func constructSearchQueryFromTitleAndAlbum(title, album string) string {
	return fmt.Sprintf("track:\"%s\" album:\"%s\"", title, album)
}

func constructSearchQueryFromTitleArtistAndAlbum(title, artist, album string) string {
	return fmt.Sprintf("track:\"%s\" artist:\"%s\" album:\"%s\"", title, artist, album)
}

func fetchTracksXML(url string) ([]byte, error) {
	// From http.Get:
	// An error is returned if the Client's CheckRedirect function fails
	// or if there was an HTTP protocol error. A non-2xx response doesn't
	// cause an error.
	resp, httpErr := http.Get(url)
	defer resp.Body.Close()

	if httpErr != nil {
		return []byte{}, TrackError{Msg: "Get request failed in fetchTracksXML.", ErrorType: UnexpectedError, OriginalError: httpErr}
	}

	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotModified) {
		if resp.StatusCode == http.StatusForbidden {
			return nil, TrackError{Msg: "Rate limit exceeded at Spotify Metadata API.", ErrorType: RateLimitError}
		}

		return nil, TrackError{Msg: fmt.Sprintf("GET request in FetchTracksXML returned status %d rather than %d or %d", resp.StatusCode, http.StatusOK, http.StatusNotModified), ErrorType: ExternalServiceError}
	}

	// Only returns error if bytes Buffer becomes too large. How can this be tested?
	body, ioutilErr := ioutil.ReadAll(resp.Body)

	if ioutilErr != nil {
		return []byte{}, TrackError{Msg: "ioutil.ReadAll failed in fetchTracksXML.", ErrorType: UnexpectedError, OriginalError: ioutilErr}
	}

	return body, nil
}

func (s Searcher) extractSingleTrackFromXML(xml_data []byte) (Track, error) {
	track_list, err := extractTracksFromXML(xml_data)

	if err != nil {
		return Track{}, err
	}

	// Return the first track where territories contains the territory string or "worldwide"
	for _, track := range track_list {
		// TODO: Change to use Compile and add error handling since this value is no longer hard coded
		// Or better still, use MustCompile at instantiation.
		re := regexp.MustCompile("(?i)(" + s.territory + "|worldwide)")

		if matched := re.MatchString(track.Territories); matched {
			return track, nil
		}
	}

	return Track{}, nil
}

func extractTracksFromXML(xml_data []byte) ([]Track, error) {
	var r Root

	err := xml.Unmarshal(xml_data, &r)

	if err != nil {
		return nil, TrackError{Msg: "Unable to unmarshal xml_data in extractTracksFromXML.", OriginalError: err, ErrorType: ExternalServiceError}
	}

	return r.TrackList, nil
}
