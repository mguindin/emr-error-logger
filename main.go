package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

var (
	dc        = flag.String("dc", "us-east-1", "datacenter")
	clusterID = flag.String("id", "", "cluster ID")
	out       = flag.Bool("o", false, "only output file name")
)

func main() {
	flag.Parse()
	emrErrorFinder := NewEMRErrorFinder(clusterID, dc)
	clusterOutput := emrErrorFinder.GetCluster()
	log.Printf("cluster error message: %s\n", *clusterOutput.Cluster.Status.StateChangeReason.Message)
	if *clusterOutput.Cluster.Status.StateChangeReason.Code == "STEP_FAILURE" {
		errorLog := emrErrorFinder.GetFailedStepErrorLog()
		if errorLog != "" {
			file := emrErrorFinder.S3FileToLocal(CreateS3GetObject(errorLog))
			if !*out {
				pagerErrorFile(file)
			} else {
				fmt.Print(file)
			}
		} else {
			log.Printf("Blank error log")
		}
	}
}

func pagerErrorFile(file string) {
	if pager, ok := os.LookupEnv("PAGER"); ok {
		cmd := exec.Command(pager, file)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Unable to run pager: %s", err.Error())
		}
		return
	}
	log.Printf("No $PAGER set, outputting file name")
	log.Println(file)
}
