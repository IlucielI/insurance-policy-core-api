package http

// Policy Handler Swagger Documentation

// ListPolicies godoc
// @Summary List polis user
// @Description Mendapatkan daftar polis milik user yang sedang login dengan filter
// @Tags Policies
// @Produce json
// @Security BearerAuth
// @Security CookieAuth
// @Param search query string false "Cari berdasarkan policy number atau product name"
// @Param status query string false "Filter berdasarkan status (active, expired, cancelled)"
// @Param product query string false "Filter berdasarkan product name"
// @Param date_from query string false "Filter tanggal mulai (YYYY-MM-DD)"
// @Param date_to query string false "Filter tanggal akhir (YYYY-MM-DD)"
// @Param limit query int false "Jumlah data per halaman" default(10)
// @Param offset query int false "Offset pagination" default(0)
// @Success 200 {object} domain.PaginatedResponse{data=[]domain.Policy} "List polis"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /policies [get]

// GetPolicy godoc
// @Summary Get detail polis
// @Description Mendapatkan detail polis berdasarkan ID
// @Tags Policies
// @Produce json
// @Param id path string true "Policy ID"
// @Success 200 {object} domain.Policy "Detail polis"
// @Failure 404 {object} domain.ErrorResponse "Policy not found"
// @Router /policies/{id} [get]

// EndorsePolicy godoc
// @Summary Endorse polis
// @Description Mengajukan endorsement (perubahan) pada polis
// @Tags Policies
// @Accept json
// @Produce json
// @Param id path string true "Policy ID"
// @Param request body domain.EndorsePolicyRequest true "Data endorsement"
// @Success 200 {object} domain.SuccessResponse "Endorsement request submitted successfully"
// @Failure 400 {object} domain.ErrorResponse "Invalid request"
// @Router /policies/{id}/endorse [post]

// RenewPolicy godoc
// @Summary Renew polis
// @Description Memperpanjang polis yang akan atau sudah expired
// @Tags Policies
// @Produce json
// @Param id path string true "Policy ID"
// @Success 200 {object} domain.SuccessResponse "Policy renewed successfully"
// @Failure 400 {object} domain.ErrorResponse "Policy cannot be renewed"
// @Router /policies/{id}/renew [post]
