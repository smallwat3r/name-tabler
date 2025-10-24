package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/phpdave11/gofpdf"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func main() {
	colsFlag := flag.Int("cols", 3, "number of columns per page (>=1)")
	rowsFlag := flag.Int("rows", 7, "number of rows per page (>=1)")
	perFlag := flag.Int("n-tab", 5, "names per table cell (>=1)")
	inFile := flag.String("file", "names.txt", "path to input names file")
	outFile := flag.String("out", "names.pdf", "output PDF file")
	flag.Parse()

	// clamp to minimums
	cols := max(1, *colsFlag)
	rows := max(1, *rowsFlag)
	namesPerTable := max(1, *perFlag)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	namesPath := *inFile
	if !filepath.IsAbs(namesPath) {
		namesPath = filepath.Join(cwd, namesPath)
	}

	names, err := readNames(namesPath)
	if err != nil {
		log.Fatal(err)
	}
	if len(names) == 0 {
		log.Fatal("names.txt is empty")
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	left, top, right, bottom := 10.0, 10.0, 10.0, 10.0
	pdf.SetMargins(left, top, right)
	pdf.SetAutoPageBreak(false, bottom)
	pdf.AddPage()

	pdf.SetFont("Helvetica", "", 16)

	pageW, pageH := pdf.GetPageSize()
	innerW := pageW - left - right
	innerH := pageH - top - bottom
	tableW := innerW / float64(cols)
	tableH := innerH / float64(rows)

	pad := 2.0
	lineH := (tableH - 2*pad) / float64(namesPerTable)

	pdf.SetLineWidth(0.2)
	pdf.SetDashPattern([]float64{2, 1}, 0)

	idx := 0
	for {
		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {
				x := left + float64(c)*tableW
				y := top + float64(r)*tableH

				pdf.Rect(x, y, tableW, tableH, "D")

				w := tableW - 2*pad
				curX := x + pad
				curY := y + pad
				pdf.SetXY(curX, curY)

				for i := 0; i < namesPerTable; i++ {
					var name string
					if idx < len(names) {
						name = utf8ToCP1252(names[idx]) // allow accents
						idx++
					}
					pdf.MultiCell(w, lineH, name, "", "C", false)
					_, cy := pdf.GetXY()
					pdf.SetXY(curX, cy)
				}
			}
		}
		if idx >= len(names) {
			break
		}
		pdf.AddPage()
	}

	out := *outFile
	if !filepath.IsAbs(out) {
		out = filepath.Join(cwd, out)
	}
	if err := pdf.OutputFileAndClose(out); err != nil {
		log.Fatal(err)
	}
	fmt.Println("PDF generated:", out)
}

func readNames(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []string
	sc := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 1<<20)

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line != "" {
			out = append(out, line)
		}
	}
	return out, sc.Err()
}

// convert UTF-8 to Windows-1252 so accented Latin characters display correctly in built-in Helvetica
func utf8ToCP1252(s string) string {
	r := transform.NewReader(strings.NewReader(s), charmap.Windows1252.NewEncoder())
	b, err := io.ReadAll(r)
	if err != nil {
		return s // fallback on error
	}
	return string(b)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
