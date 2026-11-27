import { useSyncExternalStore } from 'react';
import { getCart, setCart, subscribe } from '../store/cartStore';

// One cart shared by every page (navbar badge, cart page, checkout)
export function useCart() {
  const cart = useSyncExternalStore(subscribe, getCart);

  const count = cart.reduce((n, i) => n + i.qty, 0);
  const subtotal = cart.reduce((s, i) => s + i.price * i.qty, 0);

  return {
    cart,
    count,
    subtotal,
    increase: (id: string) =>
      setCart(getCart().map((i) => (i.id === id && i.qty < 99 ? { ...i, qty: i.qty + 1 } : i))),
    decrease: (id: string) =>
      setCart(getCart().map((i) => (i.id === id && i.qty > 1 ? { ...i, qty: i.qty - 1 } : i))),
    remove: (id: string) => setCart(getCart().filter((i) => i.id !== id)),
    clear: () => setCart([]),
  };
}
