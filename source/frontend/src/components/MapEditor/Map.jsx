import React, { useEffect, useRef, useState } from 'react';
import { MapContainer, TileLayer, useMap } from 'react-leaflet';


import './Map.css';
import 'leaflet/dist/leaflet.css';
import 'leaflet/dist/leaflet.js';

const DEFAULT_CENTER = [41.9028, 12.4964]; // Rome lat, lon
const DEFAULT_ZOOM = 13;

function MapController({ center, zoom }) {
  const map = useMap();
  useEffect(() => {
    if (!map) return;
    if (center) {
      map.setView(center, zoom || map.getZoom(), { animate: true });
    }
  }, [map, center, zoom]);
  return null;
}


function Map() {
    const [center, setCenter] = useState(null);
    const [zoom, setZoom] = useState(DEFAULT_ZOOM);
    const [geoError, setGeoError] = useState(null);

    useEffect(() => {
        let initialized = false;

        // success callback
        const success = (pos) => {
            if (initialized) return;

            const { latitude, longitude } = pos.coords;
            geoFeedback('', false, [latitude, longitude])
        };

        // error callback
        const failure = (err) => {
            if (initialized) return;
            
            console.warn('Geolocation failed or denied:', err);
            geoFeedback(err.message || 'Geolocation unavailable', true);
        };

        const geoFeedback = (msg, error=false, center=DEFAULT_CENTER, zoom=DEFAULT_ZOOM) => {
            if (error) {
                setGeoError(msg);
                const geoFeedbackDiv = document.querySelector(".js-geo-feedback");
                geoFeedbackDiv.classList.remove('d-none');

                setTimeout(() => {
                    geoFeedbackDiv.classList.add('d-none');
                }, 5000);
            }
            setCenter(center);
            setZoom(zoom);
        }

        if ('geolocation' in navigator) {
            navigator.geolocation.getCurrentPosition(success, failure, {
                enableHighAccuracy: true,
                maximumAge: 60 * 1000, // cache position for 1 minute
            });
        } else {
            geoFeedback('Geolocation not supported', true);
        }
            return () => {
            initialized = true;
        };
    }, []);

    return (
        <div style={{ height: '100%', width: '100%', minHeight: 400 }}>
            <MapContainer
                center={center || DEFAULT_CENTER}
                zoom={zoom}
                style={{ height: '100%', width: '100%' }}
                id="geomap"
                whenCreated={mapInstance => {
                // ensure map occupies full height
                setTimeout(() => mapInstance.invalidateSize(), 100);}}
            >
                <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
                
                {/* only render controller once we have a chosen center to animate to it */}
                {center && <MapController center={center} zoom={zoom} />}
            </MapContainer>
           
        
            <div className="d-none js-geo-feedback" style={{ position: 'absolute', left: 12, bottom: 12, zIndex: 1000 }}>
                <div
                style={{
                    background: 'rgba(255,255,255,0.9)',
                    padding: 8,
                    borderRadius: 6,
                    boxShadow: '0 1px 4px rgba(0,0,0,0.2)',
                    fontSize: 12,
                    color: '#222',
                }}
                >
                {!center && <div>Determinando posizione...</div>}
                {geoError && <div style={{ color: '#c00' }}>Geo: {geoError}</div>}
                </div>
            </div>
        </div>   
    )
}

export default Map