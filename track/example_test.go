package track

import (
	"fmt"
)

// Since this example runs code retrieving data from an external API
// The output may well change in the future.
func ExampleSearcher() {
	s := NewSearcher("SE")

	track, err := s.FindClosestMatch("human behaviour", "björk", "")

	if err != nil {
		fmt.Printf("An error occurred. Error: %s", err.Error())
	}

	fmt.Println(track.Href)
	// Output: spotify:track:4ry6oqlwdsooYtniYJFkt5
}
