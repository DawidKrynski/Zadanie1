import { useCart } from '../context/CartContext';
import { useNavigate } from 'react-router-dom';

export default function Cart() {
  const { cart, cartTotal, loading, updateItem, removeItem } = useCart();
  const navigate = useNavigate();

  if (loading) return <p className="status">Loading cart...</p>;
  if (!cart) return <p className="status">No cart available.</p>;

  const items = cart.items ?? [];

  return (
    <div className="page">
      <h1>Cart</h1>

      {items.length === 0 ? (
        <p className="status">Your cart is empty.</p>
      ) : (
        <>
          <table className="cart-table">
            <thead>
              <tr>
                <th>Product</th>
                <th>Price</th>
                <th>Qty</th>
                <th>Subtotal</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {items.map(item => {
                const price = item.product?.price ?? 0;

                return (
                  <tr key={item.ID}>
                    <td>{item.product?.name ?? `#${item.product_id}`}</td>
                    <td>${price.toFixed(2)}</td>
                    <td>
                      <div className="qty-ctrl">
                        <button
                          type="button"
                          onClick={() => updateItem(item.ID, Math.max(1, item.quantity - 1))}
                        >
                          -
                        </button>
                        <span>{item.quantity}</span>
                        <button
                          type="button"
                          onClick={() => updateItem(item.ID, item.quantity + 1)}
                        >
                          +
                        </button>
                      </div>
                    </td>
                    <td>${(price * item.quantity).toFixed(2)}</td>
                    <td>
                      <button
                        type="button"
                        className="btn-danger"
                        onClick={() => removeItem(item.ID)}
                      >
                        Remove
                      </button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>

          <div className="cart-footer">
            <strong>Total: ${cartTotal.toFixed(2)}</strong>
            <button
              type="button"
              className="btn-primary"
              onClick={() => navigate('/payments', { state: { cart, cartTotal } })}
            >
              Proceed to payment
            </button>
          </div>
        </>
      )}
    </div>
  );
}
