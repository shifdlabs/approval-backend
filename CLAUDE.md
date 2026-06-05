# ShifdLabs Approval Backend — CLAUDE.md

## Tech Stack
- Language: Go 1.24
- Framework: Gin v1.10
- Database: PostgreSQL via GORM v1.25
- Auth: JWT (golang-jwt/jwt/v4)
- Validation: go-playground/validator/v10
- Cache: Redis v9
- Password: bcrypt (golang.org/x/crypto)

## Arsitektur
controller/ → HTTP handler, binding JSON, validasi input
service/    → Business logic
repository/ → Query database (GORM)
model/      → Struct database
data/request/  → DTO input
data/response/ → DTO output
helper/     → JWT, error handler, validator
middleware/ → Auth middleware
router/     → Definisi route

## Aturan Wajib

### 1. Selalu return setelah utils.ErrorResponse
```go
if err != nil {
    utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
    return   // ← WAJIB
}
```

### 2. Selalu panggil helper.ValidateStruct setelah ShouldBindJSON
```go
if err := ctx.ShouldBindJSON(&payload); err != nil {
    utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
    return
}
if errs := helper.ValidateStruct(payload); len(errs) > 0 {
    utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
    return
}
```

### 3. Cek error SEBELUM pakai nilai kembalian
```go
result, err := service.DoSomething(payload)
if err != nil {          // ← cek dulu
    utils.ErrorResponse(ctx, *err)
    return
}
use(result)              // ← baru pakai
```

### 4. Tag validasi standar per field
```
Email:    validate:"required,email"
UUID:     validate:"required,uuid"
Password: validate:"required,min=8,max=200"
Role:     validate:"required,oneof=1 99"
Phone:    validate:"required,min=5,max=20"
```

### 5. Error code yang benar
- Validation error (input salah): 400
- Server/internal error (DB, bcrypt): 500

## Sinkronisasi dengan Frontend

Frontend (Vue 3 + Vuetify) memvalidasi:
- firstName/lastName: required, max 100 chars
- email: required, format email
- phone: required, min 5 max 20
- password: required, min 8, max 200
- role: oneof=1 99
- document.type: oneof=1 2
- document.priority: oneof=1 2 3
- document.state (authorize): oneof=1 2 3

Backend HARUS konsisten dengan aturan di atas.
