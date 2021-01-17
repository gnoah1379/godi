package internal

import "fmt"

const (
	MAJOR         = 1
	MINOR         = 0
	PATCH         = 0
	STATE         = ALPHA
	STATE_VERSION = 1
)

const (
	ALPHA             = "a"
	BETA              = "b"
	RELEASE_CANDIDATE = "rc"
	RELEASE           = ""
)

func Version() (version string) {
	threePart := fmt.Sprintf("%d.%d.%d", MAJOR, MINOR, PATCH)
	statePart := fmt.Sprintf("%s.%d", STATE, STATE_VERSION)
	if STATE != RELEASE {
		version = threePart + "-" + statePart
	} else {
		version = threePart
	}
	return "godi/" + version
}
