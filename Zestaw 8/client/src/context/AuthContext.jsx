import { createContext, useCallback, useContext, useEffect, useReducer } from 'react';
import api, { AUTH_TOKEN_KEY } from '../api';

const AuthContext = createContext(null);

function authReducer(state, action) {
  switch (action.type) {
    case 'SET_LOADING':
      return { ...state, loading: action.payload };
    case 'SET_AUTH':
      return {
        ...state,
        token: action.payload.token,
        user: action.payload.user,
        loading: false,
        error: null,
      };
    case 'SET_USER':
      return { ...state, user: action.payload, loading: false, error: null };
    case 'SET_ERROR':
      return { ...state, error: action.payload, loading: false };
    case 'CLEAR_AUTH':
      return { token: null, user: null, loading: false, error: null };
    default:
      return state;
  }
}

const initialToken = localStorage.getItem(AUTH_TOKEN_KEY);

const initialState = {
  token: initialToken,
  user: null,
  loading: Boolean(initialToken),
  error: null,
};

export function AuthProvider({ children }) {
  const [state, dispatch] = useReducer(authReducer, initialState);

  useEffect(() => {
    if (!state.token) {
      dispatch({ type: 'SET_LOADING', payload: false });
      return;
    }

    api.get('/auth/me')
      .then(({ data }) => dispatch({ type: 'SET_USER', payload: data }))
      .catch(() => {
        localStorage.removeItem(AUTH_TOKEN_KEY);
        dispatch({ type: 'CLEAR_AUTH' });
      });
  }, [state.token]);

  const login = useCallback(async (credentials) => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      const { data } = await api.post('/auth/login', credentials);
      localStorage.setItem(AUTH_TOKEN_KEY, data.token);
      dispatch({ type: 'SET_AUTH', payload: data });
      return data.user;
    } catch (err) {
      const message = err.response?.data?.error ?? 'Login failed';
      dispatch({ type: 'SET_ERROR', payload: message });
      throw new Error(message);
    }
  }, []);

  const register = useCallback(async (credentials) => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      const { data } = await api.post('/auth/register', credentials);
      localStorage.setItem(AUTH_TOKEN_KEY, data.token);
      dispatch({ type: 'SET_AUTH', payload: data });
      return data.user;
    } catch (err) {
      const message = err.response?.data?.error ?? 'Registration failed';
      dispatch({ type: 'SET_ERROR', payload: message });
      throw new Error(message);
    }
  }, []);

  const completeOAuth = useCallback(async (token) => {
    localStorage.setItem(AUTH_TOKEN_KEY, token);
    dispatch({ type: 'SET_AUTH', payload: { token, user: null } });

    const { data } = await api.get('/auth/me');
    dispatch({ type: 'SET_USER', payload: data });
    return data;
  }, []);

  const logout = useCallback(async () => {
    try {
      await api.post('/auth/logout');
    } finally {
      localStorage.removeItem(AUTH_TOKEN_KEY);
      dispatch({ type: 'CLEAR_AUTH' });
    }
  }, []);

  const oauthLogin = useCallback((provider) => {
    window.location.href = `${api.defaults.baseURL}/auth/${provider}/login`;
  }, []);

  return (
    <AuthContext.Provider value={{ ...state, login, register, logout, oauthLogin, completeOAuth }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
