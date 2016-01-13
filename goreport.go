package goreport

import (
	"bytes"
	"fmt"
	"github.com/mikeshimura/dbflute/df"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
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
	Records      []interface{}
	DataPos      int
	Bands        map[string]*Band
	Converter    *Converter
	PageX        float64
	PageY        float64
	CurrY        float64
	FooterY      float64
	MaxGroup     int
	Page         int
	PageTotal    bool
	NewPageForce bool
	ResetPageNo  bool
}

func (r *GoReport) SetFont(fmap []*FontMap) {
	r.Converter.Fonts = fmap
}

func (r *GoReport) NewPage(resetPageNo bool) {
	r.NewPageForce = true
	r.ResetPageNo = resetPageNo
}
func (r *GoReport) Execute(filename string) {
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
	r.Converter.Execute()
	r.Converter.Pdf.WritePdf(filename)
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
			h.page = AtoiPanic(line[7:])
			list.Add(h)
			fmt.Printf("hist %v \n", h)
		}
	}
	for i, line := range lines {
		if strings.Index(line, "{#TotalPage#}") > -1 {
			total := r.getTotalValue(i, list)
			fmt.Printf("total :%v\n", total)
			lines[i] = strings.Replace(lines[i], "{#TotalPage#}", strconv.Itoa(total), -1)
		}
	}
	buf := new(bytes.Buffer)
	for _, line := range lines {
		buf.WriteString(line + "\n")
	}
	r.Converter.Text = buf.String()
}
func (r *GoReport) getTotalValue(lineno int, list *df.List) int {
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
		fmt.Printf("page :%v\n", page)
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
	if r.CurrY+height > r.FooterY {
		r.PageBreak(false)
	}
}
func (r *GoReport) ExecutePageFooter() {
	r.CurrY = r.FooterY
	h := r.Bands[PageFooter]
	if h != nil {
		(*h).Execute()
		r.CurrY += (*h).GetHeight()
	}
}
func (r *GoReport) ExecuteSummary() {
	h := r.Bands[Summary]
	if h != nil {
		r.PageBreakCheck((*h).GetHeight())
		(*h).Execute()
		r.CurrY += (*h).GetHeight()
	}
}
func (r *GoReport) ExecutePageHeader() {
	h := r.Bands[PageHeader]
	if h != nil {
		(*h).Execute()
		r.CurrY += (*h).GetHeight()
	}
}
func (r *GoReport) ExecuteGroupHeader(level int) {
	for l := level; l > 0; l-- {
		h := r.Bands[GroupHeader+strconv.Itoa(l)]
		if h != nil {
			r.PageBreakCheck((*h).GetHeight())
			(*h).Execute()
			r.CurrY += (*h).GetHeight()
		}
	}
}
func (r *GoReport) ExecuteGroupSummary(level int) {
	for l := 1; l <= level; l++ {
		h := r.Bands[GroupSummary+strconv.Itoa(l)]
		if h != nil {
			r.PageBreakCheck((*h).GetHeight())
			(*h).Execute()
			r.CurrY += (*h).GetHeight()
		}
	}
}
func (r *GoReport) ExecuteDetail() {
	h := r.Bands[Detail]
	if h != nil {
		if r.NewPageForce {
			fmt.Println("NewPageForce")
			r.PageBreak(r.ResetPageNo)
			r.NewPageForce = false
			r.ResetPageNo = false
		}
		var deti interface{} = *h
		if r.MaxGroup > 0 {
			bfr := reflect.ValueOf(deti).MethodByName("BreakCheckBefore")
			if bfr.IsValid() == false {
				panic("BreakCheckBefore functionがDetailにありません")
			}
			res := bfr.Call([]reflect.Value{reflect.ValueOf(r)})
			level := res[0].Int()
			if level > 0 {
				r.ExecuteGroupHeader(int(level))
			}
		}
		r.PageBreakCheck((*h).GetHeight())
		(*h).Execute()
		r.CurrY += (*h).GetHeight()
		aft := reflect.ValueOf(deti).MethodByName("BreakCheckAfter")
		if aft.IsValid() == false {
			panic("BreakCheckAfter functionがDetailにありません")
		}
		res := aft.Call([]reflect.Value{reflect.ValueOf(r)})
		level := res[0].Int()
		if level > 0 {
			r.ExecuteGroupSummary(int(level))
		}
	}
}
func (r *GoReport) RegisterBand(band *Band, name string) {
	(*band).SetGoReport(r)
	r.Bands[name] = band
}
func (r *GoReport) RegisterGroupBand(band *Band, name string, level int) {
	(*band).SetGoReport(r)
	r.Bands[name+strconv.Itoa(level)] = band
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
func (r *GoReport) Var(name string, val string) {
	r.AddLine("V\t" + name + "\t" + val)
}
func (r *GoReport) SetPage(size string, orientation string) {
	switch size {
	case "A4":
		switch orientation {
		case "P":
			r.AddLine("P\tA4\tP")
			r.PageX = 210
			r.PageY = 297
		case "L":
			r.AddLine("P\tA4\tL")
			r.PageX = 297
			r.PageY = 210
		}
	}
}
func (r *GoReport) SaveDoc(fileName string) {
	ioutil.WriteFile(fileName, []byte(r.Converter.Text), os.ModePerm)
}

type Band interface {
	GetHeight() float64
	Execute()
	SetGoReport(r *GoReport)
}

func CreateGoReport() *GoReport {
	GoReport := new(GoReport)
	GoReport.Bands = make(map[string]*Band)
	GoReport.Converter = new(Converter)
	return GoReport
}
