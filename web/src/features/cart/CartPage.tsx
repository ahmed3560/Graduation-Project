import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useCart } from '../../hooks/useCart';
import { money } from '../../lib/format';
import { BoxIcon, ITEM_ICONS, TrashIcon } from '../../components/common/icons';

export default function CartPage() {
  const { cart, count, subtotal, increase, decrease, remove } = useCart();
  const navigate = useNavigate();
  const [terms, setTerms] = useState(false);
  const [termsError, setTermsError] = useState('');

  function handleCheckout() {
    if (!terms) {
      setTermsError('Please agree to the terms first.');
      return;
    }
    navigate('/checkout');
  }

  return (
    <main className="cart-page">
      <nav className="crumbs" aria-label="Breadcrumb">
        <Link to="/">Home</Link>
        <span aria-hidden="true">&rsaquo;</span>
        <span className="cur">Your cart</span>
      </nav>
      <h1 className="cart-title">Your cart</h1>

      <div className="cart-layout">
        <section className="cart-list-card" aria-label="Cart items">
          <div className="cart-head" aria-hidden="true">
            <span>Product</span>
            <span>Price</span>
            <span>Quantity</span>
            <span className="right">Total</span>
          </div>

          {cart.length === 0 ? (
            <div className="cart-empty">
              <BoxIcon />
              <h3>Your cart is empty</h3>
              <p>Looks like you haven&apos;t added anything yet.</p>
              <Link className="btn btn-solid" to="/">Start shopping</Link>
            </div>
          ) : (
            cart.map((item) => {
              const Icon = ITEM_ICONS[item.icon] ?? BoxIcon;
              return (
                <article className="cart-row" key={item.id}>
                  <div className="c-product">
                    <div className="c-thumb"><Icon /></div>
                    <div>
                      <h3>{item.name}</h3>
                      <p>{item.meta}</p>
                      <p className="c-unit">{money(item.price)} each</p>
                    </div>
                  </div>
                  <div className="c-price">{money(item.price)}</div>
                  <div className="c-qty">
                    <div className="stepper">
                      <button
                        type="button"
                        aria-label="Decrease quantity"
                        disabled={item.qty <= 1}
                        onClick={() => decrease(item.id)}
                      >
                        &minus;
                      </button>
                      <span aria-live="polite">{item.qty}</span>
                      <button type="button" aria-label="Increase quantity" onClick={() => increase(item.id)}>
                        +
                      </button>
                    </div>
                    <button
                      type="button"
                      className="c-remove"
                      aria-label={`Remove ${item.name}`}
                      onClick={() => remove(item.id)}
                    >
                      <TrashIcon />
                    </button>
                  </div>
                  <div className="c-total">{money(item.price * item.qty)}</div>
                </article>
              );
            })
          )}

          <Link className="continue" to="/">&larr; Continue shopping</Link>
        </section>

        <aside className="cart-summary" aria-label="Order summary">
          <h2>Order summary</h2>
          <div className="sum-row"><span>Items</span><span>{count}</span></div>
          <div className="sum-total">
            <span className="serif">Subtotal</span>
            <strong>{money(subtotal)}</strong>
          </div>
          <p className="sum-note">Taxes and shipping are calculated at checkout.</p>

          <label className="terms-row">
            <input
              type="checkbox"
              checked={terms}
              onChange={(e) => {
                setTerms(e.target.checked);
                setTermsError('');
              }}
            />
            <span>Agree to terms of sale as per the merchant&apos;s terms of service.</span>
          </label>
          <span className="error">{termsError}</span>

          <button type="button" className="checkout" disabled={cart.length === 0} onClick={handleCheckout}>
            Checkout
          </button>
        </aside>
      </div>
    </main>
  );
}
