package main

import (
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	d := "2022-07-08T15:56:29.784+0200"
	v, err := time.Parse(RFCJava, d)
	t.Log(v, err)
	if err != nil {
		t.FailNow()
	}
}
