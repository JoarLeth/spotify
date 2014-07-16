package track

import (
	"io/ioutil"
	"os"
	"testing"
)

func getTextFileData(t *testing.T, filename string) []byte {
	file, open_file_error := os.Open(filename)
	defer file.Close()

	if open_file_error != nil {
		t.Fatalf("Failed to open file. Error: %v", open_file_error.Error())
	}

	data, read_file_error := ioutil.ReadAll(file)

	if read_file_error != nil {
		t.Fatalf("Failed to read file. Error: %v", read_file_error.Error())
	}

	return data
}
