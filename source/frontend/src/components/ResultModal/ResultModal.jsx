import React from 'react';

function ResultModal({ state, result, onRetry, onEdit, onClose }) {
  if (state === 'closed') return null;

  const handleBackdropClick = (e) => {
    if (e.target === e.currentTarget) {
      onClose();
    }
  };

  return (
    <div className="modal fade show" style={{ display: 'block', backgroundColor: 'rgba(0,0,0,0.5)' }} onClick={handleBackdropClick}>
      <div className="modal-dialog modal-dialog-centered">
        <div className="modal-content">
          <div className="modal-header">
            <h5 className="modal-title">
              {state === 'loading' && 'Computing Route...'}
              {state === 'success' && 'Route Computed Successfully'}
              {state === 'error' && 'Computation Failed'}
            </h5>
            {state !== 'loading' && <button type="button" className="btn-close" onClick={onClose}></button>}
          </div>
          <div className="modal-body">
            {state === 'loading' && (
              <div className="d-flex justify-content-center align-items-center">
                <div className="spinner-border" role="status">
                  <span className="visually-hidden">Loading...</span>
                </div>
                <span className="ms-3">Please wait...</span>
              </div>
            )}
            {state === 'error' && (
              <div className="alert alert-danger mb-0">
                An unexpected error occurred. Please check the parameters and try again.
              </div>
            )}
            {state === 'success' && result && (
              <div>
                <p><strong>Path Length:</strong> {result.pathLength} segments</p>
                <p><strong>Total Waypoints:</strong> {result.waypointsCount}</p>
                <p><strong>Distance:</strong> {result.distance} km</p>
                <p><strong>Estimated Car Duration:</strong> {result.duration}</p>
              </div>
            )}
          </div>
          <div className="modal-footer">
            {state === 'error' && (
              <>
                <button type="button" className="btn btn-primary" onClick={onRetry}>Retry</button>
                <button type="button" className="btn btn-secondary" onClick={onClose}>Close</button>
              </>
            )}
            {state === 'success' && (
              <>
                <button type="button" className="btn btn-primary" onClick={onRetry}>Retry</button>
                <button type="button" className="btn btn-secondary" onClick={onEdit}>Edit</button>
                <button type="button" className="btn btn-outline-secondary" onClick={onClose}>Close</button>
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default ResultModal;
