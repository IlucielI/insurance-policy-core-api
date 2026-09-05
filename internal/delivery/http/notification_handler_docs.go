package http

// Notification Handler Swagger Documentation

// WebSocketHandler godoc
// @Summary WebSocket untuk real-time notifications
// @Description Membuka koneksi WebSocket untuk menerima notifikasi real-time
// @Tags Notifications
// @Param user_id query string true "User ID untuk subscribe notifications"
// @Success 101 {string} string "Switching Protocols - WebSocket connection established"
// @Failure 401 {object} domain.ErrorResponse "user_id required"
// @Router /notifications/ws [get]

// GetNotifications godoc
// @Summary Get notifikasi user
// @Description Mendapatkan daftar notifikasi untuk user tertentu
// @Tags Notifications
// @Produce json
// @Param user_id query string true "User ID"
// @Param limit query int false "Jumlah data per halaman" default(20)
// @Param offset query int false "Offset pagination" default(0)
// @Success 200 {object} object{notifications=[]domain.Notification} "List notifikasi"
// @Failure 400 {object} domain.ErrorResponse "user_id required"
// @Router /notifications [get]

// GetUnreadCount godoc
// @Summary Get jumlah notifikasi belum dibaca
// @Description Mendapatkan jumlah notifikasi yang belum dibaca untuk user
// @Tags Notifications
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {object} object{unread_count=int} "Jumlah notifikasi belum dibaca"
// @Failure 400 {object} domain.ErrorResponse "user_id required"
// @Router /notifications/unread-count [get]

// MarkAsRead godoc
// @Summary Tandai notifikasi sebagai dibaca
// @Description Menandai satu notifikasi sebagai sudah dibaca
// @Tags Notifications
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} domain.SuccessResponse "Notification marked as read"
// @Failure 400 {object} domain.ErrorResponse "notification_id required"
// @Router /notifications/{id}/read [put]

// MarkAllAsRead godoc
// @Summary Tandai semua notifikasi sebagai dibaca
// @Description Menandai semua notifikasi user sebagai sudah dibaca
// @Tags Notifications
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {object} domain.SuccessResponse "All notifications marked as read"
// @Failure 400 {object} domain.ErrorResponse "user_id required"
// @Router /notifications/read-all [put]

// DeleteNotification godoc
// @Summary Hapus notifikasi
// @Description Menghapus satu notifikasi
// @Tags Notifications
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} domain.SuccessResponse "Notification deleted"
// @Failure 400 {object} domain.ErrorResponse "notification_id required"
// @Router /notifications/{id} [delete]
