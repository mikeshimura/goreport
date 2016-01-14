package goreport

import(
	"testing"
)

func TestReport(t *testing.T) {
	rep:=CreateGoReport()
	d := new(TestDetail)
	var det Band = d
	rep.RegisterBand(&det, Detail)
	
}

type TestDetail struct {
}

func (h *TestDetail) GetHeight() float64 {
	return 10
}
func (h *TestDetail) Execute(report *GoReport) {

}