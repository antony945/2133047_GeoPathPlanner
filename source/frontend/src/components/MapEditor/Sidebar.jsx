import { useState, useEffect, useRef } from 'react';
import './Sidebar.css';
import { getCurrentPosition , geocodeNominatim } from '../../assets/js/utils/geolocation.js';

function Sidebar({ onGoto, onToggleDraw, onGenerateRandomObstacles, isMapReady, onAltitudeChange, onObstacleAltitudeChange, onCompute, isComputing, parameters, onParametersChange }) {
  const [tab, setTab] = useState('waypoint');
  const [altitudeValue, setAltitudeValue] = useState(0);
  const [altitudeUnit, setAltitudeUnit] = useState('mt');
  const [minAltitudeValue, setMinAltitudeValue] = useState(100);
  const [maxAltitudeValue, setMaxAltitudeValue] = useState(500);
  const [obstacleAltitudeUnit, setObstacleAltitudeUnit] = useState('mt');
  const [query, setQuery] = useState('');
  const [displayQuery, setDisplayQuery] = useState('');
  const [suggestions, setSuggestions] = useState([]);
  const [loading, setLoading] = useState(false);
  const searchTimer = useRef(null);

  useEffect(() => {
    if (onObstacleAltitudeChange) {
      onObstacleAltitudeChange({ min: minAltitudeValue, max: maxAltitudeValue, unit: obstacleAltitudeUnit });
    }
  }, [minAltitudeValue, maxAltitudeValue, obstacleAltitudeUnit, onObstacleAltitudeChange]);

  useEffect(() => {
    if (onAltitudeChange) {
      onAltitudeChange({ value: altitudeValue, unit: altitudeUnit });
    }
  }, [altitudeValue, altitudeUnit, onAltitudeChange]);

  useEffect(() => {
    if (!query) return setSuggestions([]);
    if (searchTimer.current) clearTimeout(searchTimer.current);
    
    searchTimer.current = setTimeout(async () => {
      setLoading(true);
      try {
        const data = await geocodeNominatim(query, { limit: 6 });
        setSuggestions(data.map(item => ({ display: item.display_name, lat: Number(item.lat), lon: Number(item.lon) })));
      } catch (err) {
        console.error('Geocode error', err);
        setSuggestions([]);
      } finally {
        setLoading(false);
      }
    }, 300);
    
    return () => {
      if (searchTimer.current) clearTimeout(searchTimer.current);
    };
  }, [query]);

  const handleSuggestionClick = (s) => {
    setDisplayQuery(s.display);
    setSuggestions([]);
    if (onGoto) onGoto({ lat: s.lat, lon: s.lon, zoom: 13 });
  };

  const handleGeolocate = async () => {
    try {
      const pos = await getCurrentPosition({ timeout: 8000 });
      if (onGoto) onGoto({ lat: pos.coords.latitude, lon: pos.coords.longitude, zoom: 15 });
    } catch (err) {
      console.warn('Geolocate failed', err);
      alert('Impossibile ottenere la posizione: ' + (err.message || err));
    }
  };

  const selectTab = (tabName) => {
    let drawingInstrument = null;
    switch (tabName){
      case "waypoint":
        drawingInstrument = 'marker';
        break;
      case "obstacles":
        drawingInstrument = 'polygon';
        break;
      default:
        drawingInstrument = null;
        break;
    }

    onToggleDraw(drawingInstrument);
    setTab(tabName);
  }

  return (
    <div className="d-flex flex-column h-100 sidebar">
      <ul className="nav nav-tabs">
        <li className="nav-item">
          <button className={`nav-link${tab === 'waypoint' ? ' active' : ''}`} onClick={() => selectTab('waypoint')}>Waypoints</button>
        </li>
        <li className="nav-item">
          <button className={`nav-link${tab === 'obstacles' ? ' active' : ''}`} onClick={() => selectTab('obstacles')}>Obstacles</button>
        </li>
        <li className="nav-item">
          <button className={`nav-link${tab === 'parameters' ? ' active' : ''}`} onClick={() => selectTab('parameters')}>Parameters</button>
        </li>
      </ul>

      <div className="flex-fill overflow-auto p-3" style={{ minHeight: 0 }}>
        {tab === 'waypoint' && (
          <div>
            <div className="row g-2 mb-3">
              <div className="col-md">
                <label htmlFor="altitude-value" className="form-label">Altitude</label>
                <input
                  type="number"
                  id="altitude-value"
                  className="form-control"
                  value={altitudeValue}
                  onChange={e => setAltitudeValue(Number(e.target.value))}
                />
              </div>
              <div className="col-md">
                <label htmlFor="altitude-unit" className="form-label">Unit</label>
                <select
                  id="altitude-unit"
                  className="form-select"
                  value={altitudeUnit}
                  onChange={e => setAltitudeUnit(e.target.value)}
                >
                  <option value="mt">Meters (m)</option>
                  <option value="ft">Feet (ft)</option>
                </select>
              </div>
            </div>

            <div className="mb-3">
              <label className="form-label">Cerca località</label>
              <input value={displayQuery ? displayQuery : query} onChange={e => setQuery(e.target.value)} type="text" className="form-control mb-2" placeholder="Nome o coordinate" />
              {loading && <div className="small text-muted">Ricerca...</div>}
              {suggestions.length > 0 && (
                <ul className="list-group mt-1">
                  {suggestions.map((s, i) => (
                    <li key={i} className="list-group-item list-group-item-action" style={{ cursor: 'pointer' }} onClick={() => handleSuggestionClick(s)}>
                      {s.display}
                    </li>
                  ))}
                </ul>
              )}

              <div className="d-flex gap-2 mt-2">
                {/* <button type="button" className="btn btn-primary" onClick={() => onToggleDraw && onToggleDraw('marker')}>
                  <i className="bi bi-geo-alt"></i> Marker
                </button> */}

                <button type="button" className="btn btn-outline-secondary" onClick={handleGeolocate}>
                  <i className="bi bi-crosshair"></i> Geolocalizza
                </button>
              </div>
            </div>

            <hr />

            <div className="mb-3">
              <label className="form-label">Importa waypoints da file (.geojson)</label>
              <input type="file" className="form-control" accept=".geojson,.json" />
            </div>
          </div>
        )}
        {tab === 'obstacles' && (
          <div>
            <div className="mb-3">
                <label className="form-label">Altitude</label>
                <div className="row g-2">
                    <div className="col-md">
                        <input type="number" className="form-control" placeholder="Min" value={minAltitudeValue} onChange={e => setMinAltitudeValue(Number(e.target.value))} />
                    </div>
                    <div className="col-md">
                        <input type="number" className="form-control" placeholder="Max" value={maxAltitudeValue} onChange={e => setMaxAltitudeValue(Number(e.target.value))} />
                    </div>
                    <div className="col-md">
                        <select className="form-select" value={obstacleAltitudeUnit} onChange={e => setObstacleAltitudeUnit(e.target.value)}>
                            <option value="mt">mt</option>
                            <option value="ft">ft</option>
                        </select>
                    </div>
                </div>
            </div>
            <div className="mb-3">
              {/* <button type="button" className="btn btn-outline-primary w-100 mb-2" onClick={() => onToggleDraw && onToggleDraw('polygon')}> <i className="bi bi-vector-pen"></i> Disegna poligono</button> */}
              <button type="button" className="btn btn-outline-success w-100" onClick={onGenerateRandomObstacles} disabled={!isMapReady}> <i className="bi bi-magic"></i> Genera ostacoli a caso</button>
            </div>

            <hr />

            <div className="mb-3">
              <label className="form-label">Importa ostacoli da file (.geojson)</label>
              <input type="file" className="form-control" accept=".geojson,.json" />
            </div>
          </div>
        )}
        {tab === 'parameters' && (
          <form onSubmit={(e) => e.preventDefault()}>
            <div className="mb-3">
              <label htmlFor="algorithm" className="form-label">Algorithm</label>
              <select id="algorithm" name="algorithm" className="form-select" value={parameters.algorithm} onChange={(e) => onParametersChange({...parameters, algorithm: e.target.value})}>
                <option value="antpath">Ant Path</option>
                <option value="rrt">RRT</option>
                <option value="rrtstar">RRT*</option>
              </select>
            </div>
            <div className="mb-3">
              <label htmlFor="goal_bias" className="form-label">Goal Bias</label>
              <input type="number" id="goal_bias" name="goal_bias" className="form-control" value={parameters.goal_bias} onChange={(e) => onParametersChange({...parameters, goal_bias: Number(e.target.value)})} min="0" max="1" step="0.05" />
            </div>
            <div className="mb-3">
              <label htmlFor="max_iterations" className="form-label">Max Iterations</label>
              <input type="number" id="max_iterations" name="max_iterations" className="form-control" value={parameters.max_iterations} onChange={(e) => onParametersChange({...parameters, max_iterations: Number(e.target.value)})} min="1" />
            </div>
            <div className="mb-3">
              <label htmlFor="step_size_mt" className="form-label">Step Size (meters)</label>
              <input type="number" id="step_size_mt" name="step_size_mt" className="form-control" value={parameters.step_size_mt} onChange={(e) => onParametersChange({...parameters, step_size_mt: Number(e.target.value)})} min="0" step="0.5" />
            </div>
            <div className="mb-3">
              <label htmlFor="sampler" className="form-label">Sampler</label>
              <select id="sampler" name="sampler" className="form-select" value={parameters.sampler} onChange={(e) => onParametersChange({...parameters, sampler: e.target.value})}>
                <option value="uniform">Uniform</option>
                <option value="halton">Halton</option>
              </select>
            </div>
            <div className="mb-3">
              <label htmlFor="seed" className="form-label">Seed</label>
              <input type="number" id="seed" name="seed" className="form-control" value={parameters.seed} onChange={(e) => onParametersChange({...parameters, seed: Number(e.target.value)})} />
            </div>
            <div className="mb-3">
              <label htmlFor="storage" className="form-label">Storage</label>
              <select id="storage" name="storage" className="form-select" value={parameters.storage} onChange={(e) => onParametersChange({...parameters, storage: e.target.value})}>
                <option value="list">List</option>
                <option value="rtree">R-Tree</option>
              </select>
            </div>
            <hr />
            <div className="d-grid">
                <button 
                    type="button" 
                    className="btn btn-primary" 
                    onClick={() => onCompute()} 
                    disabled={isComputing}
                >
                    {isComputing ? (
                        <>
                            <span className="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
                            <span className="ms-2">Computing...</span>
                        </>
                    ) : 'Compute'}
                </button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
}

export default Sidebar;