package goreport

import (
	"testing"
)

func TestReport(t *testing.T) {

}

type TestDetail struct {
}

func (h *TestDetail) GetHeight() float64 {
	return 10
}
func (h *TestDetail) Execute(report *GoReport) {

}
