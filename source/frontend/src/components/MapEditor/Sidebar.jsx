import { useState, useEffect, useRef } from 'react';
import './Sidebar.css';
import { getCurrentPosition , geocodeNominatim } from '../../assets/js/utils/geolocation.js';

function Sidebar({ onGoto, onToggleDraw, onGenerateRandomObstacles, isMapReady }) {
  const [tab, setTab] = useState('waypoint');
  const [query, setQuery] = useState('');
  const [displayQuery, setDisplayQuery] = useState('');
  const [suggestions, setSuggestions] = useState([]);
  const [loading, setLoading] = useState(false);
  const searchTimer = useRef(null);

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
            <div className="mb-3">
              <label className="form-label">Unità di misura altezza</label>
              <select className="form-select">
                <option value="m">Metri (m)</option>
                <option value="ft">Piedi (ft)</option>
              </select>
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
          <div>
            <div className="mb-3">
              <label htmlFor="algorithm-select" className="form-label">Algoritmo</label>
              <select id="algorithm-select" className="form-select">
                <option value="RTT">Veloce</option>
                <option value="RTTStar">Ottimale</option>
              </select>
            </div>
            <div className="mb-3">
              <label htmlFor="iterations-input" className="form-label">Numero di iterazioni</label>
              <input type="number" id="iterations-input" className="form-control" defaultValue="1000" min="100" step="100" />
            </div>
          </div>
        )}
      </div>
      {/* TODO: Add tab parametri: algoritmo, numero di iterazioni, ... */}
    </div>
  );
}

export default Sidebar;