package goreport

import (
	"testing"
)

func TestConverter(t *testing.T) {
	conv := new(Converter)
	font1 := FontMap{
		FontName: "IPAexゴシック",
		FileName: "example//ttf//ipaexg.ttf",
	}
	conv.Fonts = []*FontMap{&font1}
	conv.ReadFile("savetext.txt")
	conv.Execute()
	conv.Pdf.WritePdf("ConverterTest.pdf")
}
