package example

import (
	//"fmt"
	gr "github.com/mikeshimura/goreport"
	//"io/ioutil"
	"strconv"
	//"strings"
)

func Complex1Report() {
	r := gr.CreateGoReport()
	//Page Total Function
	r.PageTotal = true
	r.SumWork["g1amtcum"] = 0.0
	r.SumWork["g2amtcum"] = 0.0
	r.SumWork["g1hrcum"] = 0.0
	r.SumWork["g2hrcum"] = 0.0
	r.SumWork["g2item"] = 0.0
	font1 := gr.FontMap{
		FontName: "IPAexG",
		FileName: "ttf//ipaexg.ttf",
	}
	font2 := gr.FontMap{
		FontName: "MPBOLD",
		FileName: "ttf//mplus-1p-bold.ttf",
	}
	fonts := []*gr.FontMap{&font1, &font2}
	r.SetFonts(fonts)
	d := new(C1Detail)
	r.RegisterBand(gr.Band(*d), gr.Detail)
	h := new(C1Header)
	r.RegisterBand(gr.Band(*h), gr.PageHeader)
	f := new(C1Footer)
	r.RegisterBand(gr.Band(*f), gr.PageFooter)
	s1h := new(C1G1Header)
	r.RegisterGroupBand(gr.Band(*s1h), gr.GroupHeader, 1)
	s1 := new(C1G1Summary)
	r.RegisterGroupBand(gr.Band(*s1), gr.GroupSummary, 1)
	s2 := new(C1G2Summary)
	r.RegisterGroupBand(gr.Band(*s2), gr.GroupSummary, 2)
	r.Records = gr.ReadTextFile("invoice.txt", 12)
	//fmt.Printf("Records %v \n", r.Records)
	r.SetPage("A4", "mm", "P")
	r.SetFooterY(265)
	r.Execute("complex1.pdf")
	r.SaveText("complex1.txt")
}

type C1Detail struct {
}

func (h C1Detail) GetHeight(report gr.GoReport) float64 {
	return 6
}
func (h C1Detail) Execute(report gr.GoReport) {
	cols := report.Records[report.DataPos].([]string)
	y := 2.0
	report.Font("IPAexG", 9, "")
	report.Cell(14, y, cols[5])
	report.Cell(40, y, cols[6])
	hr := gr.ParseFloatPanic(cols[7], "")
	report.CellRight(150, y, 20, gr.AddComma(strconv.FormatFloat(hr, 'f', 1, 64))+" Hrs")
	amt := gr.ParseFloatPanic(cols[8], "")
	report.CellRight(170, y, 26, gr.AddComma(strconv.FormatFloat(amt, 'f', 2, 64))+" USD")
	report.SumWork["g1amtcum"] += amt
	report.SumWork["g2amtcum"] += amt
	report.SumWork["g1hrcum"] += hr
	report.SumWork["g2hrcum"] += hr
}
func (h C1Detail) BreakCheckBefore(report gr.GoReport) int {
	if report.DataPos == 0 {
		//max no
		return 2
	}
	curr := report.Records[report.DataPos].([]string)
	before := report.Records[report.DataPos-1].([]string)
	return h.BreakCheckSub(curr, before)
}
func (h C1Detail) BreakCheckAfter(report gr.GoReport) int {
	if report.DataPos == len(report.Records)-1 {
		//max no
		return 2
	}
	curr := report.Records[report.DataPos].([]string)
	after := report.Records[report.DataPos+1].([]string)
	return h.BreakCheckSub(curr, after)
}
func (h C1Detail) BreakCheckSub(row1 []string, row2 []string) int {
	if row1[0] != row2[0] {
		return 2
	}
	if row1[4] != row2[4] {
		return 1
	}
	return 0
}

type C1Header struct {
}

