import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import axios from 'axios';

const Signup = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    if (password !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }
    try {
      await axios.post('http://localhost:8080/signup', {
        email,
        password,
      });
      navigate('/login');
    } catch (err: any) {
      if (err.response && err.response.status === 409) {
        setError('This email is already registered.');
      } else {
        setError('Something went wrong. Please try again.');
      }
    }
  };

  return (
    <div className="hero">
      <div className="glass-card" style={{ maxWidth: '420px' }}>
        <div className="logo" style={{ fontSize: '2rem' }}>Rayan Academy</div>
        <h2 style={{ marginBottom: '20px' }}>Create Your Account</h2>
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
          <input
            className="glass-input"
            type="password"
            placeholder="Confirm Password"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            required
          />
          <button type="submit" className="btn-primary" style={{ width: '100%', marginTop: '10px' }}>
            Sign Up
          </button>
        </form>
        <p style={{ marginTop: '20px', color: 'rgba(255,255,255,0.8)' }}>
          Already have an account? <Link to="/login" className="violet-text">Log in</Link>
        </p>
      </div>
    </div>
  );
};

export default Signup;