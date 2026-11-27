import { createBrowserRouter } from 'react-router-dom';
import MainLayout from './layouts/MainLayout';
import HomePage from '../routes/HomePage';
import SignInPage from '../features/auth/SignInPage';
import SignUpPage from '../features/auth/SignUpPage';
import CartPage from '../features/cart/CartPage';
import CheckoutPage from '../features/cart/CheckoutPage';

export const router = createBrowserRouter([
  // Pages that show the navbar
  {
    element: <MainLayout />,
    children: [
      { path: '/', element: <HomePage /> },
      { path: '/cart', element: <CartPage /> },
      { path: '/checkout', element: <CheckoutPage /> },
    ],
  },
  // Pages without the navbar
  { path: '/signin', element: <SignInPage /> },
  { path: '/signup', element: <SignUpPage /> },
]);
