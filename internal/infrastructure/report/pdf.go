package report

import (
	"fmt"
	"strings"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/jung-kurt/gofpdf"
)

const (
	margin    = 20.0
	pageWidth = 210.0
	bodyWidth = pageWidth - 2*margin
	colRatio  = 0.38
)

var (
	clrDark    = []int{33, 37, 41}
	clrPrimary = []int{13, 110, 253}
	clrLightBg = []int{248, 249, 250}
	clrWhite   = []int{255, 255, 255}
)

type PDFGenerator struct{}

func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{}
}

// BillingStatementPDF generates a professional billing statement PDF
func (g *PDFGenerator) BillingStatementPDF(invoice *domain.Invoice, user *domain.User, policy *domain.Policy) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(margin, margin, margin)
	pdf.SetAutoPageBreak(true, 25)

	// header
	g.drawHeader(pdf, "Laporan Tagihan / Billing Statement")
	g.drawCompanyInfo(pdf)
	pdf.Ln(8)

	// invoice meta
	g.drawSectionTitle(pdf, "Detail Tagihan")
	g.startTable(pdf)
	g.tableRow(pdf, "Nomor Faktur", invoice.InvoiceNumber)
	g.tableRow(pdf, "Status", strings.ToUpper(invoice.Status))
	g.tableRow(pdf, "Jenis Tagihan", invoice.InvoiceType)
	g.tableRow(pdf, "Jatuh Tempo", invoice.DueDate)
	g.tableRow(pdf, "Tanggal Terbit", invoice.CreatedAt.Format("02 Jan 2006"))
	if invoice.PaidAt != nil {
		g.tableRow(pdf, "Tanggal Bayar", invoice.PaidAt.Format("02 Jan 2006"))
		g.tableRow(pdf, "Metode Bayar", invoice.PaymentMethod)
	}
	g.endTable(pdf)
	pdf.Ln(6)

	// customer info
	g.drawSectionTitle(pdf, "Informasi Pelanggan")
	g.startTable(pdf)
	if user != nil {
		g.tableRow(pdf, "Nama", user.FullName)
		g.tableRow(pdf, "Email", user.Email)
		g.tableRow(pdf, "Telepon", user.Phone)
	}
	if policy != nil {
		g.tableRow(pdf, "Nomor Polis", policy.PolicyNumber)
		g.tableRow(pdf, "Status Polis", strings.ToUpper(policy.Status))
	}
	g.endTable(pdf)
	pdf.Ln(6)

	// amount detail
	g.drawSectionTitle(pdf, "Rincian Pembayaran")
	g.startTable(pdf)
	g.tableRow(pdf, "Total Tagihan", rupiah(invoice.Amount))
	if invoice.PaidAmount > 0 {
		g.tableRow(pdf, "Jumlah Dibayar", rupiah(invoice.PaidAmount))
	}
	balance := invoice.Amount - invoice.PaidAmount
	if balance > 0 {
		g.tableRow(pdf, "Sisa Tagihan", rupiah(balance))
	}
	g.endTable(pdf)

	pdf.Ln(6)
	g.drawTotalBox(pdf, "Total", rupiah(invoice.Amount))

	pdf.Ln(15)
	g.drawFooter(pdf)

	var buf strings.Builder
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("pdf output: %w", err)
	}
	return []byte(buf.String()), nil
}

func (g *PDFGenerator) drawHeader(pdf *gofpdf.Fpdf, title string) {
	pdf.SetFillColor(clrDark[0], clrDark[1], clrDark[2])
	pdf.SetTextColor(clrWhite[0], clrWhite[1], clrWhite[2])
	pdf.SetFont("Helvetica", "B", 16)
	pdf.CellFormat(bodyWidth, 14, title, "", 1, "C", true, 0, "")
	pdf.SetTextColor(clrDark[0], clrDark[1], clrDark[2])
	pdf.Ln(4)
}

func (g *PDFGenerator) drawCompanyInfo(pdf *gofpdf.Fpdf) {
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(bodyWidth, 5, "PT Asuransi Digital Indonesia", "", 1, "C", false, 0, "")
	pdf.CellFormat(bodyWidth, 5, "Jl. Sudirman No. 123, Jakarta 12930", "", 1, "C", false, 0, "")
	pdf.CellFormat(bodyWidth, 5, "Tel: (021) 1234-5678 | Email: info@asuransidigital.id", "", 1, "C", false, 0, "")
	pdf.SetTextColor(clrDark[0], clrDark[1], clrDark[2])
}

func (g *PDFGenerator) drawSectionTitle(pdf *gofpdf.Fpdf, title string) {
	pdf.SetFont("Helvetica", "B", 11)
	pdf.SetTextColor(clrPrimary[0], clrPrimary[1], clrPrimary[2])
	pdf.CellFormat(bodyWidth, 7, title, "", 1, "L", false, 0, "")
	pdf.SetDrawColor(clrPrimary[0], clrPrimary[1], clrPrimary[2])
	pdf.Line(margin, pdf.GetY(), pageWidth-margin, pdf.GetY())
	pdf.Ln(4)
}

func (g *PDFGenerator) startTable(pdf *gofpdf.Fpdf) {
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(clrDark[0], clrDark[1], clrDark[2])
}

func (g *PDFGenerator) tableRow(pdf *gofpdf.Fpdf, label, value string) {
	wLabel := bodyWidth * colRatio
	wValue := bodyWidth * (1 - colRatio)

	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetFillColor(clrLightBg[0], clrLightBg[1], clrLightBg[2])
	pdf.CellFormat(wLabel, 8, "  "+label, "1", 0, "L", true, 0, "")

	pdf.SetFont("Helvetica", "", 10)
	pdf.SetFillColor(clrWhite[0], clrWhite[1], clrWhite[2])
	pdf.CellFormat(wValue, 8, "  "+value, "1", 1, "L", true, 0, "")
}

func (g *PDFGenerator) endTable(pdf *gofpdf.Fpdf) {}

func (g *PDFGenerator) drawTotalBox(pdf *gofpdf.Fpdf, label, value string) {
	pdf.SetFont("Helvetica", "B", 14)
	pdf.SetFillColor(clrPrimary[0], clrPrimary[1], clrPrimary[2])
	pdf.SetTextColor(clrWhite[0], clrWhite[1], clrWhite[2])
	pdf.CellFormat(bodyWidth*colRatio, 12, "  "+label, "1", 0, "L", true, 0, "")
	pdf.CellFormat(bodyWidth*(1-colRatio), 12, "  "+value, "1", 1, "L", true, 0, "")
	pdf.SetTextColor(clrDark[0], clrDark[1], clrDark[2])
}

func (g *PDFGenerator) drawFooter(pdf *gofpdf.Fpdf) {
	pdf.SetFont("Helvetica", "", 8)
	pdf.SetTextColor(150, 150, 150)
	pdf.CellFormat(bodyWidth, 4, fmt.Sprintf("Dokumen dibuat otomatis pada %s. Berlaku sebagai bukti digital.",
		time.Now().Format("02 Jan 2006 15:04 WIB")), "", 1, "C", false, 0, "")
}

func rupiah(amount int64) string {
	s := fmt.Sprintf("%d", amount)
	n := len(s)
	if n <= 3 {
		return "Rp " + s
	}
	var parts []string
	for i := n; i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		parts = append([]string{s[start:i]}, parts...)
	}
	return "Rp " + strings.Join(parts, ".")
}
