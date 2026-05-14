import { createContext, useContext, useEffect, useReducer } from 'react';
import api from '../api';

const CartContext = createContext(null);

function cartReducer(state, action) {
  switch (action.type) {
    case 'SET_CART':
      return { ...state, cart: action.payload, loading: false };
    case 'SET_LOADING':
      return { ...state, loading: action.payload };
    case 'SET_ERROR':
      return { ...state, error: action.payload, loading: false };
    case 'CLEAR_CART':
      return { ...state, cart: null };
    default:
      return state;
  }
}

const initialState = {
  cart: null,
  loading: false,
  error: null,
};

export function CartProvider({ children }) {
  const [state, dispatch] = useReducer(cartReducer, initialState);

  useEffect(() => {
    const savedCartId = localStorage.getItem('cartId');
    if (savedCartId) {
      fetchCart(savedCartId);
    } else {
      createCart();
    }
  }, []);

  async function createCart() {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      const { data } = await api.post('/carts');
      localStorage.setItem('cartId', data.ID);
      dispatch({ type: 'SET_CART', payload: data });
    } catch {
      dispatch({ type: 'SET_ERROR', payload: 'Failed to create cart' });
    }
  }

  async function fetchCart(id) {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      const { data } = await api.get(`/carts/${id}`);
      dispatch({ type: 'SET_CART', payload: data });
    } catch {
      createCart();
    }
  }

  async function addItem(productId, quantity = 1) {
    if (!state.cart) return;
    const { data } = await api.post(`/carts/${state.cart.ID}/items`, {
      product_id: productId,
      quantity,
    });
    fetchCart(state.cart.ID);
    return data;
  }

  async function updateItem(itemId, quantity) {
    if (!state.cart) return;
    await api.put(`/carts/${state.cart.ID}/items/${itemId}`, { quantity });
    fetchCart(state.cart.ID);
  }

  async function removeItem(itemId) {
    if (!state.cart) return;
    await api.delete(`/carts/${state.cart.ID}/items/${itemId}`);
    fetchCart(state.cart.ID);
  }

  async function clearCart() {
    if (!state.cart) return;
    await api.delete(`/carts/${state.cart.ID}`);
    localStorage.removeItem('cartId');
    dispatch({ type: 'CLEAR_CART' });
    createCart();
  }

  const cartTotal = state.cart?.items?.reduce(
    (sum, item) => sum + item.product.price * item.quantity,
    0
  ) ?? 0;

  return (
    <CartContext.Provider value={{ ...state, cartTotal, addItem, updateItem, removeItem, clearCart }}>
      {children}
    </CartContext.Provider>
  );
}

export function useCart() {
  const ctx = useContext(CartContext);
  if (!ctx) throw new Error('useCart must be used within CartProvider');
  return ctx;
}
