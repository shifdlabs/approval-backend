# Spesifikasi Backend — Fitur Lupa Kata Sandi

Dokumen ini mendeskripsikan dua endpoint yang dibutuhkan frontend untuk mendukung alur lupa kata sandi. Frontend sudah selesai dibangun dan menunggu endpoint ini aktif.

---

## Alur Lengkap

```
[Login Page]
  └─ Klik "Lupa kata sandi?"
       │
       ▼
[GET /forgot-password]          ← halaman frontend (tidak ada API)
  └─ User isi email → submit
       │
       ▼
POST /auth/forgot-password      ← ① Endpoint pertama
  └─ Backend kirim email berisi link:
     https://<domain>/reset-password?token=<secure_token>
       │
       ▼
[User klik link di email]
[GET /reset-password?token=...] ← halaman frontend (tidak ada API)
  └─ User isi password baru → submit
       │
       ▼
POST /auth/reset-password       ← ② Endpoint kedua
  └─ Backend validasi token → update password
```

---

## ① POST `/auth/forgot-password`

Menerima email, membuat token reset, dan mengirim email ke pengguna.

### Request

```http
POST /auth/forgot-password
Content-Type: application/json
```

```json
{
  "email": "user@example.com"
}
```

| Field   | Type   | Validasi             |
|---------|--------|----------------------|
| `email` | string | wajib, format email valid |

### Response — Sukses `200 OK`

```json
{
  "success": true,
  "code": 200,
  "message": "Password reset email has been sent."
}
```

### Response — Email tidak ditemukan `404 Not Found`

> **Catatan keamanan:** Direkomendasikan tetap mengembalikan `200` meskipun email tidak terdaftar, untuk mencegah enumerasi akun (_user enumeration_). Namun jika ingin eksplisit, gunakan `404`.

```json
{
  "success": false,
  "code": 404,
  "message": "No account found with this email address."
}
```

### Response — Rate limit `429 Too Many Requests`

```json
{
  "success": false,
  "code": 429,
  "message": "Too many requests. Please try again later."
}
```

### Catatan Implementasi

- Generate token reset yang **kriptografis aman** (misal: `crypto/rand` 32 byte → hex/base64url).
- Simpan token di database (tabel baru atau kolom di tabel users) beserta:
  - `user_id`
  - `token` (hash dari token, bukan plaintext)
  - `expires_at` (waktu kedaluwarsa, disarankan **15–30 menit**)
  - `used_at` (nullable, diisi saat token sudah dipakai)
- Link yang dikirim ke email: `https://<domain>/reset-password?token=<raw_token>`
- Satu user hanya boleh punya satu token aktif — invalidasi token lama jika ada permintaan baru.
- Terapkan **rate limiting** per email/IP (misal: maks 3 request per jam).

---

## ② POST `/auth/reset-password`

Memvalidasi token dan memperbarui kata sandi pengguna.

### Request

```http
POST /auth/reset-password
Content-Type: application/json
```

```json
{
  "token": "a3f9c2e1d...",
  "new_password": "NewPass@123"
}
```

| Field          | Type   | Validasi                                            |
|----------------|--------|-----------------------------------------------------|
| `token`        | string | wajib                                               |
| `new_password` | string | wajib, min 8 karakter, harus mengandung: huruf kapital, huruf kecil, angka, dan karakter spesial (`!@#$%^&*()_+-=[]{}` dll.) |

### Regex validasi password (sama dengan frontend)

```
(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?`~]).{8,}
```

### Response — Sukses `200 OK`

```json
{
  "success": true,
  "code": 200,
  "message": "Password has been successfully reset."
}
```

### Response — Token tidak valid / kedaluwarsa `400 Bad Request`

```json
{
  "success": false,
  "code": 400,
  "message": "Reset token is invalid or has expired."
}
```

### Response — Token sudah pernah dipakai `400 Bad Request`

```json
{
  "success": false,
  "code": 400,
  "message": "Reset token has already been used."
}
```

### Response — Password tidak memenuhi syarat `422 Unprocessable Entity`

```json
{
  "success": false,
  "code": 422,
  "message": "Password does not meet the required complexity."
}
```

### Catatan Implementasi

- Hash token dari request, bandingkan dengan hash yang tersimpan di DB — **jangan simpan token plaintext**.
- Setelah berhasil: set `used_at = now()` dan langsung **invalidasi semua sesi aktif** pengguna (hapus refresh token yang ada).
- Setelah reset sukses, opsional kirim email konfirmasi "kata sandi Anda baru saja diubah".

---

## Format Response Umum

Semua endpoint di sistem ini menggunakan struktur berikut:

```json
{
  "success": boolean,
  "code": number,
  "message": string,
  "data": object | null
}
```

Untuk kedua endpoint di atas, `data` tidak diperlukan (bisa `null` atau dihilangkan).

---

## Template Email Reset Kata Sandi

Isi email yang dikirim ke pengguna harus memuat:

| Elemen     | Isi                                                                 |
|------------|---------------------------------------------------------------------|
| Subject    | `[Shifd] Permintaan Reset Kata Sandi`                               |
| Link reset | `https://<domain>/reset-password?token=<raw_token>`                 |
| Masa aktif | Informasikan bahwa link aktif selama 15–30 menit                    |
| Penutup    | Jika tidak merasa meminta reset, abaikan email ini                  |

---

## Checklist Backend

- [ ] Tabel/kolom `password_reset_tokens` di database
- [ ] `POST /auth/forgot-password` — generate & kirim token
- [ ] `POST /auth/reset-password` — validasi token & update password
- [ ] Token disimpan dalam bentuk **hash**, bukan plaintext
- [ ] Token kedaluwarsa otomatis setelah 15–30 menit
- [ ] Token langsung invalid setelah sekali dipakai
- [ ] Rate limiting per email/IP di endpoint forgot-password
- [ ] Invalidasi sesi aktif setelah reset sukses
- [ ] Email konfirmasi setelah reset (opsional)
- [ ] Cron job untuk membersihkan token kedaluwarsa (opsional)
