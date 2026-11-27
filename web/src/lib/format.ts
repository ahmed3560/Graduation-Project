export const CURRENCY = 'EGP'; // change this to $ , SAR, ... if you like

export function money(n: number): string {
  return (
    CURRENCY +
    ' ' +
    n.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
  );
}

export const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
