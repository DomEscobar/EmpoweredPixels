# Security Policy & Best Practices

## Critical Security Requirements

This document outlines the security measures implemented in EmpoweredPixels and best practices for safe deployment.

### 1. JWT Secret Management (CRITICAL)

**Requirement:** The `EP_JWT_SECRET` environment variable is **REQUIRED** and must be at least 32 characters long.

- The application will **fail to start** if `EP_JWT_SECRET` is not set or is too short.
- This is by design — the backend cannot function without a proper secret.

**How to set a proper JWT secret:**

```bash
# Generate a strong 32+ character secret using OpenSSL
openssl rand -base64 32

# Output example (copy this and use it):
# abc123xyz789def456ghi789jkl012mno345pqr678stu

# Set it in your .env file:
echo "EP_JWT_SECRET=abc123xyz789def456ghi789jkl012mno345pqr678stu" >> .env

# Or export it for the current session:
export EP_JWT_SECRET="abc123xyz789def456ghi789jkl012mno345pqr678stu"
```

**Why this matters:**
- JWT tokens are signed with this secret
- A weak secret makes tokens vulnerable to brute-force forgery
- Attackers who know the secret can create valid tokens and impersonate users

### 2. CORS Origin Validation

**Current Implementation:**
- All HTTP requests are validated against a whitelist of allowed origins
- Default allowed origins:
  - `http://localhost:3000` (dev frontend)
  - `http://localhost:5173` (Vite dev server)
  - `http://127.0.0.1:3000` (local testing)
  - `http://127.0.0.1:5173` (local testing)

**Production Deployment:**
Add your production domain(s) to the whitelist in:
```go
// backend/internal/adapter/http/middleware/cors.go
var AllowedOrigins = []string{
	"http://localhost:3000",
	"https://app.example.com",    // Add your prod domain
	"https://www.example.com",    // Add other domains
}
```

**Why this matters:**
- CORS restrictions prevent unauthorized cross-origin requests
- Without validation, attackers could make requests from their own websites
- `Access-Control-Allow-Origin: *` is a common security vulnerability

### 3. WebSocket Origin Validation

**Current Implementation:**
- WebSocket connections validate the `Origin` header before upgrade
- Uses the same whitelist as CORS (see section 2)
- Rejected origins are logged for security monitoring

**Code location:**
```go
// backend/internal/adapter/ws/match_hub.go
CheckOrigin: func(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	return isWSOriginAllowed(origin)
}
```

**Why this matters:**
- WebSocket connections are vulnerable to CSRF attacks if origin is not validated
- Default `CheckOrigin: true` accepts all origins (dangerous!)

### 4. Environment Variable Management

**All secrets must be externalized:**
- ✅ `EP_JWT_SECRET` - JWT signing secret
- ✅ `POSTGRES_PASSWORD` - Database password
- ✅ Other API keys, endpoints, etc.

**Never commit secrets to version control:**
- `.env` is listed in `.gitignore` (commits are rejected automatically)
- Use `.env.example` as a template (safe to commit, has placeholder values)
- Copy `.env.example` to `.env` and fill in real values

```bash
cp .env.example .env
# Edit .env with your real secrets
nano .env
```

### 5. Input Validation & Error Handling

**Implemented Middleware:**
- `ValidateJSON` — Validates request body JSON structure
- `ValidateAuthHeader` — Ensures `Authorization: Bearer <token>` format
- `ValidateQueryParams` — Validates URL query parameters

**Error Responses:**
- Invalid JSON → 400 Bad Request
- Missing auth header format → 400 Bad Request
- Invalid query params → 400 Bad Request
- WebSocket errors are logged with full context (not silent failures)

**Why this matters:**
- Garbage input can crash or exploit services
- Silent failures hide security issues during monitoring

### 6. HTTPS/TLS in Production

**Current Implementation:** HTTP only (suitable for local dev with Docker Compose)

