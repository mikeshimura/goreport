package goreport

import (
	"fmt"
	"github.com/signintech/gopdf"
	"io/ioutil"
	"strconv"
	"strings"
)

type Converter struct {
	Pdf    *gopdf.GoPdf
	Text   string
	Fonts  []*FontMap
	ConvPt float64
}

//var p.ConvPt float64 = 2.834645669

//Read UTF-8 encoding file
func (p *Converter) ReadFile(fileName string) error {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	p.Text = strings.Replace(string(buf), "\r", "", -1)
	fmt.Println("text:" + p.Text)
	var UTF8_BOM = []byte{239, 187, 191}
	if p.Text[0:3] == string(UTF8_BOM) {
		p.Text = p.Text[3:]
	}
	return nil
}

func (p *Converter) Execute() {
	lines := strings.Split(p.Text, "\n")
	for _, line := range lines {
		//fmt.Println("line:" + line)
		eles := strings.Split(line, "\t")
		//fmt.Printf("eles[0]:%v:len %v\n",eles[0],len(eles[0]))
		switch eles[0] {
		case "P", "P1":
			p.Page(line, eles)
		case "NP":
			p.NewPage(line, eles)
		case "F":
			p.Font(line, eles)
		case "C", "C1", "CR":
			p.Cell(line, eles)
		case "M":
			p.Move(line, eles)
		case "L", "LV", "LH", "LT":
			p.Line(line, eles)
		case "R":
			p.Rect(line, eles)
		case "I":
			p.Image(line, eles)
		default:
			fmt.Println("default:" + line + ":")
		}
	}
}
func (p *Converter) AddFont() {
	for _, font := range p.Fonts {
		err := p.Pdf.AddTTFFont(font.FontName, font.FileName)
		if err != nil {
			panic(err)
		}
	}

}
func (p *Converter) Page(line string, eles []string) {
	p.Pdf = new(gopdf.GoPdf)
	switch eles[0] {
	case "P":
		CheckLength(line, eles, 4)
		switch eles[2] {
		case "A4":
			if eles[3] == "P" {
				p.Start(595.28, 841.89)
			} else if eles[3] == "L" {
				p.Start(841.89, 595.28)
			} else {
				panic("Page Orientation accept P or L")
			}
			p.SetConv(eles[1])
		default:
			panic("This size not supported yet:" + eles[2])
		}
	case "P1":
		CheckLength(line, eles, 4)
		p.SetConv(eles[1])
		p.Start(ParseFloatPanic(eles[2])*p.ConvPt,
			ParseFloatPanic(eles[3])*p.ConvPt)
	}
	p.AddFont()
	p.Pdf.AddPage()
}
func (p *Converter) SetConv(ut string) {
	switch ut {
	case "mm":
		p.ConvPt = 2.834645669
	case "pt":
		p.ConvPt = 1
	case "in":

		p.ConvPt = 72
	default:
		panic("This unit is not specified :" + ut)
	}
}
func (p *Converter) NewPage(line string, eles []string) {
	p.Pdf.AddPage()
}
func (p *Converter) Start(w float64, h float64) {
	p.Pdf.Start(gopdf.Config{Unit: "pt",
		PageSize: gopdf.Rect{W: w, H: h}}) //595.28, 841.89 = A4
}
func (p *Converter) Font(line string, eles []string) {
	CheckLength(line, eles, 4)
	err := p.Pdf.SetFont(eles[1], eles[2], AtoiPanic(eles[3]))
	if err != nil {
		panic(err)
	}
}
func (p *Converter) Rect(line string, eles []string) {
	CheckLength(line, eles, 5)
	p.Pdf.Line(ParseFloatPanic(eles[1])*p.ConvPt,
		ParseFloatPanic(eles[2])*p.ConvPt,
		ParseFloatPanic(eles[3])*p.ConvPt,
		ParseFloatPanic(eles[2])*p.ConvPt)
	p.Pdf.Line(ParseFloatPanic(eles[1])*p.ConvPt,
		ParseFloatPanic(eles[2])*p.ConvPt,
		ParseFloatPanic(eles[1])*p.ConvPt,
		ParseFloatPanic(eles[4])*p.ConvPt)
	p.Pdf.Line(ParseFloatPanic(eles[1])*p.ConvPt,
		ParseFloatPanic(eles[4])*p.ConvPt,
		ParseFloatPanic(eles[3])*p.ConvPt,
		ParseFloatPanic(eles[4])*p.ConvPt)
	p.Pdf.Line(ParseFloatPanic(eles[3])*p.ConvPt,
		ParseFloatPanic(eles[2])*p.ConvPt,
		ParseFloatPanic(eles[3])*p.ConvPt,
		ParseFloatPanic(eles[4])*p.ConvPt)
}
func (p *Converter) Image(line string, eles []string) {
	CheckLength(line, eles, 6)
	r := new(gopdf.Rect)
	r.W = ParseFloatPanic(eles[4])*p.ConvPt -
		ParseFloatPanic(eles[2])*p.ConvPt
	r.H = ParseFloatPanic(eles[5])*p.ConvPt -
		ParseFloatPanic(eles[3])*p.ConvPt
	p.Pdf.Image(eles[1], ParseFloatPanic(eles[2])*p.ConvPt,
		ParseFloatPanic(eles[3])*p.ConvPt, r)
}
func (p *Converter) Line(line string, eles []string) {
	switch eles[0] {
	case "L":
		CheckLength(line, eles, 5)
		p.Pdf.Line(ParseFloatPanic(eles[1])*p.ConvPt,
			ParseFloatPanic(eles[2])*p.ConvPt,
			ParseFloatPanic(eles[3])*p.ConvPt,
			ParseFloatPanic(eles[4])*p.ConvPt)
	case "LH":
		CheckLength(line, eles, 4)
		p.Pdf.Line(ParseFloatPanic(eles[1])*p.ConvPt,
			ParseFloatPanic(eles[2])*p.ConvPt,
			ParseFloatPanic(eles[3])*p.ConvPt,
			ParseFloatPanic(eles[2])*p.ConvPt)
	case "LV":
		CheckLength(line, eles, 4)
		p.Pdf.Line(ParseFloatPanic(eles[1])*p.ConvPt,
			ParseFloatPanic(eles[2])*p.ConvPt,
			ParseFloatPanic(eles[1])*p.ConvPt,
			ParseFloatPanic(eles[3])*p.ConvPt)
	case "LT":
		//lineType "dashed" ,"dotted"
		CheckLength(line, eles, 3)
		lineType := eles[1]
		if lineType == "" {
			lineType = "straight"
		}
		p.Pdf.SetLineType(lineType)
		p.Pdf.SetLineWidth(ParseFloatPanic(eles[2]))
	}

}
func CheckLength(line string, eles []string, no int) {
	if len(eles) < no {
		panic("Column short:" + line)
	}
}
func (p *Converter) Cell(line string, eles []string) {
	switch eles[0] {
	case "C":
		CheckLength(line, eles, 6)
		err := p.Pdf.SetFont(eles[1], "", AtoiPanic(eles[2]))
		if err != nil {
			panic(err)
		}
		p.MoveSub(eles[3], eles[4])
		p.Pdf.Cell(nil, eles[5])
	case "C1":
		CheckLength(line, eles, 4)
		p.MoveSub(eles[1], eles[2])
		p.Pdf.Cell(nil, eles[3])
	case "CR":
		CheckLength(line, eles, 5)
		tw, err := p.Pdf.MeasureTextWidth(eles[4])
		if err != nil {
			panic(err)
		}
		x := ParseFloatPanic(eles[1]) * p.ConvPt
		y := ParseFloatPanic(eles[2]) * p.ConvPt
		w := ParseFloatPanic(eles[3]) * p.ConvPt
		finalx := x + w - tw
		p.Pdf.SetX(finalx)
		p.Pdf.SetY(y)
		p.Pdf.Cell(nil, eles[4])
	}
}
func (p *Converter) Move(line string, eles []string) {
	CheckLength(line, eles, 3)
	p.MoveSub(eles[1], eles[2])
}
func (p *Converter) MoveSub(sx string, sy string) {
	p.Pdf.SetX(ParseFloatPanic(sx) * p.ConvPt)
	p.Pdf.SetY(ParseFloatPanic(sy) * p.ConvPt)
}
func AtoiPanic(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(s + " not Integer")
	}
	return i
}
func ParseFloatPanic(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(s + " not Numeric")
	}
	return f
}

type FontMap struct {
	FontName string
	FileName string
}
