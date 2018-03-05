package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/emr"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	dc        = flag.String("dc", "us-east-1", "datacenter")
	clusterID = flag.String("id", "", "cluster ID")
)

func main() {
	flag.Parse()
	awsConfig := aws.NewConfig().WithRegion(*dc)
	sess := session.Must(session.NewSession())
	emrClient := emr.New(sess, awsConfig)
	s3Client := s3.New(sess, awsConfig)
	clusterOutput := getCluster(emrClient, clusterID)
	log.Printf("cluster error message: %s\n", *clusterOutput.Cluster.Status.StateChangeReason.Message)
	if *clusterOutput.Cluster.Status.StateChangeReason.Code == "STEP_FAILURE" {
		errorLog := getFailedStepErrorLog(emrClient, clusterID)
		if errorLog != "" {
			s3FileToLocal(s3Client, createS3GetObject(errorLog))
		} else {
			log.Printf("Blank error log")
		}
	}
}

func s3FileToLocal(s3Client *s3.S3, objInput *s3.GetObjectInput) {
	obj, err := s3Client.GetObject(objInput)
	if err != nil {
		log.Fatalf("Unable to fetch error log %s", err.Error())
	}
	fo, err := os.Create("error-log.log")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := fo.Close(); err != nil {
			log.Fatalf("Error closing file being written: %s", err.Error())
		}
	}()
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := obj.Body.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatalf("Error reading from s3 file: %s", err.Error())
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := fo.Write(buf[:n]); err != nil {
			log.Fatalf("Error writing out file: %s", err.Error())
		}
	}
	log.Print("Wrote log to `error-log.log`")
}

func getCluster(client *emr.EMR, clusterID *string) *emr.DescribeClusterOutput {
	describeClusterInput := &emr.DescribeClusterInput{
		ClusterId: clusterID,
	}
	out, err := client.DescribeCluster(describeClusterInput)
	if err != nil {
		log.Fatal("Unable to get cluster information")
	}

	if !strings.Contains(*out.Cluster.Status.State, "ERROR") {
		log.Fatal("Cluster terminated normally, no errors to fetch")
	}
	log.Printf("cluster state %s\n", *out.Cluster.Status.State)
	return out
}

func getFailedStepErrorLog(client *emr.EMR, clusterID *string) string {
	listStepsInput := &emr.ListStepsInput{
		ClusterId: clusterID,
	}

	pageNum := 0
	var logFile string
	err := client.ListStepsPages(listStepsInput,
		func(page *emr.ListStepsOutput, lastPage bool) bool {
			for _, step := range page.Steps {
				if *step.Status.State == "FAILED" {
					log.Printf("Step %v failed\n", *step.Id)
					log.Printf("failure log %s", *step.Status.FailureDetails.LogFile)
					logFile = *step.Status.FailureDetails.LogFile
				}
			}
			pageNum++
			return !lastPage
		})
	if err != nil {
		log.Fatal("Unable to fetch list of steps")
	}

	return logFile
}

func createS3GetObject(logFile string) *s3.GetObjectInput {
	s3LessLog := strings.Replace(logFile, "s3://", "", 1)
	splits := strings.SplitN(s3LessLog, "/", 2)
	bucket := splits[0]
	key := splits[1]
	if !strings.Contains(key, "stderr.gz") {
		key = key + "stderr.gz"
	}
	log.Printf("s3 bucket: %s, s3 key: %s", bucket, key)
	return &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
}
