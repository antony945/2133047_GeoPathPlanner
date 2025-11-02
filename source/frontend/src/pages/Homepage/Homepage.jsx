/* eslint-disable no-unused-vars */
import React, { useCallback, useRef, useState } from 'react';
import Sidebar from "../../components/MapEditor/Sidebar";
import Map from "../../components/MapEditor/Map";
import ResultModal from '../../components/ResultModal/ResultModal';

function Homepage() {
  const [modalState, setModalState] = useState('closed'); // closed, loading, success, error
  const [lastComputation, setLastComputation] = useState({ params: null, result: null });
  const mapRef = useRef();
  const [drawMode, setDrawMode] = useState('marker');
  const [waypoints, setWaypoints] = useState([]);
  const [obstacles, setObstacles] = useState([]);
  const [isMapReady, setIsMapReady] = useState(false);
  const [currentAltitude, setCurrentAltitude] = useState({ value: 0, unit: 'm' });
  const [currentObstacleAltitude, setCurrentObstacleAltitude] = useState({ min: 100, max: 500, unit: 'm' });

  const handleObstacleAltitudeChange = useCallback((altitude) => {
    setCurrentObstacleAltitude(altitude);
  }, []);

  const handleAltitudeChange = useCallback((altitude) => {
    setCurrentAltitude(altitude);
  }, []);

  function handleSelectLocation({ lat, lon }) {
    mapRef.current?.goTo({ lat, lon, zoom: 14 });
  }

  const handleToggleDraw = (mode) => {
    setDrawMode(mode);
  };

  const handleFeatureCreated = (feature, type) => {
    if (type === 'marker') {
      const newWaypoint = {
        ...feature,
        properties: {
          ...feature.properties,
          altitudeValue: currentAltitude.value,
          altitudeUnit: currentAltitude.unit,
        }
      };
      setWaypoints(current => [...current, newWaypoint]);
    } else if (type === 'polygon') {
      const newObstacle = {
        ...feature,
        properties: {
          ...feature.properties,
          minAltitudeValue: currentObstacleAltitude.min,
          maxAltitudeValue: currentObstacleAltitude.max,
          altitudeUnit: currentObstacleAltitude.unit,
        }
      };
      setObstacles(current => [...current, newObstacle]);
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
        properties: {
          minAltitudeValue: currentObstacleAltitude.min,
          maxAltitudeValue: currentObstacleAltitude.max,
          altitudeUnit: currentObstacleAltitude.unit,
        },
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

  const handleCompute = (params) => {
    setModalState('loading');
    setLastComputation(prev => ({ ...prev, params }));

    const bounds = mapRef.current.getBounds();
    if (!bounds) {
        alert("Could not get map bounds.");
        setModalState('closed');
        return;
    }
    const southWest = bounds.getSouthWest();
    const northEast = bounds.getNorthEast();

    const searchVolume = {
        type: "Feature",
        geometry: {
            type: "Polygon",
            coordinates: [[
                [southWest.lng, northEast.lat],
                [northEast.lng, northEast.lat],
                [northEast.lng, southWest.lat],
                [southWest.lng, southWest.lat],
                [southWest.lng, northEast.lat]
            ]]
        },
        properties: {
            minAltitudeValue: -999999,
            maxAltitudeValue: 999999,
            altitudeUnit: "m"
        }
    };

    const requestPayload = {
        waypoints: waypoints,
        constraints: obstacles,
        search_volume: searchVolume,
        parameters: {
            ...params,
            max_step_size_mt: params.step_size_mt // Renaming for the backend
        }
    };
    delete requestPayload.parameters.step_size_mt;

    console.log("Request Payload:", JSON.stringify(requestPayload, null, 2));

    // Simulate API call
    setTimeout(() => {
        const isSuccess = Math.random() > 0.3; // 70% chance of success
        if (isSuccess) {
            const fakeResult = {
                pathLength: Math.floor(Math.random() * 50) + 10,
                waypointsCount: waypoints.length,
                distance: (Math.random() * 100).toFixed(2),
                duration: `${Math.floor(Math.random() * 60) + 5} minutes`
            };
            setLastComputation(prev => ({ ...prev, result: fakeResult }));
            setModalState('success');
        } else {
            setModalState('error');
        }
    }, 2000);
  };

  const handleRetry = () => {
    if (lastComputation.params) {
        handleCompute(lastComputation.params);
    }
  };

  const handleModalClose = () => {
    setModalState('closed');
  };

  return (
    <div className="container-fluid px-0" style={{ overflow: "hidden" }}>
      <div className="row no-gutters" style={{ height: "calc(100vh - 64px)" }}>
        {/* Sidebar, tall, sticky, nav-tabs, content */}
        <div className="col-3 bg-light border-end p-0">
          <Sidebar
            onGoto={handleSelectLocation}
            onToggleDraw={handleToggleDraw}
            onGenerateRandomObstacles={handleGenerateRandomObstacles}
            onAltitudeChange={handleAltitudeChange}
            onObstacleAltitudeChange={handleObstacleAltitudeChange}
            onCompute={handleCompute}
            isComputing={modalState === 'loading'}
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
      <ResultModal 
        state={modalState}
        result={lastComputation.result}
        onRetry={handleRetry}
        onEdit={handleModalClose}
        onClose={handleModalClose}
      />
    </div>
  );
}

export default Homepage;