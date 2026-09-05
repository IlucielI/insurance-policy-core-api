package http

// Billing Handler Swagger Documentation

// GetInvoices godoc
// @Summary Get invoices user
// @Description Mendapatkan daftar invoice untuk user yang sedang login
// @Tags Billing
// @Produce json
// @Security BearerAuth
// @Security CookieAuth
// @Param limit query int false "Jumlah data per halaman" default(10)
// @Param offset query int false "Offset pagination" default(0)
// @Success 200 {object} domain.PaginatedResponse{data=[]domain.Invoice} "List invoices"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /billing/invoices [get]

// PayInvoice godoc
// @Summary Bayar invoice
// @Description Membuat transaksi pembayaran untuk invoice melalui Midtrans payment gateway
// @Tags Billing
// @Accept json
// @Produce json
// @Param request body domain.CreatePaymentRequest true "Invoice ID yang akan dibayar"
// @Success 200 {object} domain.PaymentResponse "Payment transaction created dengan snap token"
// @Failure 400 {object} domain.ErrorResponse "Invalid request atau invoice tidak ditemukan"
// @Router /billing/pay [post]

// GetPaymentHistory godoc
// @Summary Get riwayat pembayaran
// @Description Mendapatkan riwayat pembayaran user yang sedang login
// @Tags Billing
// @Produce json
// @Security BearerAuth
// @Security CookieAuth
// @Param limit query int false "Jumlah data per halaman" default(10)
// @Param offset query int false "Offset pagination" default(0)
// @Success 200 {object} domain.PaginatedResponse{data=[]domain.Invoice} "Payment history"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /billing/history [get]

// HandlePaymentWebhook godoc
// @Summary Webhook Midtrans (Public)
// @Description Menerima notifikasi dari Midtrans tentang status pembayaran
// @Tags Webhooks
// @Accept json
// @Produce json
// @Param notification body object true "Midtrans webhook notification"
// @Success 200 {object} object{status=string} "Webhook processed"
// @Failure 400 {object} domain.ErrorResponse "Invalid webhook payload"
// @Router /webhooks/payment [post]

// GetAllInvoices godoc
// @Summary Get semua invoices (Admin)
// @Description Mendapatkan daftar semua invoice dengan filter
// @Tags Admin - Billing
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter berdasarkan status (pending, paid, overdue, cancelled)"
// @Param limit query int false "Jumlah data per halaman" default(10)
// @Param offset query int false "Offset pagination" default(0)
// @Success 200 {object} domain.PaginatedResponse{data=[]domain.Invoice} "List invoices"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin or finance role"
// @Router /admin/billing/invoices [get]

// CreateInvoice godoc
// @Summary Buat invoice baru (Admin)
// @Description Membuat invoice baru untuk polis tertentu
// @Tags Admin - Billing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.CreateInvoiceRequest true "Data invoice"
// @Success 201 {object} object{message=string,invoice=domain.Invoice} "Invoice berhasil dibuat"
// @Failure 400 {object} domain.ErrorResponse "Invalid request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin or finance role"
// @Router /admin/billing/invoices [post]

// UpdateInvoiceStatus godoc
// @Summary Update status invoice (Admin)
// @Description Mengubah status invoice secara manual
// @Tags Admin - Billing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Invoice ID"
// @Param request body domain.UpdateInvoiceStatusRequest true "Status baru"
// @Success 200 {object} domain.SuccessResponse "Invoice status updated"
// @Failure 400 {object} domain.ErrorResponse "Invalid request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin or finance role"
// @Router /admin/billing/invoices/{id}/status [put]
