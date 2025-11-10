/* eslint-disable no-unused-vars */
import L from 'leaflet';
import React, { useCallback, useRef, useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import Sidebar from "../../components/MapEditor/Sidebar";
import Map from "../../components/MapEditor/Map";
import ResultModal from '../../components/ResultModal/ResultModal';
import { useAuth } from '../../context/AuthContext';

import { apiRouting } from '../../services/api';

function Homepage() {
  const [modalState, setModalState] = useState('closed');
  const [lastComputation, setLastComputation] = useState({ params: null, result: null });
  const mapRef = useRef();
  const [drawMode, setDrawMode] = useState('marker');
  const [parameters, setParameters] = useState({
    algorithm: 'rrtstar',
    goal_bias: 0.1,
    max_iterations: 10000,
    step_size_mt: 20.0,
    sampler: 'uniform',
    seed: 10,
    storage: 'rtree'
  });
  const [waypoints, setWaypoints] = useState([]);
  const [obstacles, setObstacles] = useState([]);
  const [computedRoute, setComputedRoute] = useState(null);
  const [isMapReady, setIsMapReady] = useState(false);
  const [nextId, setNextId] = useState(0);
  const [showEditControls, setShowEditControls] = useState(true);
  const [isEditingHistoryRoute, setIsEditingHistoryRoute] = useState(false);
  
  const { user } = useAuth();

  useEffect(() => {
    const routeToEditJson = sessionStorage.getItem('routeToEdit');
    if (routeToEditJson) {
      const routeToEdit = JSON.parse(routeToEditJson);
      sessionStorage.removeItem('routeToEdit');
      console.log("route to edit", routeToEdit)

      let currentId = nextId;
      const waypointsWithIds = (routeToEdit.waypoints || []).map(wp => {
        if (wp.id === undefined) {
          return { ...wp, id: currentId++ };
        }
        return wp;
      });
      const obstaclesWithIds = (routeToEdit.constraints || []).map(obs => {
        if (obs.id === undefined) {
          return { ...obs, id: currentId++ };
        }
        return obs;
      });
      setNextId(currentId);

      setWaypoints(waypointsWithIds);
      setObstacles(obstaclesWithIds);
      setComputedRoute(routeToEdit.result?.route || null);
      setShowEditControls(false);
      setIsEditingHistoryRoute(true);

      if (routeToEdit.parameters) {
        const newParams = {
            ...routeToEdit.parameters,
            step_size_mt: routeToEdit.parameters.max_step_size_mt || 20.0
        };
        delete newParams.max_step_size_mt;
        setParameters(newParams);
      }
    }
  }, []);
  const [currentAltitude, setCurrentAltitude] = useState({ value: 0, unit: 'mt' });
  const [currentObstacleAltitude, setCurrentObstacleAltitude] = useState({ min: 100, max: 500, unit: 'mt' });

  const handleObstacleAltitudeChange = useCallback((altitude) => {
    setCurrentObstacleAltitude(altitude);
  }, []);

  const handleParametersChange = useCallback((newParams) => {
    setParameters(newParams);
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
    const newId = nextId;
    setNextId(prevId => prevId + 1);

    if (type === 'marker') {
      const newWaypoint = {
        ...feature,
        id: newId,
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
        id: newId,
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

    let currentId = nextId;
    const newObstacles = [];
    const numObstacles = Math.floor(Math.random() * 2) + 3; // 3 to 5

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
        id: currentId++,
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

    setNextId(currentId);

    // Add new obstacles to map and update state
    newObstacles.forEach(obstacle => {
      const tooltip = `Altitude: ${obstacle.properties.minAltitudeValue}${obstacle.properties.altitudeUnit} - ${obstacle.properties.maxAltitudeValue}${obstacle.properties.altitudeUnit}`;
      mapRef.current.addFeature(obstacle, tooltip);
    });
    setObstacles(newObstacles);
  };

  const handleMapReady = useCallback(() => {
    setIsMapReady(true);
  }, []);

  useEffect(() => {
    if (isMapReady && mapRef.current) {
        mapRef.current.clearWaypoints();
        mapRef.current.clearObstacles();
        waypoints.forEach((wp, index) => {
            const tooltip = `Waypoint ${index + 1}: ${wp.properties.altitudeValue}${wp.properties.altitudeUnit}`;
            mapRef.current.addFeature(wp, tooltip);
        });
        obstacles.forEach(obs => {
            const tooltip = `Altitude: ${obs.properties.minAltitudeValue}${obs.properties.altitudeUnit} - ${obs.properties.maxAltitudeValue}${obs.properties.altitudeUnit}`;
            mapRef.current.addFeature(obs, tooltip);
        });
    }
  }, [waypoints, obstacles, isMapReady]);

  useEffect(() => {
    if (isMapReady && mapRef.current && computedRoute) {
        mapRef.current.drawRoute(computedRoute);
    }
  }, [computedRoute, isMapReady]);

  const handleCompute = async () => {
    setModalState('loading');
    setLastComputation(prev => ({ ...prev, params: parameters }));
    setComputedRoute(null);
    if (mapRef.current) {
        mapRef.current.clearRoute();
    }

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
            altitudeUnit: "mt"
        }
    };

    const requestPayload = {
        waypoints: waypoints,
        constraints: obstacles,
        search_volume: searchVolume,
        parameters: {
            ...parameters,
            max_step_size_mt: parameters.step_size_mt // Renaming for the backend
        }
    };
    delete requestPayload.parameters.step_size_mt;

    console.log("Request Payload:", JSON.stringify(requestPayload, null, 2));

    try {
      let url = '/routes/compute';
      if (user) {
        url = `${url}?user_id=${user.id}`
      }

      const response = await apiRouting.post(url, requestPayload);
      console.log("Response:", response.data);
      const resultData = response.data;

      setLastComputation(prev => ({ ...prev, result: resultData }));

      if (resultData?.route_found) {
          setComputedRoute(resultData.route);
          setModalState('success');
          setShowEditControls(false);
      } else {
          setComputedRoute(null);
          setModalState('error');
      }
    } catch (error) {
      console.error("Error computing route:", error);
      setLastComputation(prev => ({ ...prev, result: { message: error.message } }));
      setModalState('error');
    }
  };

  const handleRetry = () => {
    if (lastComputation.params) {
        handleCompute();
    }
  };

  const handleModalClose = () => {
    setModalState('closed');
    setShowEditControls(true);
    if (computedRoute) {
      if (mapRef.current) {
        mapRef.current.clearRoute();
      }
      setComputedRoute(null);
    }
  };

  const handleFeatureDeleted = useCallback(() => {
    console.log("handleFeature Deleted");
    console.log("feature deleted: computedRoute", computedRoute);
    setWaypoints([]);
    setObstacles([]);
    if (computedRoute) {
      if (mapRef.current) {
        mapRef.current.clearRoute();
      }
      setComputedRoute(null);
    }
  }, [computedRoute]);

  const handleFeatureEdited = useCallback((e) => {
    console.log("Handle Feature Edited");
    console.log("[handleFeatureEdited] e: ", e);
e.layers.eachLayer(layer => {
      console.log("layer", layer);
      const editedGeoJSON = layer.toGeoJSON();
      const editedId = layer.options.id;

      console.log("editedGeoJSON", editedGeoJSON);

      if (editedId === undefined) {
        console.error("Edited feature has no ID!");
        return;
      }

      if (layer instanceof L.Marker) {
        setWaypoints(current =>
          current.map(wp => (wp.id === editedId ? { ...wp, geometry: editedGeoJSON.geometry } : wp))
        );
      } else if (layer instanceof L.Polygon) {
        setObstacles(current =>
          current.map(obs => (obs.id === editedId ? { ...obs, geometry: editedGeoJSON.geometry } : obs))
        );
      }
    });
  }, []);

  const handleEnableRouteEditing = () => {
    setIsEditingHistoryRoute(false);
    if (mapRef.current) {
      mapRef.current.clearRoute();
    }
    setComputedRoute(null);
    setShowEditControls(true);
  };

  const handleFileUpload = (file, type) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      try {
        const geojson = JSON.parse(e.target.result);
        if (geojson.type === 'FeatureCollection' && Array.isArray(geojson.features)) {
          let currentId = nextId;
          const featuresWithIds = geojson.features.map(f => ({ ...f, id: currentId++ }));
          setNextId(currentId);

          const layerGroup = L.geoJSON(featuresWithIds);
          const bounds = layerGroup.getBounds();

          if (bounds.isValid() && mapRef.current) {
            mapRef.current.fitBounds(bounds);
          }

          if (type === 'waypoints') {
            setWaypoints(featuresWithIds);
          } else if (type === 'obstacles') {
            setObstacles(featuresWithIds);
          }
        } else {
          alert('Invalid GeoJSON format. Must be a FeatureCollection.');
        }
      } catch (error) {
        console.error('Error parsing GeoJSON file:', error);
        alert('Error reading or parsing the GeoJSON file.');
      }
    };
    reader.readAsText(file);
  };

  return (
    <div className="container-fluid px-0" style={{ overflowX: 'hidden' }}>
      <div className="row no-gutters" style={{ height: "calc(100vh - 64px)" }}>
        <div className="col-3 bg-light border-end p-0" style={{ position: 'relative' }}>
          <Sidebar
            onGoto={handleSelectLocation}
            onToggleDraw={handleToggleDraw}
            onGenerateRandomObstacles={handleGenerateRandomObstacles}
            onAltitudeChange={handleAltitudeChange}
            onObstacleAltitudeChange={handleObstacleAltitudeChange}
            onParametersChange={handleParametersChange}
            parameters={parameters}
            onCompute={handleCompute}
            isComputing={modalState === 'loading'}
            onRequestGeolocate={() => { }}
            isMapReady={isMapReady}
            onClearMap={() => window.location.reload()}
            isEditingHistoryRoute={isEditingHistoryRoute}
            onEnableRouteEditing={handleEnableRouteEditing}
            onFileUpload={handleFileUpload}
          />
          <ResultModal 
            state={modalState}
            result={lastComputation.result}
            onRetry={handleRetry}
            onEdit={handleModalClose}
            onClose={handleModalClose}
          />
        </div>
        {/* Map */}
        <div className="col-9 p-0">
          <div style={{ height: "100%", width: "99%" }}>
            <Map
              ref={mapRef}
              drawMode={drawMode}
              onFeatureCreated={handleFeatureCreated}
              onMapReady={handleMapReady}
              onFeatureDeleted={handleFeatureDeleted}
              onFeatureEdited={handleFeatureEdited}
              showEditControls={showEditControls}
            />
          </div>
        </div>
      </div>
    </div>
  );
}

export default Homepage;