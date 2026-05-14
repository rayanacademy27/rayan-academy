import { Link } from 'react-router-dom';

function Home() {
  const token = localStorage.getItem('token');
  
  // Helper function to extract email from the JWT token
  const getEmailFromToken = () => {
    if (!token) return "User";
    try {
      // JWTs are divided into 3 parts by dots. The 2nd part [1] contains the data.
      const base64Url = token.split('.')[1];
      const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
      const jsonPayload = decodeURIComponent(window.atob(base64).split('').map(function(c) {
          return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
      }).join(''));

      return JSON.parse(jsonPayload).email;
    } catch (e) {
      console.error("Token decoding failed", e);
      return "User";
    }
  };

  const email = getEmailFromToken();

  const handleLogout = () => {
    localStorage.removeItem('token');
    window.location.reload(); // Refresh the page to reset the UI
  };

  return (
    <div className="hero">
      <div className="glass-card">
        <div className="logo">Rayan Academy</div>
        <h1>Master Government Exams</h1>
        <p className="tagline">
          Live classes • Smart mocks • AI analytics<br />
          <span className="violet-text">Your success, now in glass.</span>
        </p>

        {/* Show welcome message only if logged in */}
        {token && (
          <p style={{ color: 'white', marginBottom: '20px', fontWeight: '500' }}>
            Welcome back, <span className="violet-text">{email}</span>
          </p>
        )}

        <div className="cta-buttons">
          {token ? (
            <div style={{ display: 'flex', gap: '15px', justifyContent: 'center' }}>
              <Link to="/dashboard">
                <button className="btn-primary">Go to Dashboard</button>
              </Link>
              <button className="btn-secondary" onClick={handleLogout}>
                Logout
              </button>
            </div>
          ) : (
            <>
              <Link to="/signup">
                <button className="btn-primary">Get Started</button>
              </Link>
              <Link to="/login">
                <button className="btn-secondary" style={{ marginLeft: '10px' }}>
                  Explore Courses
                </button>
              </Link>
            </>
          )}
        </div>

        {/* Quick links to Mock Tests */}
        <div style={{ marginTop: '30px' }}>
            <p style={{ color: 'rgba(255,255,255,0.6)', fontSize: '0.9rem', marginBottom: '15px' }}>
                Featured Mock Tests
            </p>
            <div style={{ display: 'flex', gap: '10px', justifyContent: 'center' }}>
                <Link to="/test/ssc-cgl-mock-1">
                    <button className="btn-secondary" style={{ fontSize: '0.8rem' }}>SSC CGL Mock</button>
                </Link>
                <Link to="/test/bank-po-1">
                    <button className="btn-secondary" style={{ fontSize: '0.8rem' }}>Bank PO Mock</button>
                </Link>
            </div>
        </div>
      </div>

      <div className="floating-shapes">
        <div className="shape shape-1"></div>
        <div className="shape shape-2"></div>
        <div className="shape shape-3"></div>
      </div>
    </div>
  );
}

export default Home;