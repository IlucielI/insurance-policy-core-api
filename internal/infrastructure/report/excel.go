package report

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExcelGenerator struct{}

func NewExcelGenerator() *ExcelGenerator {
	return &ExcelGenerator{}
}

// ClaimsReportExcel generates claims report Excel sheet
func (g *ExcelGenerator) ClaimsReportExcel(rows []ClaimsReportRow, filters map[string]string) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Laporan Klaim"
	f.SetSheetName("Sheet1", sheet)

	headerStyle := g.makeStyle(f, &excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"0D6EFD"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border:    []excelize.Border{{Type: "bottom", Color: "212529", Style: 1}},
	})

	dataStyle := g.makeStyle(f, &excelize.Style{
		Font:      &excelize.Font{Size: 10},
		Alignment: &excelize.Alignment{Vertical: "center"},
		Border:    []excelize.Border{{Type: "bottom", Color: "DEE2E6", Style: 1}},
	})

	titleStyle := g.makeStyle(f, &excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14, Color: "212529"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	// title
	title := "LAPORAN KLAIM"
	if fv, ok := filters["status"]; ok && fv != "" {
		title += fmt.Sprintf(" — Status: %s", fv)
	}
	f.SetCellValue(sheet, "A1", title)
	f.MergeCell(sheet, "A1", "H1")
	f.SetRowHeight(sheet, 1, 30)
	f.SetCellStyle(sheet, "A1", "H1", titleStyle)

	// subtitle
	sub := fmt.Sprintf("Tanggal Ekspor: %s | Total Data: %d", time.Now().Format("02 Jan 2006"), len(rows))
	f.SetCellValue(sheet, "A2", sub)
	f.MergeCell(sheet, "A2", "H2")
	subStyle := g.makeStyle(f, &excelize.Style{
		Font:      &excelize.Font{Size: 9, Color: "6C757D"},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	f.SetCellStyle(sheet, "A2", "H2", subStyle)

	// headers at row 4
	headers := []string{"No", "Nomor Klaim", "Nomor Polis", "Tipe", "Jumlah Klaim", "Status", "Tanggal Insiden", "Tanggal Submit"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 4)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}
	f.SetRowHeight(sheet, 4, 22)

	// column widths
	colWidths := map[string]float64{"A": 6, "B": 28, "C": 24, "D": 14, "E": 20, "F": 14, "G": 18, "H": 18}
	for col, w := range colWidths {
		f.SetColWidth(sheet, col, col, w)
	}

	// data rows starting at row 5
	for i, r := range rows {
		row := i + 5
		vals := []interface{}{r.No, r.ClaimNumber, r.PolicyNum, r.ClaimType, rupiah(r.ClaimAmount), r.Status, r.IncidentDate, r.SubmittedAt}
		for j, v := range vals {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			f.SetCellValue(sheet, cell, v)
			f.SetCellStyle(sheet, cell, cell, dataStyle)
		}
		f.SetRowHeight(sheet, row, 20)
	}

	// summary
	if len(rows) > 0 {
		summaryRow := len(rows) + 6
		var totalAmount int64
		for _, r := range rows {
			totalAmount += r.ClaimAmount
		}
		boldStyle := g.makeStyle(f, &excelize.Style{
			Font:   &excelize.Font{Bold: true, Size: 10},
			Border: []excelize.Border{{Type: "top", Color: "212529", Style: 2}},
		})
		dCell, _ := excelize.CoordinatesToCellName(4, summaryRow)
		eCell, _ := excelize.CoordinatesToCellName(5, summaryRow)
		f.SetCellValue(sheet, dCell, "Total")
		f.SetCellValue(sheet, eCell, rupiah(totalAmount))
		f.SetCellStyle(sheet, dCell, eCell, boldStyle)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("excel write: %w", err)
	}
	return buf.Bytes(), nil
}

// CustomerListExcel generates customer list Excel
func (g *ExcelGenerator) CustomerListExcel(rows []CustomerRow) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Daftar Pelanggan"
	f.SetSheetName("Sheet1", sheet)

	headerStyle := g.makeStyle(f, &excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"0D6EFD"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	dataStyle := g.makeStyle(f, &excelize.Style{
		Font:      &excelize.Font{Size: 10},
		Alignment: &excelize.Alignment{Vertical: "center"},
		Border:    []excelize.Border{{Type: "bottom", Color: "DEE2E6", Style: 1}},
	})

	titleStyle := g.makeStyle(f, &excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14, Color: "212529"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	f.SetCellValue(sheet, "A1", "DAFTAR PELANGGAN")
	f.MergeCell(sheet, "A1", "F1")
	f.SetRowHeight(sheet, 1, 30)
	f.SetCellStyle(sheet, "A1", "F1", titleStyle)

	sub := fmt.Sprintf("Tanggal Ekspor: %s | Total: %d", time.Now().Format("02 Jan 2006"), len(rows))
	f.SetCellValue(sheet, "A2", sub)
	f.MergeCell(sheet, "A2", "F2")

	headers := []string{"No", "Nama", "Email", "Telepon", "Role", "Tanggal Daftar"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 4)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	colWidths := map[string]float64{"A": 6, "B": 30, "C": 28, "D": 18, "E": 18, "F": 20}
	for col, w := range colWidths {
		f.SetColWidth(sheet, col, col, w)
	}

	for i, r := range rows {
		row := i + 5
		vals := []interface{}{r.No, r.FullName, r.Email, r.Phone, r.Role, r.CreatedAt}
		for j, v := range vals {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			f.SetCellValue(sheet, cell, v)
			f.SetCellStyle(sheet, cell, cell, dataStyle)
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("excel write: %w", err)
	}
	return buf.Bytes(), nil
}

// AnalyticsSummaryExcel generates analytics summary Excel
func (g *ExcelGenerator) AnalyticsSummaryExcel(summary *AnalyticsSummary) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Ringkasan Analitik"
	f.SetSheetName("Sheet1", sheet)

	titleStyle := g.makeStyle(f, &excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14, Color: "212529"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	sectionStyle := g.makeStyle(f, &excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12, Color: "0D6EFD"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"F8F9FA"}, Pattern: 1},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	headerStyle := g.makeStyle(f, &excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"0D6EFD"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	dataStyle := g.makeStyle(f, &excelize.Style{
		Font:      &excelize.Font{Size: 10},
		Alignment: &excelize.Alignment{Vertical: "center"},
		Border:    []excelize.Border{{Type: "bottom", Color: "DEE2E6", Style: 1}},
	})

	f.SetCellValue(sheet, "A1", "RINGKASAN ANALITIK BISNIS")
	f.MergeCell(sheet, "A1", "E1")
	f.SetRowHeight(sheet, 1, 30)
	f.SetCellStyle(sheet, "A1", "E1", titleStyle)

	sub := fmt.Sprintf("Tanggal: %s", time.Now().Format("02 Jan 2006"))
	f.SetCellValue(sheet, "A2", sub)
	f.MergeCell(sheet, "A2", "E2")

	row := 4
	f.SetColWidth(sheet, "A", "A", 30)
	f.SetColWidth(sheet, "B", "E", 16)

	// Overview
	g.writeSection(f, sheet, &row, sectionStyle, "IKHTISAR")
	overview := [][2]string{
		{"Total Polis", fmt.Sprintf("%d", summary.TotalPolicies)},
		{"Total Pelanggan", fmt.Sprintf("%d", summary.TotalCustomers)},
		{"Total Klaim", fmt.Sprintf("%d", summary.TotalClaims)},
		{"Total Tagihan", fmt.Sprintf("%d", summary.TotalInvoices)},
	}
	for _, kv := range overview {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), kv[0])
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), kv[1])
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), dataStyle)
		row++
	}
	row++

	// Policies by Status
	g.writeSection(f, sheet, &row, sectionStyle, "POLIS BERDASARKAN STATUS")
	g.writeHeaders(f, sheet, &row, headerStyle, "Status", "Jumlah")
	for _, ps := range summary.PoliciesByStatus {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), ps.Status)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), ps.Count)
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("B%d", row), dataStyle)
		row++
	}
	row++

	// Claims by Status
	g.writeSection(f, sheet, &row, sectionStyle, "KLAIM BERDASARKAN STATUS")
	g.writeHeaders(f, sheet, &row, headerStyle, "Status", "Jumlah")
	for _, cs := range summary.ClaimsByStatus {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), cs.Status)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), cs.Count)
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("B%d", row), dataStyle)
		row++
	}
	row++

	// Revenue
	g.writeSection(f, sheet, &row, sectionStyle, "RINGKASAN PENDAPATAN")
	revenue := [][2]string{
		{"Total Premi (Semua Polis)", rupiah(summary.TotalPremium)},
		{"Total Klaim Dibayar", rupiah(summary.TotalPaidClaims)},
		{"Total Tagihan", rupiah(summary.TotalBilling)},
	}
	for _, kv := range revenue {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), kv[0])
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), kv[1])
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("B%d", row), dataStyle)
		row++
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("excel write: %w", err)
	}
	return buf.Bytes(), nil
}

