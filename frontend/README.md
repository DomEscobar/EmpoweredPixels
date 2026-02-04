# Frontend (Vue + Pinia + Tailwind)

This is a clean-architecture aligned scaffold. Dependencies are not installed yet.

## Recommended setup

From `frontend/`:

```
npm init -y
npm install vue pinia vue-router
npm install -D vite @vitejs/plugin-vue typescript tailwindcss postcss autoprefixer
```

Then add Vite + Tailwind configs and run:

```
npm run dev
```

## Structure

- `src/app` app shell, router, bootstrap
- `src/shared` API client, styles, UI primitives
- `src/entities` shared entity state (Pinia)
- `src/features` feature use cases
- `src/widgets` composed UI blocks
- `src/pages` route-level pages
