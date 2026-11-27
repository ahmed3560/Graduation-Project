# Wasla – Web (React + TypeScript + Vite)

## Run
    npm install
    npm run dev      # open the link it prints (usually http://localhost:5173)
    npm run build    # production build

## Pages
| Path | File |
|---|---|
| `/` Home | `src/routes/HomePage.tsx` |
| `/signin` | `src/features/auth/SignInPage.tsx` |
| `/signup` | `src/features/auth/SignUpPage.tsx` |
| `/cart` | `src/features/cart/CartPage.tsx` |
| `/checkout` | `src/features/cart/CheckoutPage.tsx` |

Routes are defined in `src/app/router.tsx`. The cart is shared through `src/store/cartStore.ts` + `src/hooks/useCart.ts`.
Everything is demo only (no server yet).
