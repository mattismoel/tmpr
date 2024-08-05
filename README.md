# What is this?

This is a small command-line utility for retrieving the forecast at a specific location or set of coordinates.

# How do i install?

The program is simple to install if you have Golang installed on your system. Just clone the directory onto your system and run the following commands:

```bash
go build -o "./bin/tmpr" main.go # Creates a build of the program in the project diretory in directory "bin/".

```
The binary can now be run like this:

```bash
./bin/tmpr --help
```
Alternatively the binary can be run directly with Golang:

```bash
go run main.go --help
```

# Usage

Help can always be found with `tmpr --help`. The basics are as following:

Get forecast at specific coordinates:

```bash
tmpr coords --unit="<unit>" --lat="<latitude>" --lon="<longitude>"
```

Get forecast given a query:

```bash
tmpr query --unit="<unit>" --query="New York, USA"
```

> The `--unit` flag is required, and can be of following systems: `"metric", "imperial", "standard"`.

