package http

// Document Handler Swagger Documentation

// ListDocuments godoc
// @Summary List dokumen
// @Description Mendapatkan daftar dokumen yang diupload
// @Tags Documents
// @Produce json
// @Param entity_type query string false "Filter berdasarkan entity type (claim, policy, application)"
// @Param entity_id query string false "Filter berdasarkan entity ID"
// @Param limit query int false "Jumlah data per halaman" default(20)
// @Param offset query int false "Offset pagination" default(0)
// @Success 200 {object} domain.PaginatedResponse{data=[]domain.Document} "List dokumen"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /documents [get]

// UploadDocument godoc
// @Summary Upload dokumen
// @Description Upload file dokumen ke MinIO storage
// @Tags Documents
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File yang akan diupload"
// @Param entity_type formData string true "Tipe entitas (claim, policy, application)"
// @Param entity_id formData string true "ID entitas"
// @Success 200 {object} domain.UploadDocumentResponse "Document uploaded successfully"
// @Failure 400 {object} domain.ErrorResponse "No file uploaded atau parameter tidak lengkap"
// @Router /documents/upload [post]

// DownloadDocument godoc
// @Summary Download dokumen
// @Description Download file dokumen dari MinIO storage
// @Tags Documents
// @Produce application/octet-stream
// @Param id path string true "Document ID"
// @Success 200 {file} file "File dokumen"
// @Failure 404 {object} domain.ErrorResponse "Document not found"
// @Router /documents/{id}/download [get]

// DeleteDocument godoc
// @Summary Hapus dokumen
// @Description Menghapus dokumen dari database dan MinIO storage
// @Tags Documents
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} domain.SuccessResponse "Document deleted successfully"
// @Failure 404 {object} domain.ErrorResponse "Document not found"
// @Router /documents/{id} [delete]
