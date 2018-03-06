package main

import (
	"testing"
)

func TestCreateS3GetObject(t *testing.T) {
	logFileWithGz := "s3://my-bucket/path/to/file/stderr.gz"
	expectedBucket := "my-bucket"
	expectedKey := "path/to/file/stderr.gz"

	output := CreateS3GetObject(logFileWithGz)
	if *output.Bucket != expectedBucket && *output.Key != expectedKey {
		t.Errorf("Failed to output proper s3 getObject with bucket %s and key %s",
			expectedBucket,
			expectedKey,
		)
	}

	logFileWithoutGz := "s3://my-bucket/path/to/file/"
	output = CreateS3GetObject(logFileWithoutGz)
	if *output.Bucket != expectedBucket && *output.Key != expectedKey {
		t.Errorf("Failed to output proper s3 getObject with bucket %s and key %s",
			expectedBucket,
			expectedKey,
		)
	}

	emptyLogFile := ""
	if CreateS3GetObject(emptyLogFile) != nil {
		t.Error("Failed to output nil s3GetObject")
	}
}