func (h C1Header) GetHeight(report gr.GoReport) float64 {
	if report.SumWork["g2item"] == 0.0 {
		return 116
	}
	return 38
}
func (h C1Header) Execute(report gr.GoReport) {
	cols := report.Records[report.DataPos].([]string)
	y := 32.0
	if report.SumWork["g2item"] == 0.0 {
		report.Image("apple.jpg", 20, 35, 35, 50)
		report.Font("MPBOLD", 18, "")
		report.LineType("straight", 1)
		report.GrayStroke(0.9)
		report.LineV(49, 72, 90)
		report.LineV(150, 43, 67)
		report.LineV(150, 71, 95)
		report.GrayStroke(0)
		//	report.LineType("straight", 0.5)
		//	report.Rect(48, 13, 81, 21)
		report.Cell(145, 33, "TAX INVOICE")
		report.Font("MPBOLD", 9, "")
		report.Cell(139, 45, "From")
		x := 153.0
		report.Cell(x, 45, "Test Consulting Corp.")
		report.Cell(x, 51, "123 Hyde Street")
		report.Cell(x, 57, "San Francisco, Calfornia")
		report.Cell(x, 63, "USA")

		report.Cell(139, 74, "To")
		report.Cell(x, 74, cols[0])
		report.Cell(x, 80, cols[1])
		report.Cell(x, 86, cols[2])
		report.Cell(x, 92, cols[3])

		x = 14.0
		report.Cell(x, 73, "Tax Invoice No:")
		report.Cell(x, 79, "Tax Invoice Date:")
		report.Cell(x, 85, "Payment Due Date:")

		x = 52
		report.Cell(x, 73, cols[9])
		report.Cell(x, 79, cols[10])
		report.Cell(x, 85, cols[11])

		y = 110
		y = y
	}
	report.LineType("straight", 7)
	report.GrayStroke(0.9)
	report.LineH(11, y-2, 199)
	report.GrayStroke(0)
	report.Cell(14, y, "Type")
	report.Cell(40, y, "Description")
	report.Cell(161, y, "Hours")
	report.Cell(184, y, "Amount")
	report.SumWork["g2item"] = 1.0
}

type C1G1Summary struct {
}

func (h C1G1Summary) GetHeight(report gr.GoReport) float64 {
	return 7
}
func (h C1G1Summary) Execute(report gr.GoReport) {
	y := 2.0
	report.LineType("straight", 1)
	report.GrayStroke(0.9)
	report.LineH(11, 0, 199)
	report.GrayStroke(0)
	report.Font("MPBOLD", 9, "")
	report.CellRight(150, y, 20, gr.AddComma(strconv.FormatFloat(
		report.SumWork["g1hrcum"], 'f', 1, 64))+" Hrs")
	report.CellRight(170, y, 26, gr.AddComma(strconv.FormatFloat(
		report.SumWork["g1amtcum"], 'f', 2, 64))+" USD")
	report.SumWork["g1amtcum"] = 0.0
	report.SumWork["g1hrcum"] = 0.0
}

type C1G2Summary struct {
}

func (h C1G2Summary) GetHeight(report gr.GoReport) float64 {
	return 50
}
func (h C1G2Summary) Execute(report gr.GoReport) {
	report.Font("MPBOLD", 9, "")
	y := 15.0
	report.CellRight(123, y, 20, "Total:")
	report.CellRight(150, y, 20, gr.AddComma(strconv.FormatFloat(
		report.SumWork["g2hrcum"], 'f', 1, 64))+" Hrs")
	report.CellRight(170, y, 26, gr.AddComma(strconv.FormatFloat(
		report.SumWork["g2amtcum"], 'f', 2, 64))+" USD")
	y = 25.0
	report.CellRight(123, y, 20, "Tax:")
	report.CellRight(150, y, 20, "7.75%")
	tax := report.SumWork["g2amtcum"] * 0.0775
	report.CellRight(170, y, 26, gr.AddComma(strconv.FormatFloat(
		tax, 'f', 2, 64))+" USD")
	report.LineType("straight", 0.3)
	report.LineH(170, 33, 199)
	y = 39.0
	report.Font("MPBOLD", 11, "")
	report.CellRight(123, y, 20, "AMOUT DUE:")
	report.CellRight(170, y, 26, gr.AddComma(strconv.FormatFloat(
		report.SumWork["g2amtcum"]+tax, 'f', 2, 64))+" USD")
	report.NewPage(true)
	report.SumWork["g2item"] = 0.0
	report.SumWork["g2hrcum"] = 0.0
	report.SumWork["g2amtcum"] = 0.0
}

type C1G1Header struct {
}

func (h C1G1Header) GetHeight(report gr.GoReport) float64 {
	return 8
}
func (h C1G1Header) Execute(report gr.GoReport) {
	cols := report.Records[report.DataPos].([]string)
	y := 2.0
	report.Font("MPBOLD", 9, "")
	report.Cell(14, y, "SUB-TASK")
	report.Cell(40, y, cols[4])
	report.LineType("straight", 1)
	report.GrayStroke(0.9)
	report.LineH(11, 7, 199)
	report.GrayStroke(0)
}

type C1Footer struct {
}

func (h C1Footer) GetHeight(report gr.GoReport) float64 {
	return 10
}
func (h C1Footer) Execute(report gr.GoReport) {
	report.Cell(100, 12, "Page")
	report.Cell(112, 12, strconv.Itoa(report.Page))
}
