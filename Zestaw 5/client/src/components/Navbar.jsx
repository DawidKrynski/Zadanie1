import { NavLink } from 'react-router-dom';
import { useCart } from '../context/CartContext';

export default function Navbar() {
  const { cart } = useCart();
  const itemCount = cart?.items?.reduce((sum, i) => sum + i.quantity, 0) ?? 0;

  return (
    <nav className="navbar">
      <span className="logo">🛒 Shop</span>
      <div className="nav-links">
        <NavLink to="/products" className={({ isActive }) => isActive ? 'active' : ''}>Products</NavLink>
        <NavLink to="/cart" className={({ isActive }) => isActive ? 'active' : ''}>
          Cart {itemCount > 0 && <span className="badge-count">{itemCount}</span>}
        </NavLink>
        <NavLink to="/payments" className={({ isActive }) => isActive ? 'active' : ''}>Payments</NavLink>
      </div>
    </nav>
  );
}
