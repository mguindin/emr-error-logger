package main

import (
	"flag"
	"log"
)

var (
	dc        = flag.String("dc", "us-east-1", "datacenter")
	clusterID = flag.String("id", "", "cluster ID")
)

func main() {
	flag.Parse()
	emrErrorFinder := NewEMRErrorFinder(clusterID, dc)
	clusterOutput := emrErrorFinder.GetCluster()
	log.Printf("cluster error message: %s\n", *clusterOutput.Cluster.Status.StateChangeReason.Message)
	if *clusterOutput.Cluster.Status.StateChangeReason.Code == "STEP_FAILURE" {
		errorLog := emrErrorFinder.GetFailedStepErrorLog()
		if errorLog != "" {
			emrErrorFinder.S3FileToLocal(CreateS3GetObject(errorLog))
		} else {
			log.Printf("Blank error log")
		}
	}
}