func (g *ExcelGenerator) makeStyle(f *excelize.File, style *excelize.Style) int {
	id, err := f.NewStyle(style)
	if err != nil {
		return 0
	}
	return id
}

func (g *ExcelGenerator) writeSection(f *excelize.File, sheet string, row *int, styleID int, title string) {
	f.SetCellValue(sheet, fmt.Sprintf("A%d", *row), title)
	f.MergeCell(sheet, fmt.Sprintf("A%d", *row), fmt.Sprintf("E%d", *row))
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", *row), fmt.Sprintf("E%d", *row), styleID)
	*row++
}

func (g *ExcelGenerator) writeHeaders(f *excelize.File, sheet string, row *int, styleID int, cols ...string) {
	for i, col := range cols {
		cell, _ := excelize.CoordinatesToCellName(i+1, *row)
		f.SetCellValue(sheet, cell, col)
		f.SetCellStyle(sheet, cell, cell, styleID)
	}
	*row++
}

// AnalyticsSummary holds aggregated analytics data
type AnalyticsSummary struct {
	TotalPolicies    int
	TotalCustomers   int
	TotalClaims      int
	TotalInvoices    int
	TotalPremium     int64
	TotalBilling     int64
	TotalPaidClaims  int64
	PoliciesByStatus []StatusCount
	ClaimsByStatus   []StatusCount
}

type StatusCount struct {
	Status string
	Count  int
}

// ClaimsReportRow is flat struct for claims export
type ClaimsReportRow struct {
	No           int
	ClaimNumber  string
	PolicyNum    string
	ClaimType    string
	ClaimAmount  int64
	Status       string
	IncidentDate string
	SubmittedAt  string
}

// CustomerRow is flat struct for customer export
type CustomerRow struct {
	No        int
	FullName  string
	Email     string
	Phone     string
	Role      string
	CreatedAt string
}
