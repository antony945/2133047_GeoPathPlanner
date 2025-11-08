import React from 'react';

const formatDate = (dateString) => {
  const options = {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  };
  // Use 'en-GB' for d/m/Y format, and let toLocaleString handle time based on system timezone
  return new Date(dateString).toLocaleString('en-GB', options).replace(',', '');
};

const HistoryCard = ({ route, onViewOnMap, onDelete }) => {
  const statusOk = route.route_found;

  // The request object to be stored for "View on Map"
  const requestForMap = {
    waypoints: route.waypoints,
    constraints: route.constraints,
    parameters: route.parameters,
    // The full route path is also needed for display
    result: {
      route: route.route,
    },
  };

  return (
    <div className="col-md-6 col-lg-4 mb-4">
      <div className="card h-100 shadow-sm">
        <div className="card-header d-flex justify-content-between align-items-center bg-light">
          <span className={`badge ${statusOk ? 'bg-success' : 'bg-danger'}`}>
            {statusOk ? 'Success' : 'Failed'}
          </span>
          <small className="text-muted">{formatDate(route.received_at)}</small>
        </div>
        <div className="card-body pb-0">
          <p className="card-text text-muted" style={{ fontSize: '0.8rem', wordBreak: 'break-all' }}>
            <strong className="text-dark">ID:</strong> {route.request_id}
          </p>
          <ul className="list-group list-group-flush">
            <li className="list-group-item d-flex justify-content-between align-items-center px-0">
              Algorithm
              <span className="badge bg-primary rounded-pill">{route.parameters.algorithm}</span>
            </li>
            <li className="list-group-item d-flex justify-content-between align-items-center px-0">
              Waypoints
              <span className="badge bg-secondary rounded-pill">{route.waypoints.length}</span>
            </li>
            <li className="list-group-item d-flex justify-content-between align-items-center px-0">
              Constraints
              <span className="badge bg-secondary rounded-pill">{route.constraints.length}</span>
            </li>
            {statusOk && (
              <li className="list-group-item d-flex justify-content-between align-items-center px-0">
                Cost (km)
                <span className="badge bg-info text-dark rounded-pill">{route.cost_km.toFixed(2)}</span>
              </li>
            )}
          </ul>
          {!statusOk && (
            <div className="alert alert-warning mt-2 p-2" role="alert" style={{ fontSize: '0.85rem' }}>
              <strong>Reason:</strong> {route.message}
            </div>
          )}
        </div>
        <div className="card-footer bg-white border-top-0 d-flex justify-content-between">
          <button
            className="btn btn-sm btn-primary"
            onClick={() => onViewOnMap(requestForMap)}
            disabled={!statusOk}
            title={!statusOk ? "No route found to display" : "View route on map"}
          >
            View on Map
          </button>
          <button className="btn btn-sm btn-outline-danger" onClick={() => onDelete(route.request_id)}>
            Delete
          </button>
        </div>
      </div>
    </div>
  );
};

export default HistoryCard;
