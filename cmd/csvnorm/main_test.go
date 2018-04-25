package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

import (
	"github.com/pkg/errors"
)

func TestNormalizeString(t *testing.T) {
	cases := []struct{ in, want string }{
		{"hello " + string([]byte{0xff, 0xfe, 0xfd}), "hello " + string([]rune{'\uFFFD', '\uFFFD', '\uFFFD'})},
		{string([]byte{0xff, 0xfe, 0xaa}) + " world", string([]rune{'\uFFFD', '\uFFFD', '\uFFFD'}) + " world"},
	}
	for _, c := range cases {
		got := normalizeString(c.in)
		if got != c.want {
			t.Errorf("normalizeString(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestParseTimestamp(t *testing.T) {

	inputLocation, err := time.LoadLocation("US/Pacific")
	if err != nil {
		fmt.Println(errors.Wrap(err, "Could not load input timezone location with code US/Pacific"))
		os.Exit(1)
	}

	outputLocation, err := time.LoadLocation("US/Eastern")
	if err != nil {
		fmt.Println(errors.Wrap(err, "Could not load output timezone location with code US/Eastern"))
		os.Exit(1)
	}

	cases := []struct{ in, want string }{
		{"4/1/11 11:00:00 AM", "2011-04-01T14:00:00-04:00"},
		{"2/29/16 12:11:11 PM", "2016-02-29T15:11:11-05:00"},
	}
	for _, c := range cases {
		in, err := parseTimestamp(c.in, inputLocation)
		if err != nil {
			t.Errorf("parseTimestamp(%q) through an error %q", c.in, err)
		}
		got := in.In(outputLocation).Format(time.RFC3339)
		if got != c.want {
			t.Errorf("parseTimestamp(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestParseDuration(t *testing.T) {

	cases := []struct {
		in   string
		want float64
	}{
		{"1:23:32.123", float64(5012)},
		{"111:23:32.123", float64(401012)},
	}
	for _, c := range cases {
		got, err := parseDuration(c.in)
		if err != nil {
			t.Errorf("parseDuration(%q) through an error %q", c.in, err)
		}
		if got.Seconds() != c.want {
			t.Errorf("parseDuration(%q) == %q, want %f", c.in, got, c.want)
		}
	}
}
