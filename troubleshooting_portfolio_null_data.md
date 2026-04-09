# Troubleshooting: Data Portofolio Null pada Endpoint GET /portofolio

## Permasalahan
Pengguna melaporkan bahwa setelah berhasil membuat entri portofolio melalui endpoint `POST /portofolio`, ketika mencoba mengambil daftar portofolio menggunakan `GET /portofolio` (untuk pengguna yang terautentikasi), respons yang diterima adalah `{"data":null,"message":"Portofolios fetched successfully"}`. Ini mengindikasikan bahwa data portofolio tidak berhasil diambil atau didekode meskipun proses pembuatan dilaporkan berhasil.

## Analisis Awal
1.  **Pemeriksaan Rute dan Controller**: Rute `POST /portofolio` (memanggil `controllers.CreatePortofolio`) dan `GET /portofolio` (memanggil `controllers.GetMyPortofolios`) telah diperiksa dan logika dasarnya terlihat benar.
2.  **Pemeriksaan Model**: Definisi `models.Portofolio` di `models/portofolio.go` diperiksa untuk tag `json` dan `bson`. Semua field memiliki tag `json:",omitempty"` dan `bson:",omitempty"`.

## Akar Masalah
Akar masalah ditemukan pada penggunaan tag `bson:",omitempty"` pada field `IsDeleted` di struktur `models.Portofolio`.

**Lokasi Kode Bermasalah:**
[portofolio.go](file:///Users/user/Documents/Go/golang-api-portofolio/models/portofolio.go)

```go
type Portofolio struct {
	// ... field lainnya
	IsDeleted bool               `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
	// ... field lainnya
}
```

**Penjelasan:**
*   Ketika sebuah portofolio baru dibuat, nilai `IsDeleted` secara default adalah `false`.
*   Tag `bson:",omitempty"` pada field `IsDeleted` menyebabkan field ini **tidak disimpan** ke database MongoDB jika nilainya adalah `false` (atau nilai nol lainnya untuk tipe data yang relevan).
*   Akibatnya, dokumen portofolio yang tersimpan di MongoDB tidak memiliki field `is_deleted`.
*   Namun, fungsi `controllers.GetMyPortofolios` menggunakan filter `{"is_deleted": false}`. Karena dokumen di database tidak memiliki field `is_deleted` sama sekali, filter ini gagal menemukan kecocokan, sehingga mengembalikan daftar portofolio kosong (yang kemudian di-marshal menjadi `null` di respons JSON karena struktur respons).

## Solusi
Solusinya adalah memastikan bahwa field `is_deleted` selalu ada di dokumen MongoDB, terlepas dari nilainya (`true` atau `false`). Ini dicapai dengan menghapus tag `omitempty` dari field `IsDeleted` di `models.Portofolio`.

**Perubahan Kode:**
[portofolio.go](file:///Users/user/Documents/Go/golang-api-portofolio/models/portofolio.go)

```go
type Portofolio struct {
	// ... field lainnya
	IsDeleted bool               `json:"is_deleted" bson:"is_deleted"` // Tag omitempty dihapus
	// ... field lainnya
}
```

## Langkah-langkah Verifikasi
1.  **Modifikasi Model**: Mengubah definisi `IsDeleted` di `models/portofolio.go` seperti di atas.
2.  **Hapus Data Lama**: Menghapus semua dokumen portofolio yang ada di database MongoDB (karena dokumen lama tidak memiliki field `is_deleted`). Contoh perintah:
    ```bash
    docker exec -it mongodb mongosh portofolio -u admin -p secret123 --authenticationDatabase admin --eval "db.portofolio.deleteMany({})"
    ```
3.  **Buat Portofolio Baru**: Membuat entri portofolio baru melalui `POST /portofolio`.
    ```bash
    curl -L -X POST http://localhost:8081/portofolio -H "Content-Type: application/json" -H "Authorization: Bearer <YOUR_JWT_TOKEN>" -d '{"name": "My Repaired Portfolio"}'
    ```
4.  **Ambil Portofolio**: Mengambil daftar portofolio melalui `GET /portofolio`.
    ```bash
    curl -L -X GET http://localhost:8081/portofolio -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
    ```

**Hasil Verifikasi:**
Setelah menerapkan perubahan dan melakukan langkah-langkah verifikasi, endpoint `GET /portofolio` berhasil mengembalikan data portofolio yang telah dibuat, mengkonfirmasi bahwa masalah telah teratasi.

---
**Catatan Edukasi**:
Tag `omitempty` sangat berguna untuk mengurangi ukuran dokumen di database dan respons JSON dengan menghilangkan field yang memiliki nilai default atau nol. Namun, penting untuk berhati-hati saat menggunakannya pada field yang digunakan dalam kriteria pencarian (`filter`) di database, terutama jika nilai nol dari field tersebut adalah bagian penting dari logika filter. Dalam kasus seperti ini, lebih baik memastikan field selalu ada di database.
