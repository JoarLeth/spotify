package track

import (
	"fmt"
)

// Since this example runs code retrieving data from an external API
// The output may well change in the future.
func ExampleSearcher() {
	s := NewSearcher()

	track, err := s.Find("lazarus", "david byrne", "")

	if err != nil {
		fmt.Printf("An error occurred. Error: %s", err.Error())
	}

	fmt.Println(track.Uri)
	// Output: spotify:track:2NhEuDWWEeILAScaN2iPF4
}
