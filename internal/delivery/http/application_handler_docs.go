package http

// Application Handler Swagger Documentation

// CreateApplication godoc
// @Summary Buat aplikasi asuransi baru
// @Description Membuat aplikasi asuransi baru untuk produk tertentu
// @Tags Applications
// @Accept json
// @Produce json
// @Param request body domain.CreateApplicationRequest true "Data aplikasi"
// @Success 201 {object} object{message=string,application=domain.Application} "Aplikasi berhasil dibuat"
// @Failure 400 {object} domain.ErrorResponse "Invalid request body"
// @Router /applications [post]

// GetApplication godoc
// @Summary Get detail aplikasi
// @Description Mendapatkan detail aplikasi berdasarkan ID
// @Tags Applications
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} domain.Application "Detail aplikasi"
// @Failure 404 {object} domain.ErrorResponse "Application not found"
// @Router /applications/{id} [get]

// ListApplications godoc
// @Summary List aplikasi (Admin)
// @Description Mendapatkan daftar semua aplikasi dengan filter dan pagination
// @Tags Admin - Applications
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter berdasarkan status (pending, approved, rejected)"
// @Param product_id query string false "Filter berdasarkan product ID"
// @Param user_id query string false "Filter berdasarkan user ID"
// @Param limit query int false "Jumlah data per halaman" default(20)
// @Param offset query int false "Offset pagination" default(0)
// @Success 200 {object} domain.PaginatedResponse{data=[]domain.Application} "List aplikasi"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin or underwriter role"
// @Router /admin/applications [get]

// UpdateStatus godoc
// @Summary Update status aplikasi (Admin)
// @Description Mengubah status aplikasi (approve/reject)
// @Tags Admin - Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Application ID"
// @Param request body domain.UpdateApplicationStatusRequest true "Status baru dan notes"
// @Success 200 {object} domain.SuccessResponse "Status berhasil diupdate"
// @Failure 400 {object} domain.ErrorResponse "Invalid request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin or underwriter role"
// @Router /admin/applications/{id}/status [put]

// BulkUpdateStatus godoc
// @Summary Bulk update status aplikasi (Admin)
// @Description Mengubah status beberapa aplikasi sekaligus
// @Tags Admin - Applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.BulkUpdateStatusRequest true "Application IDs dan status baru"
// @Success 200 {object} domain.SuccessResponse "Bulk update berhasil"
// @Failure 400 {object} domain.ErrorResponse "Invalid request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin or underwriter role"
// @Router /admin/applications/bulk-status [post]
