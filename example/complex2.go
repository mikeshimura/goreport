package example

import (
	//"fmt"
	"time"
	"unicode"

	gr "github.com/mikeshimura/goreport"

	//"io/ioutil"
	"strconv"
	"strings"
)

const (
	C2DateFormat = "2006年01月02日"
)

func Complex2Report() {

	r := gr.CreateGoReport()
	//Page Total Function
	r.PageTotal = true
	r.SumWork["g1amtcum"] = 0.0
	r.SumWork["g2amtcum"] = 0.0
	r.SumWork["g1item"] = 0.0
	r.SumWork["g2item"] = 0.0
	font1 := gr.FontMap{
		FontName: "IPAexゴシック",
		FileName: "ttf//ipaexg.ttf",
	}
	font2 := gr.FontMap{
		FontName: "MPBOLD",
		FileName: "ttf//mplus-1p-bold.ttf",
	}
	fonts := []*gr.FontMap{&font1, &font2}
	r.SetFonts(fonts)
	d := new(C2Detail)
	r.RegisterBand(gr.Band(*d), gr.Detail)
	h := new(C2Header)
	r.RegisterBand(gr.Band(*h), gr.PageHeader)
	f := new(C2Footer)
	r.RegisterBand(gr.Band(*f), gr.PageFooter)
	s1 := new(C2G1Summary)
	r.RegisterGroupBand(gr.Band(*s1), gr.GroupSummary, 1)
	s2 := new(C2G2Summary)
	r.RegisterGroupBand(gr.Band(*s2), gr.GroupSummary, 2)
	r.Records = gr.ReadTextFile("invoice2.txt", 11)
	//fmt.Printf("Records %v \n", r.Records)
	r.SetPage("A4", "mm", "P")
	r.SetFooterY(265)
	r.Execute("complex2.pdf")
	r.SaveText("complex2.txt")
}

type C2Detail struct {
}

func (h C2Detail) GetHeight(report gr.GoReport) float64 {
	return 5
}
func (h C2Detail) Execute(report gr.GoReport) {
	slipSHow := true

	cols := report.Records[report.DataPos].([]string)
	if report.SumWork["g1item"] > 0 {
		bfr := report.Records[report.DataPos-1].([]string)
		if cols[4] == bfr[4] {
			slipSHow = false
		}
	}
	y := 1.5
	x := 25.0
	report.LineType("straight", 0.3)
	report.GrayStroke(0)
	report.Rect(x+41, 0, x+160, 5)
	report.LineV(x, 0, 5)
	report.LineV(x+23, 0, 5)
	report.LineV(x+61, 0, 5)
	report.LineV(x+99, 0, 5)
	report.LineV(x+116, 0, 5)
	report.LineV(x+135, 0, 5)
	fty := report.SumWork["__ft__"]
	//最下行なら横線を引く
	//fmt.Printf("fty %v CurrY %v\n", fty, report.CurrY)
	if fty-report.CurrY <= 5 {
		report.LineH(x, 5, x+41)
	}
	report.Font("IPAexゴシック", 9, "")
	if slipSHow {
		report.Cell(x+1, y, cols[3])
		report.Cell(x+24, y, cols[4])
	}
	report.Cell(x+42, y, cols[5])
	report.Cell(x+62, y, cols[6])
	report.CellRight(x+115, y, 0, gr.AddComma(cols[7]))
	report.CellRight(x+134, y, 0, "\u00A5"+gr.AddComma(cols[8]))
	report.CellRight(x+159, y, 0, "\u00A5"+gr.AddComma(cols[9]))
	amt := gr.AtoiPanic(cols[9], "")
	report.SumWork["g1amtcum"] += float64(amt)
	report.SumWork["g2amtcum"] += float64(amt)
	report.SumWork["g1item"]++
	report.SumWork["g2item"]++
}
func (h C2Detail) BreakCheckBefore(report gr.GoReport) int {
	if report.DataPos == 0 {
		//max no
		return 2
	}
	curr := report.Records[report.DataPos].([]string)
	before := report.Records[report.DataPos-1].([]string)
	return h.BreakCheckSub(curr, before)
}
func (h C2Detail) BreakCheckAfter(report gr.GoReport) int {
	if report.DataPos == len(report.Records)-1 {
		//max no
		return 2
	}
	curr := report.Records[report.DataPos].([]string)
	after := report.Records[report.DataPos+1].([]string)
	return h.BreakCheckSub(curr, after)
}
func (h C2Detail) BreakCheckSub(row1 []string, row2 []string) int {
	if row1[0] != row2[0] {
		return 2
	}
	if row1[4] != row2[4] {
		return 1
	}
	return 0
}

type C2Header struct {
}