**Production Requirements:**
- All communication must use HTTPS/TLS
- Use a reverse proxy (nginx, Caddy) to terminate TLS
- Obtain certificates from Let's Encrypt (free) or your provider

**Example nginx config:**
```nginx
server {
	listen 443 ssl http2;
	server_name app.example.com;
	
	ssl_certificate /etc/letsencrypt/live/app.example.com/fullchain.pem;
	ssl_certificate_key /etc/letsencrypt/live/app.example.com/privkey.pem;
	
	location / {
		proxy_pass http://backend:54321;
		proxy_set_header Host $host;
		proxy_set_header X-Real-IP $remote_addr;
	}
}
```

### 7. Database Security

**Current Configuration:**
- PostgreSQL user: `ep_user` (not default `postgres`)
- Database: `empoweredpixels`
- Connection: Requires password authentication

**Production Requirements:**
1. Use a strong, randomly generated password (minimum 32 characters)
2. Enable SSL mode in production:
   ```
   EP_DATABASE_URL=postgres://ep_user:password@db.prod.local/empoweredpixels?sslmode=require
   ```
3. Restrict network access — database should only be accessible from backend
4. Enable PostgreSQL audit logging for sensitive operations
5. Regular backups with encryption

### 8. Logging & Monitoring

**Current Implementation:**
- All validation failures are logged
- WebSocket errors include stack context
- Origin rejections are logged for security analysis

**Recommended Additions:**
- Aggregate logs to a centralized system (ELK, Splunk, etc.)
- Alert on repeated validation failures (potential attacks)
- Monitor for failed auth attempts and origin rejections
- Review logs weekly for anomalies

### 9. Dependency Security

**Scan for vulnerabilities:**
```bash
cd backend
go list -json -m all | nancy sleuth
# or
go list -u -m all
```

**Keep dependencies updated:**
```bash
go get -u ./...
go mod tidy
```

### 10. Testing Security Fixes

**Run test suite:**
```bash
cd backend
go test ./internal/adapter/http/middleware/...
```

**Key tests included:**
- `cors_test.go` — Validates CORS origin whitelist enforcement
- `validation_test.go` — Validates input validation middleware
- `config_test.go` — Documents JWT secret requirements

## Incident Response

If you suspect a security issue:

1. **Do not commit the issue** to version control
2. **Rotate secrets immediately** if compromised:
   ```bash
   openssl rand -base64 32  # Generate new JWT secret
   # Update .env and restart services
   docker-compose restart backend
   ```
3. **Review logs** for unauthorized access
4. **Update CORS/WS whitelist** if origins were misconfigured
5. **Disable affected services** if necessary

## Deployment Checklist

Before deploying to production:

- [ ] `EP_JWT_SECRET` is set to a 32+ character random value
- [ ] `POSTGRES_PASSWORD` is set to a 32+ character random value
- [ ] All hardcoded secrets removed from code/config
- [ ] CORS whitelist updated with production domain(s)
- [ ] WebSocket origin whitelist updated (matches CORS)
- [ ] HTTPS/TLS enabled (reverse proxy in front)
- [ ] Database SSL mode enabled (`sslmode=require`)
- [ ] `.env` file is NOT checked into version control
- [ ] Logging & monitoring configured
- [ ] Dependency scan passed (no known vulnerabilities)
- [ ] All tests passing including security tests

## References

- [OWASP Top 10](https://owasp.org/Top10/)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- [CORS Security](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
- [WebSocket Security](https://owasp.org/www-community/attacks/csrf)
- [PostgreSQL Security](https://www.postgresql.org/docs/current/sql-syntax.html)

## Questions?

If you have security concerns, please review this document and the inline code comments in:
- `backend/internal/adapter/http/middleware/cors.go`
- `backend/internal/adapter/http/middleware/validation.go`
- `backend/internal/adapter/ws/match_hub.go`
- `backend/internal/config/config.go`
