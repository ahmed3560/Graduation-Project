import { ChangeEvent, FormEvent, useState } from 'react';
import { Link } from 'react-router-dom';
import { useCart } from '../../hooks/useCart';
import { EMAIL_RE, money } from '../../lib/format';
import { CardIcon, CashIcon, CheckIcon } from '../../components/common/icons';

type Pay = 'cash' | 'card';

type Form = {
  coName: string;
  coPhone: string;
  coEmail: string;
  coCity: string;
  coArea: string;
  coAddress: string;
  coNotes: string;
  cardNumber: string;
  cardName: string;
  cardExpiry: string;
  cardCvv: string;
};

type FieldId = keyof Form;

type Confirmation = {
  orderNo: string;
  first: string;
  address: string;
  payText: string;
  total: number;
};

const EMPTY: Form = {
  coName: '', coPhone: '', coEmail: '', coCity: '', coArea: '', coAddress: '', coNotes: '',
  cardNumber: '', cardName: '', cardExpiry: '', cardCvv: '',
};

function expiryOk(v: string): boolean {
  const m = /^(\d{2})\/(\d{2})$/.exec(v);
  if (!m) return false;
  const month = +m[1];
  const year = 2000 + +m[2];
  if (month < 1 || month > 12) return false;
  return new Date(year, month, 1) > new Date(); // valid until the end of its month
}

