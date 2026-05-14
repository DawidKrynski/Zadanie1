import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function Login() {
  const [mode, setMode] = useState('login');
  const [form, setForm] = useState({ name: '', email: '', password: '' });
  const [localError, setLocalError] = useState('');
  const navigate = useNavigate();
  const { login, register, oauthLogin, loading, error } = useAuth();

  function handleChange(e) {
    setForm(current => ({ ...current, [e.target.name]: e.target.value }));
  }

  async function handleSubmit(e) {
    e.preventDefault();
    setLocalError('');

    try {
      if (mode === 'register') {
        await register(form);
      } else {
        await login({ email: form.email, password: form.password });
      }
      navigate('/products', { replace: true });
    } catch (err) {
      setLocalError(err.message);
    }
  }

  const visibleError = localError || error;

  return (
    <div className="auth-layout">
      <section className="auth-panel">
        <div className="auth-switch" role="tablist" aria-label="Authentication mode">
          <button
            type="button"
            className={mode === 'login' ? 'active' : ''}
            onClick={() => setMode('login')}
          >
            Login
          </button>
          <button
            type="button"
            className={mode === 'register' ? 'active' : ''}
            onClick={() => setMode('register')}
          >
            Register
          </button>
        </div>

        <form className="auth-form" onSubmit={handleSubmit}>
          {mode === 'register' && (
            <label>
              Name
              <input name="name" value={form.name} onChange={handleChange} placeholder="Jan Kowalski" />
            </label>
          )}

          <label>
            Email
            <input
              name="email"
              type="email"
              value={form.email}
              onChange={handleChange}
              placeholder="jan@example.com"
              required
            />
          </label>

          <label>
            Password
            <input
              name="password"
              type="password"
              value={form.password}
              onChange={handleChange}
              minLength={8}
              required
            />
          </label>

          {visibleError && <p className="status error">{visibleError}</p>}

          <button className="btn-primary" type="submit" disabled={loading}>
            {loading ? 'Please wait...' : mode === 'register' ? 'Create account' : 'Login'}
          </button>
        </form>

        <div className="oauth-row">
          <button type="button" className="btn-secondary" onClick={() => oauthLogin('google')}>
            Google
          </button>
          <button type="button" className="btn-secondary" onClick={() => oauthLogin('github')}>
            GitHub
          </button>
        </div>
      </section>
    </div>
  );
}
