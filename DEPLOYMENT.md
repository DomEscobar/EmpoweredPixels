# Deployment Guide

## Automated Deployment

**Git Push triggers automatic deployment to VPS.**

### How it works:
1. Push to `feature/shop-mvp-implementation` (or main)
2. GitHub Actions / VPS pulls latest code
3. Backend rebuilds automatically
4. Frontend (dist/) is already built and committed
5. Database migrations run on startup

### Manual Deployment (if needed):
```bash
ssh root@152.53.118.78
cd /opt/empoweredpixels
git pull origin feature/shop-mvp-implementation
cd backend && go build -o api ./cmd/api && systemctl restart empoweredpixels-api
```

### Live URLs:
- Frontend: http://152.53.118.78:49100
- Backend: http://152.53.118.78:49101

## Shop MVP Status
✅ Backend API complete
✅ Frontend complete  
✅ Database migrations
✅ Build verified
🔄 Awaiting VPS deployment

## Next: Attunement System
See `/docs/GAME_DESIGN.md` for specifications.
