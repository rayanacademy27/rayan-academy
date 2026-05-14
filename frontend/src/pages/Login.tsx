import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import axios from 'axios';

const Login = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      const res = await axios.post('http://localhost:8080/login', {
        email,
        password,
      });
      localStorage.setItem('token', res.data.token);
      navigate('/');
    } catch (err) {
      setError('Invalid email or password. Please try again.');
    }
  };

  return (
    <div className="hero">
      <div className="glass-card" style={{ maxWidth: '420px' }}>
        <div className="logo" style={{ fontSize: '2rem' }}>Rayan Academy</div>
        <h2 style={{ marginBottom: '20px' }}>Welcome Back</h2>
        <form onSubmit={handleSubmit}>
          {error && <p className="error-text">{error}</p>}
          <input
            className="glass-input"
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          <input
            className="glass-input"
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          <button type="submit" className="btn-primary" style={{ width: '100%', marginTop: '10px' }}>
            Log In
          </button>
        </form>
        <p style={{ marginTop: '20px', color: 'rgba(255,255,255,0.8)' }}>
          Don't have an account? <Link to="/signup" className="violet-text">Sign up</Link>
        </p>
      </div>
    </div>
  );
};

export default Login;