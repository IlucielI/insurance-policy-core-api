package http

// Product Handler Swagger Documentation

// ListProducts godoc
// @Summary List produk asuransi
// @Description Mendapatkan daftar produk asuransi dengan filter dan pagination
// @Tags Products
// @Produce json
// @Param category query string false "Filter berdasarkan kategori (health, life, vehicle, property)"
// @Param is_active query boolean false "Filter produk aktif/nonaktif"
// @Param limit query int false "Jumlah data per halaman" default(10)
// @Param offset query int false "Offset pagination" default(0)
// @Success 200 {object} domain.PaginatedResponse{data=[]domain.Product} "List produk"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /products [get]

// GetProduct godoc
// @Summary Get detail produk
// @Description Mendapatkan detail lengkap produk asuransi berdasarkan ID
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} domain.Product "Detail produk"
// @Failure 404 {object} domain.ErrorResponse "Product not found"
// @Router /products/{id} [get]

// CalculatePremium godoc
// @Summary Hitung premi asuransi
// @Description Menghitung estimasi premi berdasarkan usia, sum assured, dan payment term
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body domain.CalculatePremiumRequest true "Parameter kalkulasi premi"
// @Success 200 {object} domain.CalculatePremiumResponse "Hasil kalkulasi premi"
// @Failure 400 {object} domain.ErrorResponse "Invalid request atau parameter tidak valid"
// @Router /products/{id}/calculate-premium [post]

// SemanticSearch godoc
// @Summary Semantic search produk
// @Description Mencari produk asuransi menggunakan natural language query dengan semantic search
// @Tags Products
// @Accept json
// @Produce json
// @Param request body domain.SemanticSearchRequest true "Query pencarian"
// @Success 200 {object} object{query=string,results=[]domain.Product,count=int} "Hasil pencarian"
// @Failure 400 {object} domain.ErrorResponse "Query is required"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /products/search [post]
