import React, { useRef, useEffect, useState, forwardRef, useImperativeHandle } from 'react';
import { MapContainer, TileLayer, useMap, FeatureGroup } from 'react-leaflet';
import { EditControl } from 'react-leaflet-draw';
import L from 'leaflet';

import { getCurrentPosition } from '../../assets/js/utils/geolocation.js';

import './Map.css';
import 'leaflet/dist/leaflet.css';
import 'leaflet-draw/dist/leaflet.draw.css';

// FIX: Manually import and configure Leaflet's default icons
// This prevents issues with Vite/Webpack bundling where icon paths are lost.
import markerIcon2x from 'leaflet/dist/images/marker-icon-2x.png';
import markerIcon from 'leaflet/dist/images/marker-icon.png';
import markerShadow from 'leaflet/dist/images/marker-shadow.png';

delete L.Icon.Default.prototype._getIconUrl;

L.Icon.Default.mergeOptions({
  iconRetinaUrl: markerIcon2x,
  iconUrl: markerIcon,
  shadowUrl: markerShadow,
});

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

const Map = forwardRef(({ initialCenter = DEFAULT_CENTER, initialZoom = DEFAULT_ZOOM, drawMode, onFeatureCreated, onMapReady }, ref) => {
    const [center, setCenter] = useState(initialCenter);
    const [zoom, setZoom] = useState(initialZoom);
    const [geoError, setGeoError] = useState(null);
    const featureGroupRef = useRef();
    const leafletMapRef = useRef(null);

    useImperativeHandle(ref, () => ({
        goTo: ({ lat, lon, zoom }) => {
            setCenter([lat, lon]);
            setZoom(zoom);
        },
        getBounds: () => {
            return leafletMapRef.current ? leafletMapRef.current.getBounds() : null;
        },
        clearObstacles: () => {
            const featureGroup = featureGroupRef.current;
            if (!featureGroup) return;
            const layersToRemove = [];
            featureGroup.eachLayer(layer => {
                if (layer instanceof L.Polygon) {
                    layersToRemove.push(layer);
                }
            });
            layersToRemove.forEach(layer => featureGroup.removeLayer(layer));
        },
        addObstacle: (obstacleGeoJson) => {
            const featureGroup = featureGroupRef.current;
            if (!featureGroup) return;
            const newLayer = L.geoJSON(obstacleGeoJson);
            newLayer.eachLayer(layer => featureGroup.addLayer(layer));
        }
    }));

    useEffect(() => {
        let initialized = false;

        const success = (pos) => {
            if (initialized) return;
            const { latitude, longitude } = pos.coords;
            geoFeedback('', false, [latitude, longitude]);
        };

        const failure = (err) => {
            if (initialized) return;
            console.warn('Geolocation failed or denied:', err);
            geoFeedback(err.message || 'Geolocation unavailable', true);
        };

        const geoFeedback = (msg, error = false, center = DEFAULT_CENTER, zoom = DEFAULT_ZOOM) => {
            if (error) {
                setGeoError(msg);
                const geoFeedbackDiv = document.querySelector(".js-geo-feedback");
                if (geoFeedbackDiv) {
                    geoFeedbackDiv.classList.remove('d-none');
                    setTimeout(() => {
                        geoFeedbackDiv.classList.add('d-none');
                    }, 5000);
                }
            }
            setCenter(center);
            setZoom(zoom);
        };

        getCurrentPosition().then(success).catch(failure);

        return () => {
            initialized = true;
        };
    }, []);

    const _onCreated = (e) => {
        if (onFeatureCreated) {
            const { layerType, layer } = e;
            onFeatureCreated(layer.toGeoJSON(), layerType);
        }
    };

    return (
        <div style={{ height: '100%', width: '100%', minHeight: 400 }}>
            <MapContainer
                ref={leafletMapRef}
                center={center}
                zoom={zoom}
                style={{ height: '100%', width: '99%' }}
                id="geomap"
                whenReady={onMapReady}
            >
                <FeatureGroup ref={featureGroupRef}>
                    <EditControl
                        position="topright"
                        onCreated={_onCreated}
                        onEdited={() => { }}
                        onDeleted={() => { }}
                        draw={{
                            rectangle: false,
                            circle: false,
                            circlemarker: false,
                            polyline: false,
                            marker: drawMode === 'marker',
                            polygon: drawMode === 'polygon',
                        }}
                    />
                </FeatureGroup>
                <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />

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
});

export default Map;