export default function CheckoutPage() {
  const { cart, subtotal, clear } = useCart();
  const [form, setForm] = useState<Form>(EMPTY);
  const [errors, setErrors] = useState<Partial<Record<FieldId, string>>>({});
  const [pay, setPay] = useState<Pay>('cash');
  const [placing, setPlacing] = useState(false);
  const [done, setDone] = useState<Confirmation | null>(null);

  // Update a field and clear its error as soon as the customer starts fixing it
  function set(id: FieldId, value: string) {
    setForm((f) => ({ ...f, [id]: value }));
    setErrors((e) => ({ ...e, [id]: '' }));
  }
  const onText = (id: FieldId) => (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) =>
    set(id, e.target.value);

  // Auto-format the card inputs
  const onCardNumber = (e: ChangeEvent<HTMLInputElement>) => {
    const d = e.target.value.replace(/\D/g, '').slice(0, 19);
    set('cardNumber', d.replace(/(.{4})/g, '$1 ').trim());
  };
  const onExpiry = (e: ChangeEvent<HTMLInputElement>) => {
    const d = e.target.value.replace(/\D/g, '').slice(0, 4);
    set('cardExpiry', d.length > 2 ? `${d.slice(0, 2)}/${d.slice(2)}` : d);
  };
  const onCvv = (e: ChangeEvent<HTMLInputElement>) =>
    set('cardCvv', e.target.value.replace(/\D/g, '').slice(0, 4));

  function validate(): FieldId[] {
    const next: Partial<Record<FieldId, string>> = {};
    const bad: FieldId[] = []; // invalid fields, in page order
    function check(id: FieldId, ok: boolean, msg: string) {
      next[id] = ok ? '' : msg;
      if (!ok) bad.push(id);
    }
    const v = (id: FieldId) => form[id].trim();

    check('coName', v('coName').length >= 3, 'Please enter your full name.');
    check('coPhone', /^\+?[0-9\s-]{8,16}$/.test(v('coPhone')), 'Enter a valid phone number.');
    check('coEmail', EMAIL_RE.test(v('coEmail')), 'Enter a valid email address.');
    check('coCity', v('coCity') !== '', 'Please enter your city.');
    check('coArea', v('coArea') !== '', 'Please enter your area.');
    check('coAddress', v('coAddress').length >= 6, 'Please enter your full street address.');

    if (pay === 'card') {
      const digits = form.cardNumber.replace(/\D/g, '');
      check('cardNumber', digits.length >= 13 && digits.length <= 19, 'Enter a valid card number.');
      check('cardName', v('cardName').length >= 3, 'Enter the name on the card.');
      check('cardExpiry', expiryOk(v('cardExpiry')), 'Enter a valid expiry date (MM/YY).');
      check('cardCvv', /^\d{3,4}$/.test(form.cardCvv), 'Enter the 3 or 4 digit CVV.');
    }
    setErrors(next);
    return bad;
  }

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    const bad = validate();
    if (bad.length) {
      document.getElementById(bad[0])?.focus();
      return;
    }

    setPlacing(true);
    const total = subtotal;

    // Demo only: there is no server, so nothing is sent or charged
    setTimeout(() => {
      const orderNo = 'WSL-' + Math.floor(100000 + Math.random() * 900000);
      setDone({
        orderNo,
        first: form.coName.trim().split(/\s+/)[0],
        address: `${form.coAddress.trim()}, ${form.coArea.trim()}, ${form.coCity.trim()}`,
        payText:
          pay === 'card'
            ? 'Card ending ' + form.cardNumber.replace(/\D/g, '').slice(-4)
            : 'Cash on delivery',
        total,
      });
      clear(); // empty the cart (card details are never saved)
      setPlacing(false);
      window.scrollTo({ top: 0, behavior: 'smooth' });
    }, 700);
  }

  const err = (id: FieldId) => errors[id] ?? '';
  const cls = (id: FieldId) => (err(id) ? 'invalid' : '');

  const crumbs = (
    <nav className="crumbs" aria-label="Breadcrumb">
      <Link to="/">Home</Link>
      <span aria-hidden="true">&rsaquo;</span>
      <Link to="/cart">Your cart</Link>
      <span aria-hidden="true">&rsaquo;</span>
      <span className="cur">Checkout</span>
    </nav>
  );

  /* ===== Order confirmed ===== */
  if (done) {
    return (
      <main className="cart-page">
        {crumbs}
        <h1 className="cart-title">Order confirmed</h1>
        <div className="cart-layout done">
          <section className="co-card">
            <div className="co-done">
              <div className="done-icon"><CheckIcon /></div>
              <h2>Thank you, {done.first}!</h2>
              <p>Your order has been placed successfully.</p>
              <dl className="done-details">
                <div><dt>Order number</dt><dd>{done.orderNo}</dd></div>
                <div><dt>Deliver to</dt><dd>{done.address}</dd></div>
                <div><dt>Payment</dt><dd>{done.payText}</dd></div>
                <div><dt>Total</dt><dd>{money(done.total)}</dd></div>
              </dl>
              <Link className="btn btn-solid" to="/">Back to home</Link>
            </div>
          </section>
        </div>
      </main>
    );
  }

  /* ===== Empty cart: nothing to check out ===== */
  if (cart.length === 0) {
    return (
      <main className="cart-page">
        {crumbs}
        <h1 className="cart-title">Checkout</h1>
        <div className="cart-layout done">
          <section className="co-card">
            <div className="co-done">
              <h2>Your cart is empty</h2>
              <p>Add something to your cart before checking out.</p>
              <Link className="btn btn-solid" to="/">Start shopping</Link>
            </div>
          </section>
        </div>
      </main>
    );
  }

  return (
    <main className="cart-page">
      {crumbs}
      <h1 className="cart-title">Checkout</h1>

      <div className="cart-layout">
        <section className="co-card">
          <form id="checkoutForm" onSubmit={handleSubmit} noValidate>
            <div className="co-section">
              <h2 className="co-sec-title"><span className="co-num">1</span>Contact details</h2>
              <div className="co-grid">
                <div className="field full">
                  <label htmlFor="coName">Full name</label>
                  <input id="coName" type="text" placeholder="e.g. Ahmed Mohamed" autoComplete="name"
                    className={cls('coName')} value={form.coName} onChange={onText('coName')} />
                  <span className="error">{err('coName')}</span>
                </div>
                <div className="field">
                  <label htmlFor="coPhone">Phone number</label>
                  <input id="coPhone" type="tel" placeholder="01XXXXXXXXX" autoComplete="tel"
                    className={cls('coPhone')} value={form.coPhone} onChange={onText('coPhone')} />
                  <span className="error">{err('coPhone')}</span>
                </div>
                <div className="field">
                  <label htmlFor="coEmail">Email</label>
                  <input id="coEmail" type="email" placeholder="you@example.com" autoComplete="email"
                    className={cls('coEmail')} value={form.coEmail} onChange={onText('coEmail')} />
                  <span className="error">{err('coEmail')}</span>
                </div>
              </div>
            </div>

            <div className="co-section">
              <h2 className="co-sec-title"><span className="co-num">2</span>Delivery address</h2>
              <div className="co-grid">
                <div className="field">
                  <label htmlFor="coCity">City</label>
                  <input id="coCity" type="text" placeholder="e.g. Giza" autoComplete="address-level2"
                    className={cls('coCity')} value={form.coCity} onChange={onText('coCity')} />
                  <span className="error">{err('coCity')}</span>
                </div>
                <div className="field">
                  <label htmlFor="coArea">Area / District</label>
                  <input id="coArea" type="text" placeholder="e.g. Dokki" autoComplete="address-level3"
                    className={cls('coArea')} value={form.coArea} onChange={onText('coArea')} />
                  <span className="error">{err('coArea')}</span>
                </div>
                <div className="field full">
                  <label htmlFor="coAddress">Street address</label>
                  <input id="coAddress" type="text" placeholder="Street, building, floor, apartment"
                    autoComplete="street-address" className={cls('coAddress')} value={form.coAddress}
                    onChange={onText('coAddress')} />
                  <span className="error">{err('coAddress')}</span>
                </div>
                <div className="field full">
                  <label htmlFor="coNotes">Delivery notes <span className="opt">(optional)</span></label>
                  <textarea id="coNotes" rows={2} placeholder="Landmark, best time to call, etc."
                    value={form.coNotes} onChange={onText('coNotes')} />
                </div>
              </div>
            </div>

            <div className="co-section">
              <h2 className="co-sec-title"><span className="co-num">3</span>Payment method</h2>
              <div className="pay-options">
                <label className={'pay-opt' + (pay === 'cash' ? ' selected' : '')}>
                  <input type="radio" name="pay" value="cash" checked={pay === 'cash'}
                    onChange={() => setPay('cash')} />
                  <CashIcon />
                  <span><b>Cash on delivery</b><small>Pay when you receive</small></span>
                </label>
                <label className={'pay-opt' + (pay === 'card' ? ' selected' : '')}>
                  <input type="radio" name="pay" value="card" checked={pay === 'card'}
                    onChange={() => setPay('card')} />
                  <CardIcon />
                  <span><b>Credit / Debit card</b><small>Visa, Mastercard</small></span>
                </label>
              </div>

              {pay === 'card' && (
                <div className="co-grid card-grid">
                  <div className="field full">
                    <label htmlFor="cardNumber">Card number</label>
                    <input id="cardNumber" type="text" inputMode="numeric" placeholder="1234 5678 9012 3456"
                      autoComplete="cc-number" maxLength={23} className={cls('cardNumber')}
                      value={form.cardNumber} onChange={onCardNumber} />
                    <span className="error">{err('cardNumber')}</span>
                  </div>
                  <div className="field full">
                    <label htmlFor="cardName">Name on card</label>
                    <input id="cardName" type="text" placeholder="As printed on the card"
                      autoComplete="cc-name" className={cls('cardName')} value={form.cardName}
                      onChange={onText('cardName')} />
                    <span className="error">{err('cardName')}</span>
                  </div>
                  <div className="field">
                    <label htmlFor="cardExpiry">Expiry date</label>
                    <input id="cardExpiry" type="text" inputMode="numeric" placeholder="MM/YY"
                      autoComplete="cc-exp" maxLength={5} className={cls('cardExpiry')}
                      value={form.cardExpiry} onChange={onExpiry} />
                    <span className="error">{err('cardExpiry')}</span>
                  </div>
                  <div className="field">
                    <label htmlFor="cardCvv">CVV</label>
                    <input id="cardCvv" type="password" inputMode="numeric" placeholder="•••"
                      autoComplete="cc-csc" maxLength={4} className={cls('cardCvv')}
                      value={form.cardCvv} onChange={onCvv} />
                    <span className="error">{err('cardCvv')}</span>
                  </div>
                </div>
              )}
            </div>
          </form>
        </section>

        <aside className="cart-summary" aria-label="Order summary">
          <h2>Your order</h2>
          <div>
            {cart.map((i) => (
              <div className="mini-item" key={i.id}>
                <span>{i.name}<small>{i.qty} × {money(i.price)}</small></span>
                <span>{money(i.price * i.qty)}</span>
              </div>
            ))}
          </div>
          <div className="sum-total"><span className="serif">Total</span><strong>{money(subtotal)}</strong></div>
          <p className="sum-note">Taxes and shipping are not added yet.</p>
          <button type="submit" form="checkoutForm" className="checkout" disabled={placing}>
            {placing ? 'Placing order…' : 'Place order'}
          </button>
          <p className="secure-note">&#128274; Your details are only used to place this order.</p>
        </aside>
      </div>
    </main>
  );
}
