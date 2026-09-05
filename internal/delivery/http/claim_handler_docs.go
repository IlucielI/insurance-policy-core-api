package http

// Claim Handler Swagger Documentation

// CreateClaim godoc
// @Summary Buat klaim baru
// @Description Mengajukan klaim baru untuk polis tertentu
// @Tags Claims
// @Accept json
// @Produce json
// @Param request body domain.CreateClaimRequest true "Data klaim"
// @Success 201 {object} object{message=string,claim=domain.Claim} "Klaim berhasil dibuat"
// @Failure 400 {object} domain.ErrorResponse "Invalid request"
// @Router /claims [post]

// GetClaim godoc
// @Summary Get detail klaim
// @Description Mendapatkan detail klaim berdasarkan ID
// @Tags Claims
// @Produce json
// @Param id path string true "Claim ID"
// @Success 200 {object} domain.Claim "Detail klaim"
// @Failure 404 {object} domain.ErrorResponse "Claim not found"
// @Router /claims/{id} [get]

// UploadDocument godoc
// @Summary Upload dokumen klaim
// @Description Upload dokumen pendukung untuk klaim (medical receipt, police report, dll)
// @Tags Claims
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Claim ID"
// @Param file formData file true "File dokumen"
// @Success 200 {object} domain.SuccessResponse "Document uploaded successfully"
// @Failure 400 {object} domain.ErrorResponse "No file uploaded"
// @Router /claims/{id}/documents [put]

// GetClaimTimeline godoc
// @Summary Get timeline klaim
// @Description Mendapatkan riwayat perubahan status dan aktivitas klaim
// @Tags Claims
// @Produce json
// @Param id path string true "Claim ID"
// @Success 200 {object} object{timeline=[]object} "Timeline klaim"
// @Failure 404 {object} domain.ErrorResponse "Claim not found"
// @Router /claims/{id}/timeline [get]

// ListClaims godoc
// @Summary List klaim (Admin)
// @Description Mendapatkan daftar semua klaim dengan filter dan pagination
// @Tags Admin - Claims
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter berdasarkan status (pending, investigating, approved, rejected)"
// @Param policy_id query string false "Filter berdasarkan policy ID"
// @Param user_id query string false "Filter berdasarkan user ID"
// @Param date_from query string false "Filter tanggal mulai (YYYY-MM-DD)"
// @Param date_to query string false "Filter tanggal akhir (YYYY-MM-DD)"
// @Param limit query int false "Jumlah data per halaman" default(20)
// @Param offset query int false "Offset pagination" default(0)
// @Success 200 {object} domain.PaginatedResponse{data=[]domain.Claim} "List klaim"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin or claims_officer role"
// @Router /admin/claims [get]

// UpdateClaimStatus godoc
// @Summary Update status klaim (Admin)
// @Description Mengubah status klaim
// @Tags Admin - Claims
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Claim ID"
// @Param request body domain.UpdateClaimStatusRequest true "Status baru dan notes"
// @Success 200 {object} domain.SuccessResponse "Status berhasil diupdate"
// @Failure 400 {object} domain.ErrorResponse "Invalid request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin or claims_officer role"
// @Router /admin/claims/{id}/status [put]

// ApproveClaim godoc
// @Summary Approve klaim (Admin)
// @Description Menyetujui klaim dan menentukan jumlah approved amount
// @Tags Admin - Claims
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Claim ID"
// @Param request body domain.ApproveClaimRequest true "Approved amount dan notes"
// @Success 200 {object} domain.SuccessResponse "Claim approved successfully"
// @Failure 400 {object} domain.ErrorResponse "Invalid request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin or claims_officer role"
// @Router /admin/claims/{id}/approve [put]

// BulkUpdateClaimStatus godoc
// @Summary Bulk update status klaim (Admin)
// @Description Mengubah status beberapa klaim sekaligus
// @Tags Admin - Claims
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.BulkUpdateClaimStatusRequest true "Claim IDs dan status baru"
// @Success 200 {object} domain.SuccessResponse "Bulk update berhasil"
// @Failure 400 {object} domain.ErrorResponse "Invalid request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin or claims_officer role"
// @Router /admin/claims/bulk-status [post]
