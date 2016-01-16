# goreport
Golang Pdf Report Generator  
[日本語の説明はこちら](https://github.com/mikeshimura/dbflute/wiki/Tutorial)

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
Second ster: Generate Pdf from Text data.

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



##Sample program
