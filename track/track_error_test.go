package track

import (
	"testing"
)

type test_error struct {
	msg string
}

func (te test_error) Error() string {
	return "testing: " + te.msg
}

func TestTrackError(t *testing.T) {
	msg := "Testing TrackError."
	var err error
	err = TrackError{Msg: msg, ErrorType: UnexpectedError}

	expectedMessage := errorPrefix + msg
	actualMessage := err.Error()

	if expectedMessage != actualMessage {
		t.Errorf("Unexpected error message.\nExpexted: %s\nActual: %s", expectedMessage, actualMessage)
	}

	terr, isTrackError := err.(TrackError)

	if !isTrackError {
		t.Error("Expected error to be of type TrackError.")
	}

	if terr.ErrorType != UnexpectedError {
		t.Error("Expected ErrorCode to be UnexpectedError.")
	}
}

func TestTrackErrorOriginalError(t *testing.T) {
	msg := "Testing TrackError."
	var err error
	oerr := test_error{msg: "Testing."}

	err = TrackError{Msg: msg, OriginalError: oerr, ErrorType: UnexpectedError}

	expectedMessage := errorPrefix + msg + " Original error: testing: Testing."
	actualMessage := err.Error()

	if expectedMessage != actualMessage {
		t.Errorf("Unexpected error message.\nExpexted: %s\nActual: %s", expectedMessage, actualMessage)
	}
}
