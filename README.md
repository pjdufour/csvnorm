[![Build Status](https://travis-ci.org/pjdufour/csvnorm.svg)](https://travis-ci.org/pjdufour/csvnorm)

# csvnorm

# Description

**csvnorm** is a tool for normalizing CSV files ingested via `stdin` and output to `stdout`.  **csvnorm** is built in [Go](https://golang.org/).

# Usage

To use **csvnorm**, download the appropriate executable from https://github.com/pjdufour/csvnorm/releases.  Executables for Linux, Darwin (for Mac OSX), and Windows are provided.

```
Usage: csvnorm [-input_timezone INPUT_TIMEZONE] [-output_timezone OUTPUT_TIMEZONE] [-help] [-version]
  -help
    	Print help
  -input_timezone string
    	Default timezone for timestamp input.  Matches names in the IANA Time Zone database. (default "US/Pacific")
  -output_timezone string
    	Timezone for timestamp output.  Matches names in the IANA Time Zone database. (default "US/Eastern")
  -version
    	Prints version to stdout
```

A typical use case would be to use bash to pipe a csv through csvnorm and save to a file.

```
cat raw.csv | ./csvnorm_linux_amd64 > normalized.csv
```

In a deployment, the executables can be renamed as simply `csvnorm` for convenience.

# Building

**csvnorm** is built in [Go](https://golang.org/).  To get started, download the applicable version of Go from the [Downloads](https://golang.org/dl/) page and then follow the [Installation Instructions](https://golang.org/doc/install).  Then install all necessary dependencies with:

```
go get -d ./...
```

Finally, to build **csvnorm** run the `scripts/build.sh` script to build executables for Linux, Windows, and Darwin.

```
bash scripts/build.sh
```

Advanced builds can be created using standard Go methods.  You can also create your own simple build script or run `go build` directly as needed.  Consult the [How to Write Go Code](https://golang.org/doc/code.html) page for more information.

# Testing

To run the unit tests use Go's test suite as described on the [Command go](https://golang.org/cmd/go/) page.

```
cd cmd/csvnorm
go test
```

You can also run the tests using the long form `go test github.com/pjdufour/csvnorm/cmd/csvnorm`.  For an end-to-end test against the files in the `examples` folder, run the `scripts/test.sh` script.

```
bash scripts/test.sh
```

# Contributing

We are currently accepting pull requests for this repository.  We'd love to have your contributions!  Please see [Contributing.md](https://github.com/pjdufour/csvnorm/blob/master/CONTRIBUTING.md) for how to get started.

# License

This work is distributed under the **MIT License**.  See **LICENSE** file.
