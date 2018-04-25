package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

import (
	"github.com/pkg/errors"
)

var CSVNORM_VERSION = "0.0.1"
var CSVNORM_USAGE = "Usage: csvnorm [-input_timezone INPUT_TIMEZONE] [-output_timezone OUTPUT_TIMEZONE] [-help] [-version]"

func parseTimestamp(str string, loc *time.Location) (time.Time, error) {
	var timestamp time.Time
	var err error
	formats := []string{
		"1/2/06 15:04:05 PM",
		"1/2/06 15:04:05 PM (MST)",
	}
	for _, f := range formats {
		timestamp, err = time.ParseInLocation(f, str, loc)
		if err != nil {
			errors.Wrap(err, "Error parsing timestamp with format "+f)
			continue
		}
		return timestamp, nil
	}
	return time.Time{}, err
}

func parseDuration(s string) (time.Duration, error) {
	var v time.Duration

	if len(s) == 0 {
		return v, errors.New("duration string is blank")
	}

	parts := strings.Split(s, ":")

	hours_int, err := strconv.Atoi(parts[0])
	if err != nil {
		return v, errors.Wrap(err, "error parsing duration "+s)
	}

	minutes_int := 0
	if len(parts) > 1 {
		minutes_int, err = strconv.Atoi(parts[1])
		if err != nil {
			return v, errors.Wrap(err, "error parsing duration "+s)
		}
	}

	seconds_float := 0.0
	if len(parts) > 2 {
		seconds_float, err = strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return v, errors.Wrap(err, "error parsing duration "+s)
		}
	}

	v = time.Hour*time.Duration(hours_int) + time.Minute*time.Duration(minutes_int) + time.Second*time.Duration(seconds_float)

	return v, nil
}

func normalizeString(s string) string {
	runes := make([]rune, len(s))
	for i, v := range s {
		if utf8.ValidRune(v) {
			runes[i] = v
		} else {
			runes[i] = '\uFFFD'
		}
	}
	return string(runes)
}

func main() {

	var defaultInputTimezoneCode string
	var outputTimezoneCode string

	var version bool
	var help bool

	flag.StringVar(&defaultInputTimezoneCode, "input_timezone", "US/Pacific", "Default timezone for timestamp input.  Matches names in the IANA Time Zone database.")
	flag.StringVar(&outputTimezoneCode, "output_timezone", "US/Eastern", "Timezone for timestamp output.  Matches names in the IANA Time Zone database.")
	flag.BoolVar(&version, "version", false, "Prints version to stdout")
	flag.BoolVar(&help, "help", false, "Print help")

	flag.Parse()

	if help || (len(os.Args) == 2 && os.Args[1] == "help") {
		fmt.Println(CSVNORM_USAGE)
		flag.PrintDefaults()
		os.Exit(0)
	} else if version || (len(os.Args) == 2 && os.Args[1] == "version") {
		fmt.Println(CSVNORM_VERSION)
		os.Exit(0)
	} else if flag.NArg() > 0 {
		fmt.Println("Error: Provided invalid arguments:", strings.Join(flag.Args(), ", "))
		fmt.Println("Run \"csvnorm -help\" for more information.")
		os.Exit(0)
	}

	defaultInputLocation, err := time.LoadLocation(defaultInputTimezoneCode)
	if err != nil {
		fmt.Println(errors.Wrap(err, "Could not load default input timezone location with code "+defaultInputTimezoneCode))
		os.Exit(1)
	}

	outputLocation, err := time.LoadLocation(outputTimezoneCode)
	if err != nil {
		fmt.Println(errors.Wrap(err, "Could not load output timezone location with code "+outputTimezoneCode))
		os.Exit(1)
	}

	reader := csv.NewReader(bufio.NewReader(os.Stdin))
	writer := csv.NewWriter(bufio.NewWriter(os.Stdout))

	header, err := reader.Read()
	if err != nil {
		if err != io.EOF {
			panic(err)
		}
	}

	writer.Write(header)
	writer.Flush()

	for {

		inRow, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		input := map[string]string{}
		for i, h := range header {
			input[strings.ToLower(h)] = inRow[i]
		}

		fullname := normalizeString(input["fullname"])
		address := normalizeString(input["address"])
		notes := normalizeString(input["notes"])

		timestamp, err := parseTimestamp(normalizeString(input["timestamp"]), defaultInputLocation)
		if err != nil {
			fmt.Fprintf(os.Stderr, errors.Wrap(err, "Error parsing timestamp").Error()+"\n")
			continue
		}

		fooDuration, err := parseDuration(normalizeString(input["fooduration"]))
		if err != nil {
			fmt.Fprintf(os.Stderr, errors.Wrap(err, "Error parsing FooDuration").Error()+"\n")
			continue
		}

		barDuration, err := parseDuration(normalizeString(input["barduration"]))
		if err != nil {
			fmt.Fprintf(os.Stderr, errors.Wrap(err, "Error parsing BarDuration").Error()+"\n")
			continue
		}

		zipcode_int, err := strconv.Atoi(normalizeString(input["zip"]))
		if err != nil {
			fmt.Fprintf(os.Stderr, errors.Wrap(err, "Error parsing zipcode").Error()+"\n")
			continue
		}

		data := map[string]string{
			"FullName":      strings.ToUpper(fullname),
			"Address":       address,
			"Timestamp":     timestamp.In(outputLocation).Format(time.RFC3339),
			"ZIP":           fmt.Sprintf("%05d", zipcode_int),
			"FooDuration":   strconv.FormatFloat(fooDuration.Seconds(), 'f', -1, 64),
			"BarDuration":   strconv.FormatFloat(barDuration.Seconds(), 'f', -1, 64),
			"TotalDuration": strconv.FormatFloat((fooDuration.Seconds() + barDuration.Seconds()), 'f', -1, 64),
			"Notes":         notes,
		}

		outRow := make([]string, len(header))
		for i, h := range header {
			outRow[i] = data[h]
		}
		writer.Write(outRow)
		writer.Flush()

	}

}
