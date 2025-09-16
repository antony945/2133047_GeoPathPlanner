import { BrowserRouter, Routes, Route, Link, Navigate, NavLink } from 'react-router-dom';
import Homepage from './pages/Homepage/Homepage';

function App() {
   return (
      <BrowserRouter>
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
               <i className="bi bi-geo-alt-fill" /> {/* If using Bootstrap Icons, otherwise use emoji/map icon */}
               </span>
               GeoPathPlanner
            </Link>
            {/* Collapse/toggler for mobile */}
            <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarGeoPath" aria-controls="navbarGeoPath" aria-expanded="false" aria-label="Toggle navigation">
               <span className="navbar-toggler-icon" />
            </button>
            <div className="collapse navbar-collapse" id="navbarGeoPath">
               <ul className="navbar-nav ms-auto mb-2 mb-lg-0">
               <li className="nav-item">
                  <NavLink to="/" className="nav-link" style={({ isActive }) => ({ color: isActive ? "#4FC0D0" : "#fff" })}>
                     Home
                  </NavLink>
               </li>
               <li className="nav-item">
                  <NavLink to="/about" className="nav-link" style={({ isActive }) => ({ color: isActive ? "#4FC0D0" : "#fff" })}>
                     About
                  </NavLink>
               </li>
               <li className="nav-item">
                  <NavLink to="/dashboard" className="nav-link" style={({ isActive }) => ({ color: isActive ? "#4FC0D0" : "#fff" })}>
                     Dashboard
                  </NavLink>
               </li>
               {/* Future user dropdown/avatar: */}
               {/* <li className="nav-item dropdown">
                  <a className="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false" style={{ color: "#fff" }}>
                     <i className="bi bi-person-circle"></i> User
                  </a>
                  <ul className="dropdown-menu dropdown-menu-end">
                     <li><Link className="dropdown-item" to="/profile">Profile</Link></li>
                     <li><button className="dropdown-item">Logout</button></li>
                  </ul>
               </li> */}
               </ul>
            </div>
         </div>
         </nav>
         <Routes>
            <Route path="/" element={<Homepage/>} />
            <Route path="/about" element={<h1>About Page</h1>} />
            <Route path="/dashboard" element={<h1>Dashboard Page</h1>} />
            <Route path="*" element={<Navigate to="/" replace />} />
         </Routes>
      </BrowserRouter>
  )
}

export default App
