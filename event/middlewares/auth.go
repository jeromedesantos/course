package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeromedesantos/rest-api/utils"
)

func Authenticate(context *gin.Context) {

	var (
		token	string
		userId	int64
		err		error
	)

	// cek token
	token = context.Request.Header.Get("Authorization")
	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Not authorized.",
		})
		return
	}

	// verify token
	userId, err = utils.VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Not authorized.",
		})
		return
	}

	context.Set("userId", userId) // --> simpan userId ke context agar bisa diakses di handler/routes berikutnya
	context.Next() // --> opsional jika ditaruh di akhir, karena Gin otomatis melanjutkan ke handler berikutnya.
	// Catatan: context.Next() berfungsi mengeksekusi handler berikutnya terlebih dahulu, lalu kembali ke baris setelahnya (berguna untuk post-processing seperti hitung latency/logging).
}
/*
Context.AbortWithStatusJSON vs Context.JSON

# AbortWithStatusJSON :
Mengirim response JSON ke client DAN memanggil Abort(), sehingga handler/controller berikutnya di antrean tidak akan dijalankan.

# JSON :
Hanya mengirim response JSON ke client, TIDAK membatalkan antrean. Gin tetap akan mengeksekusi handler berikutnya meskipun ada return.

-------------------------------

PERBEDAAN ARSITEKTUR & CARA KERJA context.Next():

# Express (Manual):
Middleware lanjut HANYA jika kamu panggil next(). Kalau mau stop, cukup jangan panggil next().

# Gin (Otomatis):
Middleware dan handler disimpan di dalam antrean array yang otomatis di-loop sampai habis.
- c.Abort(): Menghentikan antrean agar handler di belakangnya batal dieksekusi.
- c.Next(): Menjeda middleware saat ini -> menjalankan handler berikutnya sampai selesai -> lalu kembali ke middleware ini untuk menjalankan baris di bawah c.Next().
*/