# API Spec — Fitur Lupa Kata Sandi

Dokumen ini ditujukan untuk tim frontend. Berisi spesifikasi dua endpoint yang mendukung alur lupa kata sandi.

---

## Alur

```
[Halaman Login]
  └─ Klik "Lupa kata sandi?"
       │
       ▼
[GET /forgot-password]          ← halaman frontend (tidak ada API)
  └─ User isi email → submit
       │
       ▼
POST /api/auth/forgot-password  ← ① Endpoint pertama
  └─ Backend kirim email berisi link:
     https://<domain>/reset-password?token=<token>
       │
       ▼
[User klik link di email]
[GET /reset-password?token=...] ← halaman frontend, baca query param "token"
  └─ User isi password baru → submit
       │
       ▼
POST /api/auth/reset-password   ← ② Endpoint kedua
  └─ Backend validasi token → update password
```

---

## ① POST `/api/auth/forgot-password`

### Request

```http
POST /api/auth/forgot-password
Content-Type: application/json
```

```json
{
  "email": "user@example.com"
}
```

| Field   | Type   | Validasi              |
|---------|--------|-----------------------|
| `email` | string | wajib, format email   |

### Response

**Sukses `200 OK`**
```json
{
  "success": true,
  "code": 200,
  "message": "Password reset email has been sent.",
  "data": null
}
```

**Validasi gagal `400 Bad Request`**
```json
{
  "success": false,
  "code": 400,
  "message": "Bad Request",
  "data": null
}
```

**Rate limit `429 Too Many Requests`**
```json
{
  "success": false,
  "code": 429,
  "message": "Too many requests. Please try again later.",
  "data": null
}
```

**Gagal kirim email `500`**
```json
{
  "success": false,
  "code": 500,
  "message": "Failed to send reset email. Please try again.",
  "data": null
}
```

> **Catatan:** Response `200` selalu dikembalikan meskipun email tidak terdaftar di sistem. Frontend tidak perlu membedakan kasus ini — cukup tampilkan pesan generik seperti _"Jika email terdaftar, link reset akan dikirim."_

---

## ② POST `/api/auth/reset-password`

### Request

```http
POST /api/auth/reset-password
Content-Type: application/json
```

```json
{
  "token": "<token dari query param URL>",
  "new_password": "NewPass@123"
}
```

| Field          | Type   | Validasi                                      |
|----------------|--------|-----------------------------------------------|
| `token`        | string | wajib                                         |
| `new_password` | string | wajib, min 8 karakter, wajib mengandung: huruf kapital, huruf kecil, angka, karakter spesial |

### Validasi password (regex — sama dengan frontend)

```
(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?`~]).{8,}
```

### Response

**Sukses `200 OK`**
```json
{
  "success": true,
  "code": 200,
  "message": "Password has been successfully reset.",
  "data": null
}
```

**Validasi input gagal `400 Bad Request`**
```json
{
  "success": false,
  "code": 400,
  "message": "Bad Request",
  "data": null
}
```

**Token tidak valid atau kedaluwarsa `400 Bad Request`**
```json
{
  "success": false,
  "code": 400,
  "message": "Reset token is invalid or has expired.",
  "data": null
}
```

**Token sudah pernah dipakai `400 Bad Request`**
```json
{
  "success": false,
  "code": 400,
  "message": "Reset token has already been used.",
  "data": null
}
```

**Password tidak memenuhi syarat `422 Unprocessable Entity`**
```json
{
  "success": false,
  "code": 422,
  "message": "Password does not meet the required complexity.",
  "data": null
}
```

---

## Catatan Implementasi Frontend

### Cara ambil token dari URL

```js
const urlParams = new URLSearchParams(window.location.search)
const token = urlParams.get('token')
```

### Penanganan error yang disarankan

| Code | Pesan untuk user |
|------|-----------------|
| `400` (token) | "Link reset tidak valid atau sudah kedaluwarsa. Silakan minta link baru." |
| `400` (token used) | "Link reset sudah pernah digunakan. Silakan minta link baru." |
| `422` | "Password tidak memenuhi syarat. Pastikan mengandung huruf kapital, huruf kecil, angka, dan karakter spesial." |
| `429` | "Terlalu banyak permintaan. Coba lagi dalam 1 jam." |
| `500` | "Terjadi kesalahan. Silakan coba lagi." |

### Token berlaku selama 30 menit

Jika user membuka link lebih dari 30 menit setelah email diterima, token sudah kedaluwarsa dan harus meminta link baru.
