package main

import (
	"testing"
)

func TestGetFailedBootstrapStepNumber(t *testing.T) {
	goodMsg := "bootstrap action 4 returned non-zero"
	if getFailedBootstrapStepNumber(goodMsg) != 4 {
		t.Error("Regex failed to match proper number (4) from msg")
	}

	nonNumMsg := "bootstrap action flux returned non-zero"
	if getFailedBootstrapStepNumber(nonNumMsg) != 0 {
		t.Error("Regex failed to return 0 from mis-matched cluster message")
	}

	negNumMsg := "bootstrap action -1 returned non-zero"
	if getFailedBootstrapStepNumber(negNumMsg) != 0 {
		t.Error("Regex failed to return 0 from mis-matched cluster message")
	}
}
