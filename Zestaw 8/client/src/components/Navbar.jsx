import { NavLink } from 'react-router-dom';
import { useCart } from '../context/CartContext';
import { useAuth } from '../context/AuthContext';

export default function Navbar() {
  const { cart } = useCart();
  const { user, logout } = useAuth();
  const itemCount = cart?.items?.reduce((sum, i) => sum + i.quantity, 0) ?? 0;

  return (
    <nav className="navbar">
      <span className="logo">Shop</span>
      <div className="nav-links">
        <NavLink to="/products" className={({ isActive }) => isActive ? 'active' : ''}>Products</NavLink>
        <NavLink to="/cart" className={({ isActive }) => isActive ? 'active' : ''}>
          Cart {itemCount > 0 && <span className="badge-count">{itemCount}</span>}
        </NavLink>
        <NavLink to="/payments" className={({ isActive }) => isActive ? 'active' : ''}>Payments</NavLink>
        {user ? (
          <button type="button" className="nav-button" onClick={logout}>
            {(user.name || user.email)} / Logout
          </button>
        ) : (
          <NavLink to="/login" className={({ isActive }) => isActive ? 'active' : ''}>Login</NavLink>
        )}
      </div>
    </nav>
  );
}
