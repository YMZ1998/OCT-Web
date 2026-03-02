# Backend quick verification (Windows)

## 1) Prepare `.env`
Create `backend/.env` with required variables:

```env
PORT=8080
MONGO_URI=mongodb://127.0.0.1:27017
MONGO_DB=oct
JWT_SECRET=change_me
```

## 2) Start backend

```powershell
cd backend
go run ./cmd/main.go
```

## 3) Verify APIs end-to-end
Open another PowerShell window and run:

```powershell
cd backend
powershell -ExecutionPolicy Bypass -File .\scripts\verify-api.ps1
```

It validates:
- `POST /api/v1/register`
- `POST /api/v1/login`
- `GET /api/v1/user/:id` with JWT

You can override params:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\verify-api.ps1 -BaseUrl "http://localhost:8080" -Username "demo1" -Password "12345678" -Email "demo1@example.com"
```
