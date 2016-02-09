package server

import (
	"testing"
)

func TestGetTime(t *testing.T) {
	in  := "AAPL,0,20160127,093000,96.090000000,100"
	expected := "093000"
	out := getTime(in)

	if out != expected { t.Errorf("Expected time %s. Got: %s", expected, out)}
}