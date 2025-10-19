import { BrowserRouter, Routes, Route, Link, Navigate, NavLink } from 'react-router-dom';
import Homepage from './pages/Homepage/Homepage';
import LoginPage from './pages/Login/Login';
import ProfilePage from './pages/Profile/Profile';
import ProtectedRoute from './components/ProtectedRoute/ProtectedRoute';
import { AuthProvider, useAuth } from './context/AuthContext';

// Navbar Component
function Navigation() {
  const { isAuthenticated, user, logout } = useAuth();

  return (
    <nav className="navbar navbar-expand-lg" style={{ background: "#176B87" }}>
      <div className="container">
        {/* Brand */}
        <Link className="navbar-brand d-flex align-items-center" to="/" style={{ color: "#fff", fontWeight: 'bold', letterSpacing: 1 }}>
          <span
            className="me-2 d-inline-flex align-items-center justify-content-center"
            style={{
              width: 38,
              height: 38,
              background: "#4FC0D0",
              borderRadius: "50%",
              color: "#fff",
              fontWeight: "bold",
              fontSize: 20,
            }}
          >
            <i className="bi bi-geo-alt-fill" />
          </span>
          GeoPathPlanner
        </Link>
        {/* Toggler */}
        <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarGeoPath" aria-controls="navbarGeoPath" aria-expanded="false" aria-label="Toggle navigation">
          <span className="navbar-toggler-icon" />
        </button>
        {/* Links */}
        <div className="collapse navbar-collapse" id="navbarGeoPath">
          <ul className="navbar-nav ms-auto mb-2 mb-lg-0 align-items-center">
            <li className="nav-item">
              <NavLink to="/" className="nav-link" style={({ isActive }) => ({ color: isActive ? "#4FC0D0" : "#fff" })}>
                Home
              </NavLink>
            </li>

            
            {isAuthenticated ? (
              <li className="nav-item dropdown">
                <a className="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false" style={{ color: "#fff" }}>
                  <i className="bi bi-person-circle me-1"></i> {user?.username}
                </a>
                <ul className="dropdown-menu dropdown-menu-end">
                  <li><Link className="dropdown-item" to="/profile">Profile</Link></li>
                  <li><hr className="dropdown-divider" /></li>
                  <li><button className="dropdown-item" onClick={logout}>Logout</button></li>
                </ul>
              </li>
            ) : (
              <li className="nav-item">
                <NavLink to="/login" className="btn btn-outline-light" style={{'--bs-btn-hover-bg': '#4FC0D0', '--bs-btn-hover-border-color': '#4FC0D0'}}>
                  Login / Register
                </NavLink>
              </li>
            )}
          </ul>
        </div>
      </div>
    </nav>
  );
}

// Main App Layout
function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Navigation />
        <Routes>
          {/* Public Routes */}
          <Route path="/" element={<Homepage />} />
          <Route path="/login" element={<LoginPage />} />
          
          {/* Protected Routes */}

          <Route 
            path="/profile" 
            element={
              <ProtectedRoute>
                <ProfilePage />
              </ProtectedRoute>
            } 
          />

          {/* Fallback Route */}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;

