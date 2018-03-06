package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/aws/aws-sdk-go/service/emr"
)

var (
	dc             = flag.String("dc", "us-east-1", "datacenter")
	clusterID      = flag.String("id", "", "cluster ID")
	out            = flag.Bool("o", false, "only output file name")
	bootstrapRegex = regexp.MustCompile(`bootstrap action (\d+) returned`)
)

func main() {
	flag.Parse()
	emrErrorFinder := NewEMRErrorFinder(clusterID, dc)
	clusterOutput := emrErrorFinder.GetCluster()
	clusterMsg := *clusterOutput.Cluster.Status.StateChangeReason.Message
	clusterCode := *clusterOutput.Cluster.Status.StateChangeReason.Code
	log.Printf("cluster error message: %s, Code: %s\n", clusterMsg, clusterCode)
	var errorLog string
	if clusterCode == emr.ClusterStateChangeReasonCodeStepFailure {
		errorLog = emrErrorFinder.GetFailedStepErrorLog()
	} else if clusterCode == emr.ClusterStateChangeReasonCodeBootstrapFailure {
		step := getFailedBootstrapStepNumber(clusterMsg)
		errorLog = emrErrorFinder.GetBootstrapFailureErrorLog(step)
	}
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

func getFailedBootstrapStepNumber(clusterMsg string) (stepNo int) {
	matches := bootstrapRegex.FindStringSubmatch(clusterMsg)
	if len(matches) < 2 {
		log.Print("Unable to get failed bootstrap step number")
		return
	}
	stepNo, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Print("Unable to get failed bootstrap step number")
		return
	}
	return
}
