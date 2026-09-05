package http

// Auth Handler Swagger Documentation

// Register godoc
// @Summary Register user baru
// @Description Mendaftarkan akun user baru dengan email, password, nama lengkap, dan nomor telepon
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "Request body untuk registrasi"
// @Success 201 {object} domain.SuccessResponse "User berhasil didaftarkan"
// @Failure 400 {object} domain.ErrorResponse "Invalid request body atau email sudah terdaftar"
// @Router /auth/register [post]

// Login godoc
// @Summary Login user
// @Description Autentikasi user dengan email dan password, mengembalikan JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Kredensial login"
// @Success 200 {object} domain.LoginResponse "Login berhasil dengan token dan informasi user"
// @Failure 400 {object} domain.ErrorResponse "Invalid request body"
// @Failure 401 {object} domain.ErrorResponse "Email atau password salah"
// @Router /auth/login [post]

// Logout godoc
// @Summary Logout user
// @Description Menghapus cookie autentikasi dan melakukan logout
// @Tags Auth
// @Produce json
// @Success 200 {object} domain.SuccessResponse "Logout berhasil"
// @Router /auth/logout [post]

// Me godoc
// @Summary Get current user
// @Description Mendapatkan informasi user yang sedang login berdasarkan JWT token
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Security CookieAuth
// @Success 200 {object} object{user=domain.User,roles=[]string} "Informasi user"
// @Failure 401 {object} domain.ErrorResponse "Not authenticated"
// @Router /auth/me [get]

// ForgotPassword godoc
// @Summary Request password reset
// @Description Mengirim email reset password ke alamat email terdaftar
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.ForgotPasswordRequest true "Email untuk reset password"
// @Success 200 {object} domain.SuccessResponse "Email reset password telah dikirim"
// @Failure 400 {object} domain.ErrorResponse "Email is required"
// @Router /auth/forgot-password [post]

// ResetPassword godoc
// @Summary Reset password
// @Description Mengatur password baru menggunakan token reset yang diterima via email
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.ResetPasswordRequest true "Token dan password baru"
// @Success 200 {object} domain.SuccessResponse "Password berhasil direset"
// @Failure 400 {object} domain.ErrorResponse "Token invalid atau password terlalu pendek"
// @Router /auth/reset-password [post]

// ListUsers godoc
// @Summary List semua users (Admin)
// @Description Mendapatkan daftar semua user dengan pagination
// @Tags Admin - Users
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Jumlah data per halaman" default(20)
// @Param offset query int false "Offset pagination" default(0)
// @Success 200 {object} domain.PaginatedResponse{data=[]domain.User} "List users"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin role"
// @Router /admin/users [get]

// AssignRole godoc
// @Summary Assign role ke user (Admin)
// @Description Menambahkan role ke user tertentu
// @Tags Admin - Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body domain.AssignRoleRequest true "Role ID yang akan di-assign"
// @Success 200 {object} domain.SuccessResponse "Role berhasil di-assign"
// @Failure 400 {object} domain.ErrorResponse "Invalid request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin role"
// @Router /admin/users/{id}/roles [post]

// RemoveRole godoc
// @Summary Remove role dari user (Admin)
// @Description Menghapus role dari user tertentu
// @Tags Admin - Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body domain.AssignRoleRequest true "Role ID yang akan dihapus"
// @Success 200 {object} domain.SuccessResponse "Role berhasil dihapus"
// @Failure 400 {object} domain.ErrorResponse "Invalid request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin role"
// @Router /admin/users/{id}/roles [delete]

// GetRoles godoc
// @Summary Get semua roles (Admin)
// @Description Mendapatkan daftar semua role yang tersedia di sistem
// @Tags Admin - Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{data=[]domain.Role} "List roles"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin role"
// @Router /admin/roles [get]
