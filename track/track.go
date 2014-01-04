package track

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Root struct {
	XMLName   xml.Name `xml:"tracks"`
	TrackList []Track  `xml:"track"`
}

type Track struct {
	Name    string   `xml:"name"`
	Artists []string `xml:"artist>name"`
	Album   string   `xml:"album>name"`
	Href    string   `xml:"href,attr"`
}

func Fetch(title, artist, album string) (Track, error) {
	return Track{}, nil
}

//fmt.Sprintf("http://ws.spotify.com/search/1/track?q=track:%s artist:%s album:%s", name, artist, album)
//fmt.Sprintf("http://ws.spotify.com/search/1/track?q=track:%s artist:%s", name, artist)

func ConstructSearchURL(title, artist, album string) string {
	if len(title) > 0 {
		url := fmt.Sprintf("http://ws.spotify.com/search/1/track?q=track:%s", title)

		if len(artist) > 0 {
			url += fmt.Sprintf(" artist:%s", artist)
		}

		url += fmt.Sprintf(" album:%s", album)

		return url
	}

	return ""
}

func FetchTracksXML(url string) ([]byte, error) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotModified) {
		return nil, fmt.Errorf("spotify: GET request in FetchTracksXML returned status %d rather than %d or %d", resp.StatusCode, http.StatusOK, http.StatusNotModified)
	}

	// Only returns error if bytes Buffer mecomes to large. How can this be tested?
	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}

func ExtractSingleTrackFromXML(xml_data []byte) (Track, error) {
	track_list, err := extractTracksFromXML(xml_data)

	if err != nil {
		return Track{}, err
	}

	return track_list[0], nil
}

func extractTracksFromXML(xml_data []byte) ([]Track, error) {
	var r Root

	err := xml.Unmarshal(xml_data, &r)

	if err != nil {
		return nil, errors.New("spotify: unable to unmarshal xml_data in extractTracksFromXML")
	}

	return r.TrackList, nil
}
