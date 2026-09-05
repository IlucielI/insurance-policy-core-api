package http

// Audit, Chat, Email Handler Swagger Documentation

// GetAuditLog godoc
// @Summary Get detail audit log (Admin)
// @Description Mendapatkan detail audit log berdasarkan ID
// @Tags Admin - Audit Logs
// @Produce json
// @Security BearerAuth
// @Param id path string true "Audit Log ID"
// @Success 200 {object} domain.AuditLog "Detail audit log"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin role"
// @Router /admin/audit-logs/{id} [get]

// ListAuditLogs godoc
// @Summary List audit logs (Admin)
// @Description Mendapatkan daftar audit logs dengan filter
// @Tags Admin - Audit Logs
// @Produce json
// @Security BearerAuth
// @Param user_id query string false "Filter berdasarkan user ID"
// @Param action query string false "Filter berdasarkan action"
// @Param entity_type query string false "Filter berdasarkan entity type"
// @Param entity_id query string false "Filter berdasarkan entity ID"
// @Param date_from query string false "Tanggal mulai (YYYY-MM-DD)"
// @Param date_to query string false "Tanggal akhir (YYYY-MM-DD)"
// @Param limit query int false "Jumlah data per halaman" default(50)
// @Param offset query int false "Offset pagination" default(0)
// @Success 200 {object} domain.PaginatedResponse{data=[]domain.AuditLog} "List audit logs"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin role"
// @Router /admin/audit-logs [get]

// GetAvailableActions godoc
// @Summary Get available audit actions (Admin)
// @Description Mendapatkan daftar action yang tersedia di audit logs
// @Tags Admin - Audit Logs
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{actions=[]string} "List actions"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden"
// @Router /admin/audit-logs/actions [get]

// GetEntityTypes godoc
// @Summary Get available entity types (Admin)
// @Description Mendapatkan daftar entity type yang tersedia di audit logs
// @Tags Admin - Audit Logs
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{entity_types=[]string} "List entity types"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden"
// @Router /admin/audit-logs/entity-types [get]

// GetEntityAuditLogs godoc
// @Summary Get audit logs untuk entity tertentu (Admin)
// @Description Mendapatkan audit trail untuk entity tertentu
// @Tags Admin - Audit Logs
// @Produce json
// @Security BearerAuth
// @Param type path string true "Entity type (application, claim, policy, user)"
// @Param id path string true "Entity ID"
// @Success 200 {object} object{data=[]domain.AuditLog} "Entity audit logs"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden"
// @Router /admin/audit-logs/entity/{type}/{id} [get]

// SendMessage godoc
// @Summary Kirim pesan chat
// @Description Mengirim pesan ke chatbot AI untuk mendapatkan bantuan
// @Tags Chat
// @Accept json
// @Produce json
// @Param request body domain.SendMessageRequest true "User ID dan pesan"
// @Success 200 {object} object{reply=string,message_id=string} "Balasan dari AI"
// @Failure 400 {object} domain.ErrorResponse "Invalid request"
// @Router /chat [post]

// GetHistory godoc
// @Summary Get riwayat chat
// @Description Mendapatkan riwayat percakapan chat untuk user
// @Tags Chat
// @Produce json
// @Param user_id query string true "User ID"
// @Param limit query int false "Jumlah pesan" default(50)
// @Success 200 {object} object{messages=[]object} "Chat history"
// @Failure 400 {object} domain.ErrorResponse "user_id required"
// @Router /chat/history [get]

// SendApplicationApprovedEmail godoc
// @Summary Kirim email approval aplikasi (Admin)
// @Description Mengirim email notifikasi approval aplikasi ke user
// @Tags Admin - Email
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Application ID"
// @Param request body domain.SendApplicationApprovedEmailRequest true "Additional message (optional)"
// @Success 200 {object} domain.SuccessResponse "Email sent successfully"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin or underwriter role"
// @Router /admin/email/application-approved/{id} [post]

// SendClaimStatusEmail godoc
// @Summary Kirim email status klaim (Admin)
// @Description Mengirim email notifikasi perubahan status klaim
// @Tags Admin - Email
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Claim ID"
// @Param request body domain.SendClaimStatusEmailRequest true "Custom message (optional)"
// @Success 200 {object} domain.SuccessResponse "Email sent successfully"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin or claims_officer role"
// @Router /admin/email/claim-status/{id} [post]

// BulkSendEmail godoc
// @Summary Kirim email massal (Admin)
// @Description Mengirim email ke multiple users sekaligus
// @Tags Admin - Email
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.BulkSendEmailRequest true "User IDs, subject, dan body email"
// @Success 200 {object} domain.SuccessResponse "Bulk email sent"
// @Failure 400 {object} domain.ErrorResponse "Invalid request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin role"
// @Router /admin/email/bulk-send [post]

// PreviewEmail godoc
// @Summary Preview template email (Admin)
// @Description Melihat preview template email sebelum dikirim
// @Tags Admin - Email
// @Produce json
// @Security BearerAuth
// @Param type query string true "Tipe email (application_approved, claim_status, payment_reminder)"
// @Success 200 {object} object{preview=string} "HTML preview"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden"
// @Router /admin/email/preview/{type} [get]

// ReviewApplication godoc
// @Summary AI review aplikasi (Admin)
// @Description Melakukan AI-powered review terhadap aplikasi menggunakan LLM
// @Tags Admin - AI
// @Produce json
// @Security BearerAuth
// @Param id path string true "Application ID"
// @Success 200 {object} object{review=string,risk_level=string,recommendation=string} "AI review result"
// @Failure 400 {object} domain.ErrorResponse "Invalid application ID"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin or underwriter role"
// @Router /admin/ai-review/{id} [post]

// GetPreferences godoc
// @Summary Get notification preferences
// @Description Mendapatkan preferensi notifikasi user yang sedang login
// @Tags Notification Preferences
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.NotificationPreferencesRequest "Notification preferences"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /notification-preferences [get]

// UpdatePreferences godoc
// @Summary Update notification preferences
// @Description Mengupdate preferensi notifikasi user
// @Tags Notification Preferences
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.NotificationPreferencesRequest true "Preferensi notifikasi"
// @Success 200 {object} domain.SuccessResponse "Preferences updated successfully"
// @Failure 400 {object} domain.ErrorResponse "Invalid request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /notification-preferences [put]
