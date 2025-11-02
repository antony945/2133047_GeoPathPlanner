import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import ResultModal from '../../components/ResultModal/ResultModal';

const mockRoutes = [
  {
    id: 1,
    name: 'Brussels Exploration',
    createdAt: '2025-10-28T10:30:00Z',
    request: {
      waypoints: [{ type: 'Feature', geometry: { type: 'Point', coordinates: [4.43, 50.87] }, properties: { altitudeValue: 200, altitudeUnit: 'm' } }, { type: 'Feature', geometry: { type: 'Point', coordinates: [4.46, 50.88] }, properties: { altitudeValue: 300, altitudeUnit: 'm' } }],
      constraints: [{ type: 'Feature', geometry: { type: 'Polygon', coordinates: [[ [4.45, 50.87], [4.45, 50.88], [4.44, 50.88], [4.45, 50.87] ]] }, properties: { minAltitudeValue: 100, maxAltitudeValue: 500, altitudeUnit: 'm' } }],
      parameters: { algorithm: 'rrtstar', goal_bias: 0.1, max_iterations: 10000, step_size_mt: 20, sampler: 'uniform', seed: 10, storage: 'rtree' },
    },
    result: {
      pathLength: 34,
      waypointsCount: 2,
      distance: '12.50',
      duration: '15 minutes',
    },
  },
  {
    id: 2,
    name: 'City Center Fly-by',
    createdAt: '2025-10-27T15:00:00Z',
    request: {
      waypoints: [{ type: 'Feature', geometry: { type: 'Point', coordinates: [12.49, 41.90] }, properties: { altitudeValue: 100, altitudeUnit: 'ft' } }, { type: 'Feature', geometry: { type: 'Point', coordinates: [12.51, 41.89] }, properties: { altitudeValue: 150, altitudeUnit: 'ft' } }],
      constraints: [],
      parameters: { algorithm: 'rrt', goal_bias: 0.2, max_iterations: 5000, step_size_mt: 30, sampler: 'halton', seed: 42, storage: 'list' },
    },
    result: {
      pathLength: 22,
      waypointsCount: 2,
      distance: '5.80',
      duration: '8 minutes',
    },
  },
];

function HistoryPage() {
  const [routes, setRoutes] = useState(mockRoutes);
  const [modalState, setModalState] = useState({ type: 'closed', data: null }); // type: 'view', 'delete'
  const navigate = useNavigate();

  const handleView = (route) => {
    setModalState({ type: 'view', data: route.result });
  };

  const handleDelete = (routeId) => {
    setModalState({ type: 'delete', data: { id: routeId } });
  };

  const confirmDelete = () => {
    setRoutes(prev => prev.filter(r => r.id !== modalState.data.id));
    setModalState({ type: 'closed', data: null });
  };

  const handleEdit = (route) => {
    navigate('/', { state: { routeToEdit: route.request } });
  };

  const closeModal = () => {
    setModalState({ type: 'closed', data: null });
  };

  return (
    <div className="container mt-4">
      <h1 className="mb-4">Route History</h1>
      <div className="row">
        {routes.map(route => (
          <div key={route.id} className="col-md-6 col-lg-4 mb-4">
            <div className="card h-100">
              <div className="card-body">
                <h5 className="card-title">{route.name}</h5>
                <p className="card-text">
                  <small className="text-muted">Created: {new Date(route.createdAt).toLocaleString()}</small>
                </p>
                <p>
                  {route.request.waypoints.length} waypoints, {route.request.constraints.length} obstacles.
                </p>
              </div>
              <div className="card-footer d-flex justify-content-between">
                <button className="btn btn-sm btn-outline-primary" onClick={() => handleView(route)}>View</button>
                <button className="btn btn-sm btn-outline-secondary" onClick={() => handleEdit(route)}>Edit</button>
                <button className="btn btn-sm btn-outline-danger" onClick={() => handleDelete(route.id)}>Delete</button>
              </div>
            </div>
          </div>
        ))}
        {routes.length === 0 && <p>No saved routes found.</p>}
      </div>

      {/* View Modal */}
      {modalState.type === 'view' && (
        <ResultModal
          state="success"
          result={modalState.data}
          onClose={closeModal}
          onEdit={closeModal} 
          onRetry={() => {}}
        />
      )}

      {/* Delete Confirmation Modal */}
      {modalState.type === 'delete' && (
        <div className="modal fade show" style={{ display: 'block', backgroundColor: 'rgba(0,0,0,0.5)' }}>
          <div className="modal-dialog modal-dialog-centered">
            <div className="modal-content">
              <div className="modal-header">
                <h5 className="modal-title">Confirm Deletion</h5>
                <button type="button" className="btn-close" onClick={closeModal}></button>
              </div>
              <div className="modal-body">
                <p>Are you sure you want to delete this route? This action cannot be undone.</p>
              </div>
              <div className="modal-footer">
                <button type="button" className="btn btn-secondary" onClick={closeModal}>Cancel</button>
                <button type="button" className="btn btn-danger" onClick={confirmDelete}>Delete</button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default HistoryPage;
