import type { CartItem } from '../types/cart';

const CART_KEY = 'waslaCart';

// Sample items shown until real products are connected
const DEFAULT_CART: CartItem[] = [
  { id: 'p1', name: 'Express Delivery', meta: 'Same-day · Cairo & Giza', price: 85, qty: 2, icon: 'truck' },
  { id: 'p2', name: 'Packaging Box Set', meta: 'Size M · 10 pieces', price: 120, qty: 1, icon: 'box' },
  { id: 'p3', name: 'Shipment Insurance', meta: 'Covers up to 5,000', price: 45, qty: 1, icon: 'shield' },
];

function read(): CartItem[] {
  try {
    const raw = localStorage.getItem(CART_KEY);
    if (raw) return JSON.parse(raw) as CartItem[];
  } catch {
    /* storage not available */
  }
  return DEFAULT_CART.map((i) => ({ ...i }));
}

let cache: CartItem[] | null = null;
const listeners = new Set<() => void>();

export function getCart(): CartItem[] {
  if (cache === null) cache = read();
  return cache;
}

export function setCart(next: CartItem[]): void {
  cache = next;
  try {
    localStorage.setItem(CART_KEY, JSON.stringify(next));
  } catch {
    /* ignore */
  }
  listeners.forEach((l) => l());
}

export function subscribe(listener: () => void): () => void {
  listeners.add(listener);
  return () => {
    listeners.delete(listener);
  };
}
