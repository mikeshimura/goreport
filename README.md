# goreport
Golang Pdf Report Generator  
[日本語の説明はこちら](https://github.com/mikeshimura/goreport/wiki/%E6%97%A5%E6%9C%AC%E8%AA%9E%E8%AA%AC%E6%98%8E)

This program use github.com/signintech/gopdf for pdf generation.

Simple report sample
![Simple1](https://bytebucket.org/mikeshimura/goreport/wiki/image/simple1.jpg "Simple1")  

[pdf](https://bytebucket.org/mikeshimura/goreport/wiki/pdf/simple1.pdf)  


Medium report sample
![Medium1](https://bytebucket.org/mikeshimura/goreport/wiki/image/medium1.jpg "Medium1")

[pdf](https://bytebucket.org/mikeshimura/goreport/wiki/pdf/medium1.pdf)  



Complex report sample
![Complex1](https://bytebucket.org/mikeshimura/goreport/wiki/image/complex1.jpg "Complex1")

[pdf](https://bytebucket.org/mikeshimura/goreport/wiki/pdf/complex1.pdf)  


![Complex2](https://bytebucket.org/mikeshimura/goreport/wiki/image/complex2.jpg "Complex2")

[pdf](https://bytebucket.org/mikeshimura/goreport/wiki/pdf/complex2.pdf)  

##Installation
```
go get -u github.com/mikeshimura/goreport
```
##concept
- Following Bands are available.  
PageHeader  
GroupHeader2  
GroupHeader1  
Detail  
GroupSummary1  
GroupSummry2  
Summary  
PageFooter

- Groups can be any number

- User defined Band structure required to implement Band interface.  
Only two functions are required.


```
GetHeight(report GoReport) float64
Execute(report GoReport)
```

- Two step executiion.  
First step: Generate Text data.  
Second step: Generate Pdf from Text data.

- Above two step execution enable very flexible usability.  
You may generate Text data by program, then any kind of pdf can be generated.

- I use above flexibity to insert total pages data after generation automatically.

- Band height can be changed program, therefore conditional display can be achieved.

- Data source is stored as []interface{}, then any kind of data type can be used. For example, string array, entity object, map etc.

- Any Ttf Font can be used

##Setup Commands
- Font Setting Sample
```
font1 := gr.FontMap{
		FontName: "IPAex",
		FileName: "ttf//ipaexg.ttf",
	}
fonts := []*gr.FontMap{&font1}
r.SetFonts(fonts)
```
- Page Setting  
 SetPage(size string, unit string, orientation string)  
 //size A4 or LTR, unit mm, pt or in  

 SetPageByDimension(unit string, width float64, height float64)  
- Normal Print limit setting. PageFooter(if defined) will be written after reach this limt.  
 SetFooterY(footerY float64)  

 SetFooterYbyFooterHeight(footerHeight float64)  
 //Sheet height - footerHeight will be set  

##Draw Commands

- Font setting  
Font(fontName string, size int, style string)  
//style "" or "U" (underline)  

 TextColor(red int, green int, blue int) //Set Font color  
 GrayFill(grayScale float64) //Set grayScale for black font  

- Text Draw  
 Cell(x float64, y float64, content string)  
 CellRight(x float64, y float64, w float64, content string)  //Right Justify  

- Line Draw  
 LineType(ltype string, width float64)
 //lineType "dashed" ,"dotted","straight" ""="straight"  
 GrayStroke(grayScale float64)  //Set grayScale  
 LineH(x1 float64, y float64, x2 float64) // Horizontal Line  
 LineV(x float64, y1 float64, y2 float64) // Vertical line  
 Line(x1 float64, y1 float64, x2 float64, y2 float64)  

- Shape Draw  
 Rect(x1 float64, y1 float64, x2 float64, y2 float64)  
 Oval(x1 float64, y1 float64, x2 float64, y2 float64)  

- Image  Draw  
  Image(path string, x1 float64, y1 float64, x2 float64, y2 float64)  
	
##Genarate Commands
-  Execute(filename string)  
Genarate PDF File.

- GetBytesPdf() (ret []byte)  
Create byte stream

##License  

goreport is released under the MIT License. It is copyrighted by Masanobu Shimura. (Gmail mikeshimura)

##Limitation  

- Font style not allow B(bold) and I(italic).
- Line, Rect and Oval are Black and Gray only.

##Sample program

[sample source](https://github.com/mikeshimura/goreport/tree/master/example)
```go
package example

import (
	gr "github.com/mikeshimura/goreport"
	"strconv"
)

func Simple1Report() {
	r := gr.CreateGoReport()
	//var accumrate amount
	r.SumWork["amountcum="]=0.0
	font1 := gr.FontMap{
		FontName: "IPAex",
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
	r.Records = gr.ReadTextFile("sales1.txt",7)
	r.SetPage("A4", "mm","L")
	r.SetFooterY(190)
	r.Execute("simple1.pdf")
}

type S1Detail struct {
}

func (h S1Detail) GetHeight(report gr.GoReport) float64 {
	return 10
}
func (h S1Detail) Execute(report gr.GoReport) {
	cols := report.Records[report.DataPos].([]string)
	report.Font("IPAex", 12, "")
	y:=2.0
	report.Cell(15, y, cols[0])
	report.Cell(30, y, cols[1])
	report.Cell(60, y, cols[2])
	report.Cell(90, y, cols[3])
	report.Cell(120, y, cols[4])
	report.CellRight(135, y,25, cols[5])
	report.CellRight(160, y,20, cols[6])
	amt:=ParseFloatNoError(cols[5])*ParseFloatNoError(cols[6])
	report.SumWork["amountcum="]+=amt
	report.CellRight(180, y,30, strconv.FormatFloat(amt,'f',2,64))
}

type S1Header struct {
}

func (h S1Header) GetHeight(report gr.GoReport) float64 {
	return 30
}
func (h S1Header) Execute(report gr.GoReport) {
	report.Font("IPAex", 14, "")
	report.Cell(50, 15, "Sales Report")
	report.Font("", 12, "")
	report.Cell(240, 20, "page")
	report.Cell(260, 20, strconv.Itoa(report.Page))
	y:=23.0
	report.Cell(15, y, "D No")
	report.Cell(30, y, "Dept")
	report.Cell(60, y, "Order")
	report.Cell(90, y, "Stock")
	report.Cell(120, y,"Name")
	report.CellRight(135, y,25, "Unit Price")
	report.CellRight(160, y,20, "Qty")
	report.CellRight(190, y,20, "Amount")
}

type S1Summary struct {
}

func (h S1Summary) GetHeight(report gr.GoReport) float64 {
	return 10
}
func (h S1Summary) Execute(report gr.GoReport) {
	report.Cell(160, 2,"Total")
	report.CellRight(180, 2,30, strconv.FormatFloat(
			report.SumWork["amountcum="],'f',2,64))
}

func ParseFloatNoError(s string)float64{
	f,_:=strconv.ParseFloat(s,64)
	return f
}
```
