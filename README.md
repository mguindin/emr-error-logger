# EMR Error Log Downloader

This is a simple binary that writes out the error log (stored in S3) from a failed EMR cluster and opens it in $PAGER for viewing

## Prereqs
You need to have AWS credentials stored in a proper place (e.g.
`~/.aws/credentials`)

## Installation
First install `vgo`:
```sh
  go get -u golang.org/x/vgo
```

Then install:
```sh
  vgo install
```

Usage:
```sh
Usage of emr-error-logger:
  -dc string
        datacenter (default "us-east-1")
  -id string
        cluster ID
  -o    only output file name
```

## TODO:
1. Handle bootstrap errors
1. Write tests