func (h C2Header) GetHeight(report gr.GoReport) float64 {
	if report.SumWork["g2item"] == 0.0 {
		return 90
	}
	return 37
}
func (h C2Header) Execute(report gr.GoReport) {
	cols := report.Records[report.DataPos].([]string)
	y := 32.0
	x := 25.0
	if report.SumWork["g2item"] == 0.0 {
		numConv := unicode.SpecialCase{
			// 半角の 0 から 9 に対する変換ルール
			unicode.CaseRange{
				0x0030, // Lo: 半角の 0
				0x0039, // Hi: 半角の 9
				[unicode.MaxCase]rune{
					0xff10 - 0x0030, // UpperCase で全角に変換
					0,               // LowerCase では変換しない
					0xff10 - 0x0030, // TitleCase で全角に変換
				},
			},
		}
		report.Font("IPAexゴシック", 10, "")

		now := strings.ToUpperSpecial(
			numConv, time.Now().Format(C2DateFormat))
		report.CellRight(182, 12, 0, now)
		report.CellRight(182, 19, 0, "請求書番号:"+cols[2])
		report.Font("MPBOLD", 16, "")
		report.Cell(92, 30, "請求書")
		report.Font("MPBOLD", 11, "")
		x = 118.0
		report.Cell(x, 46, "サンプル商事株式会社")
		report.Cell(x, 51, "山田太郎")
		report.Font("IPAexゴシック", 9, "")
		report.Cell(x, 58, "〒181-0001")
		report.Cell(x, 62, "東京都三鷹市井の頭5-12-12")
		report.Cell(x, 72, "TEL:0422-22-2222")
		report.Cell(x, 76, "FAX:0422-22-2223")
		report.Cell(x, 80, "info@mitakashoji.jp")
		cols = cols
		x = 25
		report.Font("MPBOLD", 12, "")
		company := cols[1] + " 御中"
		// 注意 Pdf.SetFontは呼び出し順が違う
		report.Converter.Pdf.SetFont("IPAexゴシック", "", 12)
		w, _ := report.Converter.Pdf.MeasureTextWidth(company)
		report.Cell(x, 46, company)
		report.LineType("straight", 0.3)
		report.GrayStroke(0.5)
		report.LineH(x, 50, x+w/report.ConvPt)
		report.Font("IPAexゴシック", 9, "")
		report.Cell(x, 58, "下記のとおりご請求申し上げます。")
		report.Font("MPBOLD", 12, "")
		report.Cell(x, 70, "ご請求金額")
		report.CellRight(x+72, 70, 0, "￥"+gr.AddComma(cols[10])+"-")
		report.LineH(x, 74, x+72)
		report.GrayStroke(0)
		y = 85
	}
	report.LineType("straight", 5)
	report.GrayStroke(0.85)
	report.LineH(x, y, x+160)
	report.LineType("straight", 0.3)
	report.GrayStroke(0)
	report.Rect(x, y, x+160, y+5)
	report.LineV(x+23, y, y+5)
	report.LineV(x+41, y, y+5)
	report.LineV(x+61, y, y+5)
	report.LineV(x+99, y, y+5)
	report.LineV(x+116, y, y+5)
	report.LineV(x+135, y, y+5)
	report.Font("IPAexゴシック", 10, "")
	yadd := 1.5
	report.Cell(x+5, y+yadd, "年月日")
	report.Cell(x+28, y+yadd, "伝票")
	report.Cell(x+47, y+yadd, "品番")
	report.Cell(x+76, y+yadd, "品名")
	report.Cell(x+104, y+yadd, "数量")
	report.Cell(x+122, y+yadd, "単価")
	report.Cell(x+144, y+yadd, "金額")
	report.SumWork["g1item"] = 0.0
	report.SumWork["g2item"] = 1.0
}

type C2G1Summary struct {
}

func (h C2G1Summary) GetHeight(report gr.GoReport) float64 {
	return 5

}
func (h C2G1Summary) Execute(report gr.GoReport) {
	x := 25.0
	y := 1.5
	report.LineType("straight", 5)
	report.GrayStroke(0.85)
	report.LineH(x, 0, x+160)
	report.LineType("straight", 0.3)
	report.GrayStroke(0)
	report.Rect(x, 0, x+160, 5)
	report.Font("IPAexゴシック", 10, "")
	report.CellRight(x+159, y, 0, "\u00A5"+gr.AddComma(
		strconv.FormatFloat(report.SumWork["g1amtcum"], 'f', 0, 64)))
	report.Cell(x+117, y, "伝票合計")
	report.SumWork["g1amtcum"] = 0.0
	report.SumWork["g1item"] = 0.0

}

type C2G2Summary struct {
}

func (h C2G2Summary) GetHeight(report gr.GoReport) float64 {
	return 20
}
func (h C2G2Summary) Execute(report gr.GoReport) {
	x := 25.0
	y := 1.5
	report.LineType("straight", 15)
	report.GrayStroke(0.85)
	report.LineH(x+116, 0, x+160)
	report.LineType("straight", 0.3)
	report.GrayStroke(0)
	report.Rect(x+116, 0, x+160, 5)
	report.Rect(x+116, 5, x+160, 10)
	report.Rect(x+116, 10, x+160, 15)
	report.Font("IPAexゴシック", 10, "")
	report.CellRight(x+159, y, 0, "\u00A5"+gr.AddComma(
		strconv.FormatFloat(report.SumWork["g2amtcum"], 'f', 0, 64)))
	report.Cell(x+117, y, "合計")
	amt := report.SumWork["g2amtcum"]
	report.CellRight(x+159, y, 0, "\u00A5"+gr.AddComma(
		strconv.FormatFloat(amt, 'f', 0, 64)))
	cons := amt * 0.08
	report.Cell(x+117, y+5, "消費税(8%)")
	report.CellRight(x+159, y+5, 0, "\u00A5"+gr.AddComma(
		strconv.FormatFloat(cons, 'f', 0, 64)))
	report.Cell(x+117, y+10, "請求金額")
	report.CellRight(x+159, y+10, 0, "\u00A5"+gr.AddComma(
		strconv.FormatFloat(amt+cons, 'f', 0, 64)))
	report.SumWork["g2amtcum"] = 0.0
	report.SumWork["g2item"] = 0.0
	report.NewPage(true)
}

type C2Footer struct {
}

func (h C2Footer) GetHeight(report gr.GoReport) float64 {
	return 10
}
func (h C2Footer) Execute(report gr.GoReport) {
	report.Cell(100, 12, "Page")
	report.Cell(112, 12, strconv.Itoa(report.Page))
}
