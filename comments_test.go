package main

import "testing"

func TestGetPort(t *testing.T) {
	defaultPortValue := ":8000"
	port := getPort()
	if port != defaultPortValue {
		t.Error("getPort default value returned as", port, ", correct value should be ':8000'")
	}
}
