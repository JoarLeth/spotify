// Package track fetches tracks using Spotify's metadata API
package track

import (
	"encoding/xml"
	"errors"
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
	preferred_territory   string
	track_search_base_url string
}

func NewSearcher(territory string) *Searcher {
	return &Searcher{
		preferred_territory:   territory,
		track_search_base_url: trackSearchBaseUrl,
	}
}

// GetClosestMatch returns a track from spotify matching title and at least one of artist and album.
// The data is fetched from Spotify's metadata API. (https://developer.spotify.com/technologies/web-api/)
// The first track containing the preferredTerritory constant, or "worldwide" in the territories string
// is returned.
//
// Please beware of rate limits;
// "The rate limit is currently 10 request per second per ip. This may change."
func (s *Searcher) FindClosestMatch(title, artist, album string) (Track, error) {
	search_queries, _ := constructSearchQuery(title, artist, album)

	url := s.track_search_base_url + "/" + search_queries[0]

	xml_data, _ := fetchTracksXML(url)

	track, _ := extractSingleTrackFromXML(xml_data)

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

	return nil, errors.New("spotify/track: A title and at least one of article and album must be supplied to constructSearchQuery.")
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
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotModified) {
		return nil, fmt.Errorf("spotify/track: GET request in FetchTracksXML returned status %d rather than %d or %d", resp.StatusCode, http.StatusOK, http.StatusNotModified)
	}

	// Only returns error if bytes Buffer becomes too large. How can this be tested?
	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}

func extractSingleTrackFromXML(xml_data []byte) (Track, error) {
	track_list, err := extractTracksFromXML(xml_data)

	if err != nil {
		return Track{}, err
	}

	// Return the first track where territories contains the preferredTerritory constant or "worldwide"
	for _, track := range track_list {
		re := regexp.MustCompile("(?i)(" + preferredTerritory + "|worldwide)")

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
		return nil, errors.New("spotify/track: unable to unmarshal xml_data in extractTracksFromXML")
	}

	return r.TrackList, nil
}
