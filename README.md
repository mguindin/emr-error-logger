# emr-error-logger

Usage:
```sh
Usage of emr-error-logger:
  -dc string
        datacenter (default "us-east-1")
  -id string
        cluster ID
```

This is a simple binary that writes out the error log (stored in S3) from a failed EMR cluster

TODO:
1. Handle bootstrap errors
1. Write tests

