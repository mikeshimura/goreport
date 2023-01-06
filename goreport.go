package goreport

import (
	"bytes"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/mikeshimura/dbflute/df"
)

const (
	PageHeader   = "PageHeader"
	PageFooter   = "PageFooter"
	Detail       = "Detail"
	Summary      = "Summary"
	GroupHeader  = "GroupHeader"
	GroupSummary = "GroupSummary"
)

type GoReport struct {
	Records   []interface{}
	DataPos   int
	Bands     map[string]*Band
	Converter *Converter
	ConvPt    float64
	PageX     float64
	PageY     float64
	CurrY     float64
	MaxGroup  int
	Page      int
	PageTotal bool
	Flags     map[string]bool
	SumWork   map[string]float64
	Vars      map[string]string
}

func (r *GoReport) SetFooterYbyFooterHeight(footerHeight float64) {
	if r.PageY == 0 {
		panic("Page size not yet specified")
	}
	r.SumWork["__ft__"] = r.PageY - footerHeight
}
func (r *GoReport) SetFooterY(footerY float64) {
	r.SumWork["__ft__"] = footerY
}
func (r *GoReport) SetFonts(fmap []*FontMap) {
	r.Converter.Fonts = fmap
}

func (r *GoReport) NewPage(resetPageNo bool) {
	//fmt.Println("report.NewPage goRep")
	r.Flags["NewPageForce"] = true
	r.Flags["ResetPageNo"] = resetPageNo
}
func (r *GoReport) Convert(exec bool) {
	if exec {
		if r.SumWork["__ft__"] == 0 {
			panic("footerY not set yet.")
		}
		r.Page = 1
		r.CurrY = 0
		r.ExecutePageHeader()
		r.AddLine("v\tPAGE\t" + strconv.Itoa(r.Page))
		for r.DataPos = 0; r.DataPos < len(r.Records); r.DataPos++ {
			r.ExecuteDetail()
		}
		r.ExecuteSummary()
		r.ExecutePageFooter()
		r.ReplacePageTotal()
	}
	r.Converter.Execute()
}
func (r *GoReport) Execute(filename string) {
	r.Convert(true)
	r.Converter.Pdf.WritePdf(filename)
}

func (r *GoReport) GetBytesPdf() (ret []byte) {
	r.Convert(true)
	ret = r.Converter.Pdf.GetBytesPdf()
	return
}

func (r *GoReport) ReplacePageTotal() {
	if r.PageTotal == false {
		return
	}
	lines := strings.Split(r.Converter.Text, "\n")
	list := new(df.List)
	for i, line := range lines {
		if len(line) < 8 {
			continue
		}
		if line[0:7] == "v\tPAGE\t" {
			h := new(pagehist)
			h.line = i
			h.page = AtoiPanic(line[7:], line)
			list.Add(h)
			//fmt.Printf("hist %v \n", h)
		}
	}
	for i, line := range lines {
		if strings.Index(line, "{#TotalPage#}") > -1 {
			total := r.getTotalPage(i, list)
			//fmt.Printf("total :%v\n", total)
			lines[i] = strings.Replace(lines[i], "{#TotalPage#}", strconv.Itoa(total), -1)
		}
	}
	buf := new(bytes.Buffer)
	for _, line := range lines {
		buf.WriteString(line + "\n")
	}
	r.Converter.Text = buf.String()
}
func (r *GoReport) getTotalPage(lineno int, list *df.List) int {
	count := 0
	page := 0
	for i, l := range list.GetAsArray() {
		if l.(*pagehist).line >= lineno {
			count = i
			break
		}
	}
	for i := count; i < list.Size(); i++ {
		newpage := list.Get(i).(*pagehist).page
		if newpage <= page {
			return page
		}
		page = newpage
		//fmt.Printf("page :%v\n", page)
	}
	return page
}

type pagehist struct {
	line int
	page int
}

