package main

import "testing"

func TestParseFlags(t *testing.T) {
	defaultPortValue := ":8000"
	defaultDisableCORSValue := false
	parseFlags()
	if port != defaultPortValue {
		t.Error("parseFlag port default value returned as", port, ", correct value should be ':8000'")
	}
	if disableCORS != defaultDisableCORSValue {
		t.Error("parseFlag disableCORS default value returned as true, correct value should be false")
	}
}
