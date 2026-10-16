# Pembacaan Data Query Database di Go (`database/sql`)

Dokumentasi ini menjelaskan mekanisme kerja pengambilan data query SQL di Go, khususnya alasan mengapa kita perlu menggunakan perulangan (`for rows.Next()`) dan bagaimana cara membacanya.

---

## 📌 Mengapa Harus Menggunakan Looping (`for rows.Next()`)?

Go menggunakan konsep **Streaming / Cursor (Penunjuk)** untuk menangani pembacaan data dari database.

1. **`db.DB.Query(...)` mengembalikan Kursor:**
   - Hasil query tidak otomatis dimuat semua ke dalam memori/array.
   - Variabel `rows` adalah sebuah penunjuk (kursor) yang awalnya berada tepat **sebelum baris pertama**.

2. **Peran `rows.Next()`:**
   - Berfungsi memajukan kursor ke baris record berikutnya.
   - Mengembalikan `true` jika masih ada data yang bisa dibaca.
   - Mengembalikan `false` jika data sudah habis atau terjadi error koneksi.

3. **Peran `rows.Scan(...)`:**
   - Hanya membaca dan menyalin **1 baris data** yang sedang ditunjuk kursor ke dalam variabel/struct penampung sementara.

4. **Efisiensi Memori (Memory Efficiency):**
   - Data dibaca secara bertahap baris demi baris, sehingga aplikasi tetap ringan dan tidak menghabiskan RAM meskipun query menghasilkan ribuan record data.

---

## 📖 Cara Membaca Kode (Line by Line)

Potongan kode pada `models/event.go`:

```go
var events []Event

// iterasi untuk query (looping)
for rows.Next() {
    var event Event
    err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
    if err != nil {
        return nil, err
    }
    
    events = append(events, event)
}

return events, nil
```

### Terjemahan ke Bahasa Manusia:
1. **`var events []Event`**
   > *"Siapkan sebuah kotak besar (slice) kosong bernama `events` untuk menampung seluruh daftar event."*
2. **`for rows.Next() {`**
   > *"Selama masih ada baris data berikutnya dari hasil query database, terus jalankan perulangan."*
3. **`var event Event`**
   > *"Siapkan satu wadah kecil kosong baru (`event`) untuk menampung data 1 baris saja."*
4. **`err := rows.Scan(&event.ID, &event.Name, ...)`**
   > *"Salin isi setiap kolom dari baris yang sedang aktif ke dalam alamat memori variabel `event`. Jika tipe data tidak cocok atau gagal, hentikan dan kirim error."*
5. **`events = append(events, event)`**
   > *"Masukkan wadah kecil `event` yang sudah terisi ke dalam kotak besar `events`."*
6. **`return events, nil`**
   > *"Setelah semua baris habis terbaca (looping berhenti), kembalikan kotak besar `events` yang sudah terisi lengkap."*

---

## 🎨 Ilustrasi Visual: Kursor dan Penampung Data

Misal tabel database memiliki **2 baris data**:

```text
TABEL DATABASE (Hasil Query: rows)
-----------------------------------------------------------------
Baris 1 : [ ID: 1 | Name: "Konser Musik" | UserID: 10 ]
Baris 2 : [ ID: 2 | Name: "Workshop Go"  | UserID: 15 ]
-----------------------------------------------------------------
```

### 1. Keadaan Awal (Sebelum Masuk Loop)
Kursor berada di luar/sebelum baris pertama.
```text
           [ KURSOR ] ➔ Berada di luar / sebelum baris 1
-----------------------------------------------------------------
Baris 1 : [ ID: 1 | Name: "Konser Musik" | UserID: 10 ]
Baris 2 : [ ID: 2 | Name: "Workshop Go"  | UserID: 15 ]

Wadah Besar (events) : [ ] (Masih Kosong)
```

### 2. Putaran ke-1 (Loop 1)
- `rows.Next()` memindahkan kursor ke **Baris 1** (`true`).
- `rows.Scan(...)` menyalin kolom Baris 1 ke `event`.
- `append(...)` memasukkannya ke `events`.
```text
Baris 1 : [ ID: 1 | Name: "Konser Musik" | UserID: 10 ] ◄─── [ KURSOR DI SINI ]
Baris 2 : [ ID: 2 | Name: "Workshop Go"  | UserID: 15 ]

Wadah Kecil (event)  : { ID: 1, Name: "Konser Musik", ... }
Wadah Besar (events) : [ {ID: 1, "Konser Musik"} ]
```

### 3. Putaran ke-2 (Loop 2)
- `rows.Next()` memindahkan kursor ke **Baris 2** (`true`).
- `rows.Scan(...)` menyalin kolom Baris 2 ke `event` yang baru.
- `append(...)` memasukkannya ke `events`.
```text
Baris 1 : [ ID: 1 | Name: "Konser Musik" | UserID: 10 ]
Baris 2 : [ ID: 2 | Name: "Workshop Go"  | UserID: 15 ] ◄─── [ KURSOR DI SINI ]

Wadah Kecil (event)  : { ID: 2, Name: "Workshop Go", ... }
Wadah Besar (events) : [ 
    {ID: 1, "Konser Musik"}, 
    {ID: 2, "Workshop Go"} 
]
```

### 4. Putaran ke-3 (Selesai)
- `rows.Next()` dipanggil lagi, namun tidak ada baris ke-3.
- Mengembalikan nilai `false` ➔ Perulangan otomatis **berhenti**.
- Slice `events` siap dikembalikan (`return events, nil`).

---

## 💡 Best Practices Tambahan

1. **Selalu Tutup Kursor dengan `defer rows.Close()`**:
   Mencegah *connection leak* (koneksi database menggantung dan tidak dilepaskan kembali ke *connection pool*).
2. **Cek `rows.Err()` Setelah Loop**:
   Memastikan perulangan berhenti karena data sudah selesai dibaca seluruhnya, bukan terputus mendadak karena error di tengah proses *stream*:
   ```go
   if err = rows.Err(); err != nil {
       return nil, err
   }
   ```
