package logspy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

var (
	sink bytes.Buffer
)

func Sink() *bytes.Buffer {
	return &sink
}

func Contains(ss string) bool {
	return strings.Contains(sink.String(), ss)
}

func ContainsJsonMember(f string) bool {
	return strings.Contains(String(), fmt.Sprintf("%q:", f))
}

// NumEntries() returns the number of non-empty entries in the log
func NumEntries() int {
	n := 0
	for _, e := range strings.Split(String(), "\n") {
		if len(strings.TrimSpace(e)) > 0 {
			n++
		}
	}
	return n
}

// NumJsonEntries() returns the number of Json objects in the log
func NumJsonEntries() (int, error) {
	n := 0
	r := strings.NewReader(sink.String())
	d := json.NewDecoder(r)
	for {
		var o interface{}
		if err := d.Decode(&o); err == io.EOF {
			break
		} else if err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

// Reset() clears captured logs, preparing the sink to capture new logs.
func Reset() {
	sink = bytes.Buffer{}
}

// String() returns a string containing the contents of the
// logs captured since the most recent Reset().
func String() string {
	return sink.String()
}
