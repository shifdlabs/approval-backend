# ShifdLabs Approval Frontend — CLAUDE.md

## Tech Stack

- **Framework**: Vue 3 + TypeScript
- **UI Library**: Vuetify 3.7.5
- **Router**: `unplugin-vue-router` (file-based routing)
- **HTTP**: `useApi` composable (wraps `ofetch`)
- **Auth**: Cookie-based (`accessToken`), guard di `src/plugins/1.router/index.ts`
- **State**: `ref` / `reactive` (Composition API)
- **Testing**: Playwright (62 tests, semua PASS)
- **Auto-imports**: `unplugin-auto-import` — validator dan composable tersedia global tanpa import manual

---

## Arsitektur Project

```
src/
  pages/              → Route otomatis dari nama file (unplugin-vue-router)
  components/dialogs/ → Dialog form (Create/Update/Change*)
  controllers/        → Logic non-UI (binding data, API call)
  @core/utils/        → validators.ts — semua validator form
  plugins/1.router/   → Auth guard (beforeEach)
  models/             → TypeScript interface
tests/e2e/            → Playwright test files
```

### Penting: Konvensi Routing File-Based

`unplugin-vue-router` mengkonversi **titik (.) dalam nama file menjadi slash (/) di path**:

| File | Route yang dihasilkan |
|------|----------------------|
| `src/pages/admin/document.numbers.vue` | `/admin/document/numbers` |
| `src/pages/admin/app.log.vue` | `/admin/app/log` |

Jangan gunakan titik di nama file kecuali memang ingin membuat nested route.

---

## Aturan Penting saat Mengerjakan Kode Ini

### 1. Validasi Form — Pola Wajib

Semua dialog form menggunakan pola ini:

```vue
<VForm ref="refVForm" @submit.prevent="onFormSubmit">
  <AppTextField
    v-model="fieldName"
    :rules="[requiredValidator, maxLengthValidator(100)]"
  />
  <VBtn type="submit">Submit</VBtn>
</VForm>
```

```ts
const onFormSubmit = async () => {
  refVForm.value?.validate().then(({ valid: isValid }) => {
    if (isValid) callApi()
    else isAllInputtedValid.value = false
  })
}
```

**Jangan pernah memanggil API langsung dari tombol tanpa melewati VForm validate.**

### 2. Dialog Hanya Menutup saat Sukses

```ts
// SALAH — dialog menutup meski API error
const { data, error } = await useApi(...)
emit('update:isDialogVisible', false)  // ← unconditional close

// BENAR
if (data.value?.success) {
  emit('update:isDialogVisible', false)
} else {
  isAllInputtedValid.value = false      // tampilkan error di dalam dialog
}
```

### 3. Validator yang Tersedia (auto-import dari `validators.ts`)

| Validator | Kegunaan |
|-----------|----------|
| `requiredValidator` | Field wajib diisi |
| `emailValidator` | Format email valid |
| `passwordValidator` | Min 8, harus ada uppercase + lowercase + digit + special char |
| `confirmedValidator` | Confirm password harus sama |
| `phoneValidator` | Format nomor telepon (opsional — pakai bersama requiredValidator) |
| `maxLengthValidator(n)` | Panjang string maksimal n karakter |

### 4. Input Type untuk Nomor Telepon

```vue
<!-- SALAH — memblok karakter + dan non-digit -->
<AppTextField type="number" v-model="phone" />

<!-- BENAR — mengizinkan +62 dan format internasional -->
<AppTextField type="tel" v-model="phone" />
```

### 5. Route Guard — Path yang Dilindungi

File: `src/plugins/1.router/index.ts`

Protected routes: `/admin/*`, `/reguler/*`, `/document/*`, `/profile`, `/preview/*`

Jika menambah halaman baru yang butuh auth, pastikan pathnya tercakup di guard.

---

## Bug yang Sudah Diperbaiki (Juni 2026)

### CRITICAL

| ID | Komponen | Fix |
|----|----------|-----|
| FE-BUG-13 | `login-controller.ts` | Hapus hardcoded credentials dari inisialisasi ref |
| FE-BUG-06 | `CreatePublicationFormatDialog.vue` | Bungkus field utama dengan VForm, hubungkan tombol ke onFormSubmit, tambah :rules |

### HIGH

| ID | Komponen | Fix |
|----|----------|-----|
| FE-BUG-14 | `src/plugins/1.router/index.ts` | Tambah `/document`, `/profile`, `/preview` ke route guard |
| FE-BUG-08 | `UpdateUserDialog.vue`, `CreateUserDialog.vue` | emit close hanya di blok success |
| FE-BUG-09 | `ChangeEmailDialog.vue`, `UpdateBiodataDialog.vue` | emit close hanya di blok success |

### MEDIUM

| ID | Komponen | Fix |
|----|----------|-----|
| FE-BUG-01 | Dialog user | Tambah `maxLengthValidator(100)` pada firstName dan lastName |
| FE-BUG-02 | `document/create.vue` | Tambah `maxLengthValidator(200)` pada Subject |
| FE-BUG-03 | Dialog user | Ubah `type="number"` → `type="tel"` pada field phone |
| FE-BUG-05 | `document/create.vue` | External Recipient: rules kondisional berdasarkan tipe dokumen |
| FE-BUG-10 | `validators.ts` | Perluas regex passwordValidator ke semua special chars ASCII |

### LOW

| ID | Komponen | Fix |
|----|----------|-----|
| FE-BUG-12 | `document/create.vue` | Hapus asterisk dari label "Reference *" |

---

## Bug Masih Terbuka

| ID | Komponen | Keterangan |
|----|----------|------------|
| FE-BUG-04 | `CreateUserDialog.vue` — `positionId` | Perlu konfirmasi product: apakah posisi wajib atau opsional? |

---

## QA Test Status (Playwright)

**62/62 PASS** — Jalankan: `pnpm exec playwright test`

| Project | Tests | Status |
|---------|-------|--------|
| Unauthenticated (login, route guard) | 15 | ✓ |
| Admin (`admin@approval.com`) | 21 | ✓ |
| Reguler (`staff1@approval.com`) | 26 | ✓ |

Auth state: `tests/e2e/.auth/` (gitignored). Jalankan ulang jika token expired.

---

## Alignment Frontend ↔ Backend

Temuan QA yang memerlukan perbaikan di backend:

| Endpoint | Temuan | Yang Backend Harus Lakukan |
|----------|--------|---------------------------|
| `PUT /user/biodata` | Menerima `firstName: ""` — tidak ada server-side validation | `validate:"required,min=1,max=100"` di DTO |
| Semua POST/PUT | `helper.ValidateStruct` tidak pernah dipanggil | Aktifkan di semua controller setelah `ShouldBindJSON` |
| Field email | Frontend validasi format email, backend belum tentu | `validate:"required,email"` di semua DTO email |
| Field `role` | Frontend `oneof=1 99`, backend perlu sinkron | `validate:"required,oneof=1 99"` |
| Field enum dokumen | `type`, `priority`, `state` perlu enum validation | `validate:"required,oneof=..."` sesuai nilai valid |
