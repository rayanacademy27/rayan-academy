import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './App.css';

function Home() {
  return (
    <div className="hero">
      <div className="glass-card">
        <div className="logo">Rayan Academy</div>
        <h1>Master Government Exams</h1>
        <p className="tagline">
          Live classes • Smart mocks • AI analytics<br />
          <span className="violet-text">Your success, now in glass.</span>
        </p>
        <div className="cta-buttons">
          <button className="btn-primary">Get Started</button>
          <button className="btn-secondary">Explore Courses</button>
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

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        {/* More routes will be added day by day */}
      </Routes>
    </Router>
  );
}

export default App;