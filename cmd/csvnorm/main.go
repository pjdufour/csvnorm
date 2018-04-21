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

	var version bool
	var help bool

	flag.BoolVar(&version, "version", false, "Prints version to stdout")
	flag.BoolVar(&help, "help", false, "Print help")

	flag.Parse()

	if help || (len(os.Args) == 2 && os.Args[1] == "help") {
		fmt.Println("Usage: csvnorm [-version] [-help]")
		flag.PrintDefaults()
		os.Exit(0)
	} else if flag.NArg() > 0 {
		fmt.Println("Error: Provided invalid arguments:", strings.Join(flag.Args(), ", "))
		fmt.Println("Run \"csvnorm -help\" for more information.")
		os.Exit(0)
	}

	if version {
		fmt.Println(CSVNORM_VERSION)
		os.Exit(0)
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

		fullname, ok := input["fullname"]
		if !ok {
			panic(errors.New("Missing FullName"))
		}
		fullname = strings.ToUpper(fullname)

		timestamp, err := time.Parse("1/2/06 15:04:05 PM", input["timestamp"])
		if err != nil {
			panic(errors.Wrap(err, "Error parsing timestamp"))
		}

		fooDuration, err := parseDuration(input["fooduration"])
		if err != nil {
			panic(errors.Wrap(err, "Error parsing FooDuration"))
		}

		barDuration, err := parseDuration(input["barduration"])
		if err != nil {
			panic(errors.Wrap(err, "Error parsing BarDuration"))
		}

		zipcode := input["zip"]
		if len(zipcode) < 5 {
			zipcode_int, err := strconv.Atoi(zipcode)
			if err != nil {
				panic(err)
			}
			zipcode = fmt.Sprintf("%05d", zipcode_int)
		}

		address := normalizeString(input["address"])
		notes := normalizeString(input["notes"])

		data := map[string]string{
			"FullName":      fullname,
			"Address":       address,
			"Timestamp":     timestamp.Format(time.RFC3339),
			"ZIP":           zipcode,
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
