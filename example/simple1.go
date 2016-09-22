package example

import (
	//"fmt"
	gr "github.com/mikeshimura/goreport"
	//"io/ioutil"
	"strconv"
	//"strings"
)

func Simple1Report() {
	r := gr.CreateGoReport()
	r.SumWork["amountcum="] = 0.0
	font1 := gr.FontMap{
		FontName: "IPAexG",
		FileName: "ttf//ipaexg.ttf",
	}
	fonts := []*gr.FontMap{&font1}
	r.SetFonts(fonts)
	d := new(S1Detail)
	r.RegisterBand(gr.Band(*d), gr.Detail)
	h := new(S1Header)
	r.RegisterBand(gr.Band(*h), gr.PageHeader)
	s := new(S1Summary)
	r.RegisterBand(gr.Band(*s), gr.Summary)
	r.Records = gr.ReadTextFile("sales1.txt", 7)
	//fmt.Printf("Records %v \n", r.Records)
	r.SetPage("A4", "mm", "L")
	r.SetFooterY(190)
	r.Execute("simple1.pdf")
	r.SaveText("simple1.txt")
}

type S1Detail struct {
}

func (h S1Detail) GetHeight(report gr.GoReport) float64 {
	return 10
}
func (h S1Detail) Execute(report gr.GoReport) {
	cols := report.Records[report.DataPos].([]string)
	report.Font("IPAexG", 12, "")
	y := 2.0
	report.Cell(15, y, cols[0])
	report.Cell(30, y, cols[1])
	report.Cell(60, y, cols[2])
	report.Cell(90, y, cols[3])
	report.Cell(120, y, cols[4])
	report.CellRight(135, y, 25, cols[5])
	report.CellRight(160, y, 20, cols[6])
	amt := ParseFloatNoError(cols[5]) * ParseFloatNoError(cols[6])
	report.SumWork["amountcum="] += amt
	report.CellRight(180, y, 30, strconv.FormatFloat(amt, 'f', 2, 64))
}
func ParseFloatNoError(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

type S1Header struct {
}

func (h S1Header) GetHeight(report gr.GoReport) float64 {
	return 30
}
func (h S1Header) Execute(report gr.GoReport) {
	report.Font("IPAexG", 14, "")
	report.Cell(50, 15, "Sales Report")
	report.Font("IPAexG", 12, "")
	report.Cell(240, 20, "page")
	report.Cell(260, 20, strconv.Itoa(report.Page))
	y := 23.0
	report.Cell(15, y, "D No")
	report.Cell(30, y, "Dept")
	report.Cell(60, y, "Order")
	report.Cell(90, y, "Stock")
	report.Cell(120, y, "Name")
	report.CellRight(135, y, 25, "Unit Price")
	report.CellRight(160, y, 20, "Qty")
	report.CellRight(190, y, 20, "Amount")
}

type S1Summary struct {
}

func (h S1Summary) GetHeight(report gr.GoReport) float64 {
	return 10
}
func (h S1Summary) Execute(report gr.GoReport) {
	report.Cell(160, 2, "Total")
	report.CellRight(180, 2, 30, strconv.FormatFloat(
		report.SumWork["amountcum="], 'f', 2, 64))
}
