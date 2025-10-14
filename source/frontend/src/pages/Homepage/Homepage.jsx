/* eslint-disable no-unused-vars */
import React, { useCallback, useRef, useState } from 'react';
import Sidebar from "../../components/MapEditor/Sidebar";
import Map from "../../components/MapEditor/Map";

function Homepage() {
  const mapRef = useRef();
  const [drawMode, setDrawMode] = useState('marker');
  const [waypoints, setWaypoints] = useState([]);
  const [obstacles, setObstacles] = useState([]);
  const [isMapReady, setIsMapReady] = useState(false);

  function handleSelectLocation({ lat, lon }) {
    mapRef.current?.goTo({ lat, lon, zoom: 14 });
  }

  const handleToggleDraw = (mode) => {
    setDrawMode(mode);
  };

  const handleFeatureCreated = (feature, type) => {
    if (type === 'marker') {
      setWaypoints(current => [...current, feature]);
    } else if (type === 'polygon') {
      setObstacles(current => [...current, feature]);
    }
  };

  /**
   * Generates a random set of polygonal obstacles on the map.
   * This function is designed to provide a quick way to populate the map with sample data.
   *
   * Process:
   * 1. It first clears any existing obstacles from both the map layer and the component's state.
   * 2. It retrieves the current geographical bounds of the visible map area.
   * 3. It then generates a random number of polygons (between 3 and 7).
   * 4. For each polygon, it generates a random number of vertices (between 3 and 6)
   *    and calculates their positions within a small radius around a random center point inside the map bounds.
   *    The vertices are sorted by angle to create a simple, non-self-intersecting polygon.
   * 5. Finally, it adds the new polygons to the map and updates the state.
   */
  const handleGenerateRandomObstacles = () => {
    if (!mapRef.current) return;

    // Clear existing obstacles from map and state
    mapRef.current.clearObstacles();

    const bounds = mapRef.current.getBounds();
    if (!bounds) return;

    const southWest = bounds.getSouthWest();
    const northEast = bounds.getNorthEast();
    const minLat = southWest.lat;
    const minLon = southWest.lng;
    const maxLat = northEast.lat;
    const maxLon = northEast.lng;

    const newObstacles = [];
    const numObstacles = Math.floor(Math.random() * 5) + 3; // 3 to 7

    for (let i = 0; i < numObstacles; i++) {
      const numVertices = Math.floor(Math.random() * 4) + 3; // 3 to 6
      const centerLon = minLon + Math.random() * (maxLon - minLon);
      const centerLat = minLat + Math.random() * (maxLat - minLat);
      
      const lonRadius = (maxLon - minLon) * 0.05 * (Math.random() * 0.5 + 0.5);
      const latRadius = (maxLat - minLat) * 0.05 * (Math.random() * 0.5 + 0.5);

      let points = [];
      for (let j = 0; j < numVertices; j++) {
        const angle = (j / numVertices) * 2 * Math.PI;
        // Add some randomness to the radius to make shapes less regular
        const lon = centerLon + Math.cos(angle) * lonRadius * (Math.random() * 0.4 + 0.8);
        const lat = centerLat + Math.sin(angle) * latRadius * (Math.random() * 0.4 + 0.8);
        points.push([lon, lat]);
      }
      points.push(points[0]); // Close the polygon

      const polygonGeoJSON = {
        type: 'Feature',
        properties: {},
        geometry: {
          type: 'Polygon',
          coordinates: [points]
        }
      };
      newObstacles.push(polygonGeoJSON);
    }

    // Add new obstacles to map and update state
    newObstacles.forEach(obstacle => {
      mapRef.current.addObstacle(obstacle);
    });
    setObstacles(newObstacles);
  };

  const handleMapReady = useCallback(() => {
    setIsMapReady(true);
  }, []);

  return (
    <div className="container-fluid px-0" style={{ overflow: "hidden" }}>
      <div className="row no-gutters" style={{ height: "calc(100vh - 64px)" }}>
        {/* Sidebar, tall, sticky, nav-tabs, content */}
        <div className="col-3 bg-light border-end p-0">
          <Sidebar
            onGoto={handleSelectLocation}
            onToggleDraw={handleToggleDraw}
            onGenerateRandomObstacles={handleGenerateRandomObstacles}
            onRequestGeolocate={() => { /* call util then goto */ }}
            isMapReady={isMapReady}
          />
        </div>
        {/* Map */}
        <div className="col-9 p-0">
          <div style={{ height: "100%", width: "100%" }}>
            <Map
              ref={mapRef}
              drawMode={drawMode}
              onFeatureCreated={handleFeatureCreated}
              onMapReady={handleMapReady}
            />
          </div>
        </div>
      </div>
    </div>
  );
}

export default Homepage;