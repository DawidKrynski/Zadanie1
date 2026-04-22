import { useState } from 'react';
import { useProducts } from '../hooks/useProducts';
import { useCart } from '../context/CartContext';

export default function Products() {
  const [filters, setFilters] = useState({ in_stock: '', min_price: '', max_price: '' });
  const { products, loading, error } = useProducts(
    // Only pass non-empty filter values
    Object.fromEntries(Object.entries(filters).filter(([, v]) => v !== ''))
  );
  const { addItem } = useCart();
  const [added, setAdded] = useState(null);

  async function handleAdd(product) {
    await addItem(product.ID);
    setAdded(product.ID);
    setTimeout(() => setAdded(null), 1500);
  }

  return (
    <div className="page">
      <h1>Products</h1>

      {/* Filters */}
      <div className="filters">
        <label>
          Min price
          <input
            type="number"
            value={filters.min_price}
            onChange={e => setFilters(f => ({ ...f, min_price: e.target.value }))}
            placeholder="0"
          />
        </label>
        <label>
          Max price
          <input
            type="number"
            value={filters.max_price}
            onChange={e => setFilters(f => ({ ...f, max_price: e.target.value }))}
            placeholder="9999"
          />
        </label>
        <label className="checkbox-label">
          <input
            type="checkbox"
            checked={filters.in_stock === 'true'}
            onChange={e => setFilters(f => ({ ...f, in_stock: e.target.checked ? 'true' : '' }))}
          />
          In stock only
        </label>
      </div>

      {loading && <p className="status">Loading…</p>}
      {error && <p className="status error">{error}</p>}

      <div className="product-grid">
        {products.map(product => (
          <div key={product.ID} className="card">
            <h3>{product.name}</h3>
            {product.category && <span className="badge">{product.category.name}</span>}
            <p className="desc">{product.description}</p>
            <div className="card-footer">
              <strong>${product.price.toFixed(2)}</strong>
              <span className={product.quantity > 0 ? 'in-stock' : 'out-stock'}>
                {product.quantity > 0 ? `${product.quantity} in stock` : 'Out of stock'}
              </span>
              <button
                disabled={product.quantity === 0 || added === product.ID}
                onClick={() => handleAdd(product)}
                className="btn-primary"
              >
                {added === product.ID ? '✓ Added' : 'Add to cart'}
              </button>
            </div>
          </div>
        ))}
      </div>

      {!loading && products.length === 0 && (
        <p className="status">No products found.</p>
      )}
    </div>
  );
}
