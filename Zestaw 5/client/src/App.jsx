import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { CartProvider } from './context/CartContext';
import Navbar from './components/Navbar';
import Products from './components/Products';
import Cart from './components/Cart';
import Payments from './components/Payments';
import './index.css';

export default function App() {
  return (
    // CartProvider wraps everything so all components share cart state via useCart hook
    <CartProvider>
      <BrowserRouter>
        <Navbar />
        <main>
          <Routes>
            <Route path="/" element={<Navigate to="/products" replace />} />
            <Route path="/products" element={<Products />} />
            <Route path="/cart" element={<Cart />} />
            <Route path="/payments" element={<Payments />} />
          </Routes>
        </main>
      </BrowserRouter>
    </CartProvider>
  );
}
