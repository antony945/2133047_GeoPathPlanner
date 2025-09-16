import { useState } from "react";
import './Sidebar.css';

function Sidebar() {
  const [tab, setTab] = useState("waypoint");

  return (
    <div className="d-flex flex-column h-100 sidebar">
      {/* Nav-tabs */}
      <ul className="nav nav-tabs">
        <li className="nav-item">
          <button className={`nav-link${tab === "waypoint" ? " active" : ""}`} onClick={() => setTab("waypoint")}>
            Waypoints
          </button>
        </li>
        <li className="nav-item">
          <button className={`nav-link${tab === "obstacles" ? " active" : ""}`} onClick={() => setTab("obstacles")}>
            Obstacles
          </button>
        </li>
      </ul>
      {/* Tab content */}
      <div className="flex-fill overflow-auto p-3" style={{ minHeight: 0 }}>
        {tab === "waypoint" ? <WaypointPanel /> : <ObstaclePanel />}
      </div>
    </div>
  );
}

function WaypointPanel() {
  return (
    <form className="text-start">
      <div className="mb-3">
        <label className="form-label">Unità di misura altezza</label>
        <select className="form-select">
          <option value="m">Metri (m)</option>
          <option value="ft">Piedi (ft)</option>
        </select>
      </div>
      <div className="mb-3">
        <label className="form-label">Cerca località</label>
        <input type="text" className="form-control mb-2" placeholder="Nome o coordinate" />
        <div className="d-flex gap-2">
          <button type="button" className="btn btn-primary">
            <i className="bi bi-geo-alt"></i> Marker
          </button>
          <button type="button" className="btn btn-outline-secondary">
            <i className="bi bi-crosshair"></i> Geolocalizza
          </button>
        </div>
      </div>
      <hr />
      <div className="mb-3">
        <label className="form-label">Importa waypoints da file (.geojson)</label>
        <input type="file" className="form-control" accept=".geojson,.json" />
      </div>
    </form>
  );
}

function ObstaclePanel() {
  return (
    <form className="text-start">
      <div className="mb-3">
        <button type="button" className="btn btn-outline-primary w-100 mb-2">
          <i className="bi bi-vector-pen"></i> Disegna poligono
        </button>
        <button type="button" className="btn btn-outline-success w-100">
          <i className="bi bi-magic"></i> Genera ostacoli a caso
        </button>
      </div>
      <hr />
      <div className="mb-3">
        <label className="form-label">Importa ostacoli da file (.geojson)</label>
        <input type="file" className="form-control" accept=".geojson,.json" />
      </div>
    </form>
  );
}

export default Sidebar;