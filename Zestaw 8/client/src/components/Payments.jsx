import { useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { useCart } from '../context/CartContext';
import api from '../api';

export default function Payments() {
  const { state: routeState } = useLocation();
  const { cart: ctxCart, cartTotal: ctxTotal, clearCart } = useCart();
  const cart = routeState?.cart ?? ctxCart;
  const cartTotal = routeState?.cartTotal ?? ctxTotal;
  const navigate = useNavigate();

  const [form, setForm] = useState({ name: '', email: '', address: '' });
  const [status, setStatus] = useState(null);
  const [error, setError] = useState('');

  function handleChange(e) {
    setForm(f => ({ ...f, [e.target.name]: e.target.value }));
  }

  async function handleSubmit(e) {
    e.preventDefault();
    if (!cart) return;

    setStatus('loading');
    setError('');

    const payload = {
      cart_id: cart.ID,
      customer: form,
      items: cart.items?.map(item => ({
        product_id: item.product_id,
        quantity: item.quantity,
      })),
    };

    try {
      await api.post('/orders', payload);
      setStatus('success');
      await clearCart();
    } catch (err) {
      setError(err.response?.data?.error ?? 'Payment failed. Try again.');
      setStatus('error');
    }
  }

  if (status === 'success') {
    return (
      <div className="page center">
        <div className="success-box">
          <h2>Order placed</h2>
          <p>Thank you, {form.name}. A confirmation will be sent to {form.email}.</p>
          <button className="btn-primary" onClick={() => navigate('/products')}>
            Continue shopping
          </button>
        </div>
      </div>
    );
  }

  if (!cart || !cart.items?.length) {
    return (
      <div className="page center">
        <p>Your cart is empty. <button className="btn-link" onClick={() => navigate('/products')}>Go shopping</button></p>
      </div>
    );
  }

  return (
    <div className="page">
      <h1>Payment</h1>

      <div className="order-summary">
        <h3>Order summary</h3>
        {cart.items.map(item => (
          <div key={item.ID} className="summary-row">
            <span>{item.product?.name} x {item.quantity}</span>
            <span>${(item.product?.price * item.quantity).toFixed(2)}</span>
          </div>
        ))}
        <div className="summary-row total">
          <strong>Total</strong>
          <strong>${cartTotal.toFixed(2)}</strong>
        </div>
      </div>

      <form onSubmit={handleSubmit} className="payment-form">
        <h3>Billing details</h3>

        <label>
          Full name
          <input name="name" value={form.name} onChange={handleChange} required placeholder="Jan Kowalski" />
        </label>

        <label>
          Email
          <input name="email" type="email" value={form.email} onChange={handleChange} required placeholder="jan@example.com" />
        </label>

        <label>
          Address
          <input name="address" value={form.address} onChange={handleChange} required placeholder="ul. Glowna 1, 00-001 Warszawa" />
        </label>

        {error && <p className="status error">{error}</p>}

        <button className="btn-primary" type="submit" disabled={status === 'loading'}>
          {status === 'loading' ? 'Processing...' : `Pay $${cartTotal.toFixed(2)}`}
        </button>
      </form>
    </div>
  );
}