func (r *GoReport) PageBreak(resetPageNo bool) {
	r.ExecutePageFooter()
	r.AddLine("NP")
	if resetPageNo {
		r.Page = 1
	} else {
		r.Page++
	}
	r.CurrY = 0
	r.ExecutePageHeader()
	r.AddLine("v\tPAGE\t" + strconv.Itoa(r.Page))
}
func (r *GoReport) PageBreakCheck(height float64) {
	if r.CurrY+height > r.SumWork["__ft__"] {
		r.PageBreak(false)
	}
}
func (r *GoReport) ExecutePageFooter() {
	r.CurrY = r.SumWork["__ft__"]
	h := r.Bands[PageFooter]
	if h != nil {
		height := (*h).GetHeight(*r)
		(*h).Execute(*r)
		r.CurrY += height
	}
}
func (r *GoReport) ExecuteSummary() {
	h := r.Bands[Summary]
	if h != nil {
		height := (*h).GetHeight(*r)
		r.PageBreakCheck(height)
		(*h).Execute(*r)
		r.CurrY += height
	}
}
func (r *GoReport) ExecutePageHeader() {
	h := r.Bands[PageHeader]
	if h != nil {
		height := (*h).GetHeight(*r)
		(*h).Execute(*r)
		r.CurrY += height
	}
}
func (r *GoReport) ExecuteGroupHeader(level int) {
	for l := level; l > 0; l-- {
		h := r.Bands[GroupHeader+strconv.Itoa(l)]
		if h != nil {
			height := (*h).GetHeight(*r)
			r.PageBreakCheck(height)
			(*h).Execute(*r)
			r.CurrY += height
		}
	}
}
func (r *GoReport) ExecuteGroupSummary(level int) {
	for l := 1; l <= level; l++ {
		h := r.Bands[GroupSummary+strconv.Itoa(l)]
		if h != nil {
			height := (*h).GetHeight(*r)
			r.PageBreakCheck(height)
			(*h).Execute(*r)
			r.CurrY += height
		}
	}
}
func (r *GoReport) ExecuteDetail() {
	h := r.Bands[Detail]
	if h != nil {
		//fmt.Printf("report.NewPage flag %v\n", r.Flags["NewPageForce"])
		if r.Flags["NewPageForce"] {
			//fmt.Println("NewPageForce")
			r.PageBreak(r.Flags["ResetPageNo"])
			r.Flags["NewPageForce"] = false
			r.Flags["ResetPageNo"] = false
		}
		var deti interface{} = *h
		if r.MaxGroup > 0 {
			bfr := reflect.ValueOf(deti).MethodByName("BreakCheckBefore")
			if bfr.IsValid() == false {
				panic("BreakCheckBefore function not exist in Detail")
			}
			res := bfr.Call([]reflect.Value{reflect.ValueOf(*r)})
			level := res[0].Int()
			if level > 0 {
				r.ExecuteGroupHeader(int(level))
			}
		}
		height := (*h).GetHeight(*r)
		r.PageBreakCheck(height)
		(*h).Execute(*r)
		r.CurrY += height
		if r.MaxGroup > 0 {
			aft := reflect.ValueOf(deti).MethodByName("BreakCheckAfter")
			if aft.IsValid() == false {
				panic("BreakCheckAfter function not exist in Detail")
			}
			res := aft.Call([]reflect.Value{reflect.ValueOf(*r)})
			level := res[0].Int()
			if level > 0 {
				r.ExecuteGroupSummary(int(level))
			}
		}
	}
}
func (r *GoReport) RegisterBand(band Band, name string) {
	r.Bands[name] = &band
}
func (r *GoReport) RegisterGroupBand(band Band, name string, level int) {
	r.Bands[name+strconv.Itoa(level)] = &band
	if r.MaxGroup < level {
		r.MaxGroup = level
	}
}

func Ftoa(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}
func (r *GoReport) AddLine(s string) {
	r.Converter.Text += s + "\n"
}
func (r *GoReport) Font(fontName string, size int, style string) {
	r.AddLine("F\t" + fontName + "\t" + style + "\t" + strconv.Itoa(size))
}
func (r *GoReport) Cell(x float64, y float64, content string) {
	r.AddLine("C1\t" + Ftoa(x) + "\t" + Ftoa(r.CurrY+y) + "\t" + content)
}
func (r *GoReport) CellRight(x float64, y float64, w float64, content string) {
	r.AddLine("CR\t" + Ftoa(x) + "\t" + Ftoa(r.CurrY+y) + "\t" +
		Ftoa(w) + "\t" + content)
}

func (r *GoReport) LineType(ltype string, width float64) {
	r.SumWork["__lw__"] = width
	r.AddLine("LT\t" + ltype + "\t" + Ftoa(width))
}
func (r *GoReport) Line(x1 float64, y1 float64, x2 float64, y2 float64) {
	r.AddLine("L\t" + Ftoa(x1) + "\t" + Ftoa(r.CurrY+y1) + "\t" + Ftoa(x2) +
		"\t" + Ftoa(r.CurrY+y2))
}
func (r *GoReport) LineH(x1 float64, y float64, x2 float64) {
	adj := r.SumWork["__lw__"] * 0.5
	r.AddLine("LH\t" + Ftoa(x1) + "\t" + Ftoa(r.CurrY+y+adj) + "\t" + Ftoa(x2))
}
func (r *GoReport) LineV(x float64, y1 float64, y2 float64) {
	adj := r.SumWork["__lw__"] * 0.5
	r.AddLine("LV\t" + Ftoa(x+adj) + "\t" + Ftoa(r.CurrY+y1) + "\t" + Ftoa(r.CurrY+y2))
}

