# API Performance Testing with k6

Project ini berisi skrip pengujian performa untuk API Portfolio menggunakan [k6](https://k6.io/).

## Prasyarat
- Instal k6 di mesin lokal atau VPS Anda:
  - macOS: `brew install k6`
  - Linux (Debian/Ubuntu): `sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69` dan ikuti panduan resmi.
  - Windows: `choco install k6`

## Cara Menjalankan Tes

### 1. Load Test (Default)
Menguji performa API dengan beban bertahap (ramp-up) hingga 20 user bersamaan.
```bash
k6 run tests/load_test.js
```

### 2. Smoke Test
Pengujian cepat untuk memastikan skrip berjalan tanpa error dengan beban minimal (1 user).
Edit `options` di `load_test.js` dan jalankan:
```bash
k6 run tests/load_test.js
```

## Skenario Pengujian
Skrip ini menguji alur berikut secara otomatis:
1. **Health Check**: Memastikan API aktif.
2. **Public List**: Mengambil data portofolio publik.
3. **Auth Flow**: Melakukan pendaftaran user baru secara acak dan login untuk mendapatkan token.
4. **Protected Area**: Menggunakan token JWT untuk mengakses data portofolio pribadi.

## Ambang Batas (Thresholds)
Tes akan dianggap gagal jika:
- Tingkat kegagalan request > 1% (`http_req_failed`).
- 95% request membutuhkan waktu lebih dari 500ms (`http_req_duration`).
