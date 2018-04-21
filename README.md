# csvnorm

# Description

**csvnorm** is a tool for normalizing CSV files.  **csvnorm** is built in [Go](https://golang.org/).

# Usage

To use **csvnorm**, download the appropriate executable from https://github.com/pjdufour/csvnorm/releases.  Executables for Linux, Darwin, and Windows are provided.

```
Usage: csvnorm_linux_amd64 [-version] [-help]
```

A typical use case would be to use bash to pipe a csv through csvnorm and save to a file.

```
cat raw.csv | ./csvnorm_linux_amd64 > normalized.csv
```

# Building

Run the `build.sh` script to build executables for Linux, Windows, and Darwin.

```
bash scripts/build.sh
```

More advanced builds can be created using standard Go methods.

# Examples


# Contributing

We are currently accepting pull requests for this repository.  We'd love to have your contributions!  Please see [Contributing.md](https://github.com/pjdufour/csvnorm/blob/master/CONTRIBUTING.md) for how to get started.

# License

This work is distributed under the **MIT License**.  See **LICENSE** file.
