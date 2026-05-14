import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function AuthCallback() {
  const [searchParams] = useSearchParams();
  const [error, setError] = useState('');
  const navigate = useNavigate();
  const { completeOAuth } = useAuth();

  useEffect(() => {
    const token = searchParams.get('token');
    const oauthError = searchParams.get('error');

    if (oauthError) {
      setError(oauthError);
      return;
    }
    if (!token) {
      setError('Missing login token');
      return;
    }

    completeOAuth(token)
      .then(() => navigate('/products', { replace: true }))
      .catch(() => setError('OAuth login could not be completed'));
  }, [completeOAuth, navigate, searchParams]);

  return (
    <div className="page center">
      {error ? (
        <div className="success-box">
          <h2>Login failed</h2>
          <p className="status error">{error}</p>
        </div>
      ) : (
        <p className="status">Completing login...</p>
      )}
    </div>
  );
}