// SumWork["__lw__"] width adjust
func (r *GoReport) Rect(x1 float64, y1 float64, x2 float64, y2 float64) {
	r.AddLine("R\t" + Ftoa(x1) + "\t" + Ftoa(r.CurrY+y1) + "\t" + Ftoa(x2) +
		"\t" + Ftoa(r.CurrY+y2))
}
func (r *GoReport) Oval(x1 float64, y1 float64, x2 float64, y2 float64) {
	r.AddLine("O\t" + Ftoa(x1) + "\t" + Ftoa(r.CurrY+y1) + "\t" + Ftoa(x2) +
		"\t" + Ftoa(r.CurrY+y2))
}
func (r *GoReport) TextColor(red int, green int, blue int) {
	r.AddLine("TC\t" + strconv.Itoa(red) + "\t" + strconv.Itoa(green) +
		"\t" + strconv.Itoa(blue))
}
func (r *GoReport) StrokeColor(red int, green int, blue int) {
	r.AddLine("SC\t" + strconv.Itoa(red) + "\t" + strconv.Itoa(green) +
		"\t" + strconv.Itoa(blue))
}
func (r *GoReport) GrayFill(grayScale float64) {
	r.AddLine("GF\t" + Ftoa(grayScale))
}
func (r *GoReport) GrayStroke(grayScale float64) {
	r.AddLine("GS\t" + Ftoa(grayScale))
}
func (r *GoReport) Image(path string, x1 float64, y1 float64, x2 float64, y2 float64) {
	r.AddLine("I\t" + path + "\t" + Ftoa(x1) + "\t" + Ftoa(r.CurrY+y1) + "\t" +
		Ftoa(x2) + "\t" + Ftoa(r.CurrY+y2))
}

func (r *GoReport) Var(name string, val string) {
	r.AddLine("V\t" + name + "\t" + val)
}

// unit mm pt in  size A4 LTR
func (r *GoReport) SetPage(size string, unit string, orientation string) {
	r.SetConv(unit)
	switch size {
	case "A4":
		switch orientation {
		case "P":
			r.AddLine("P\t" + unit + "\tA4\tP")
			r.PageX = 595.28 / r.ConvPt
			r.PageY = 841.89 / r.ConvPt
		case "L":
			r.AddLine("P\t" + unit + "\tA4\tL")
			r.PageX = 841.89 / r.ConvPt
			r.PageY = 595.28 / r.ConvPt
		}
	case "LTR":
		switch orientation {
		case "P":

			r.PageX = 612 / r.ConvPt
			r.PageY = 792 / r.ConvPt
			r.AddLine("P1\t" + unit + "\t" + strconv.FormatFloat(r.PageX, 'f', 4, 64) +
				"\t" + strconv.FormatFloat(r.PageY, 'f', 4, 64))
		case "L":

			r.PageX = 792 / r.ConvPt
			r.PageY = 612 / r.ConvPt
			r.AddLine("P1\t" + unit + "\t" + strconv.FormatFloat(r.PageX, 'f', 4, 64) +
				"\t" + strconv.FormatFloat(r.PageY, 'f', 4, 64))
		}
	}
	r.Convert(false)
}

// unit mm pt in
func (r *GoReport) SetConv(ut string) {
	switch ut {
	case "mm":
		r.ConvPt = 2.834645669
	case "pt":
		r.ConvPt = 1

	case "in":

		r.ConvPt = 72
	default:
		panic("This unit is not specified :" + ut)
	}
}
func (r *GoReport) SetPageByDimension(unit string, width float64, height float64) {
	r.SetConv(unit)
	r.PageX = width
	r.PageY = height
	r.AddLine("P1\t" + unit + "\t" + strconv.FormatFloat(width, 'f', 4, 64) +
		"\t" + strconv.FormatFloat(height, 'f', 4, 64))
	r.Convert(false)
}

func (r *GoReport) SaveText(fileName string) {
	os.WriteFile(fileName, []byte(r.Converter.Text), os.ModePerm)
}

type Band interface {
	GetHeight(report GoReport) float64
	Execute(report GoReport)
}

func CreateGoReport() *GoReport {
	GoReport := new(GoReport)
	GoReport.Bands = make(map[string]*Band)
	GoReport.Converter = new(Converter)
	GoReport.SumWork = make(map[string]float64)
	GoReport.Vars = make(map[string]string)
	GoReport.SumWork["__ft__"] = 0.0 //FooterY
	GoReport.Flags = make(map[string]bool)
	GoReport.Flags["NewPageForce"] = false
	GoReport.Flags["ResetPageNo"] = false
	return GoReport
}

type TemplateDetail struct {
}

func (h TemplateDetail) GetHeight() float64 {
	return 10
}
func (h TemplateDetail) Execute(report GoReport) {
}
