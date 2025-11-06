import React from 'react';
import { Button, Spinner, Alert } from 'react-bootstrap';
import './ResultModal.css';

function ResultModal({ state, result, onRetry, onEdit, onClose }) {
  if (state === 'closed') {
    return null;
  }

  const renderContent = () => {
    switch (state) {
      case 'loading':
        return (
          <div className="text-center">
            <Spinner animation="border" role="status" variant="primary" style={{ width: '3rem', height: '3rem' }}>
              <span className="visually-hidden">Loading...</span>
            </Spinner>
            <p className="mt-3 mb-0 fs-5">Computing route...</p>
          </div>
        );
      case 'success':
        const processingTime = result?.completed_at && result?.received_at
        ? (new Date(result.completed_at) - new Date(result.received_at)) / 1000
        : null;

        return (
          <div>
            <Alert variant="success" className="mb-3">
              <Alert.Heading>Computation Successful!</Alert.Heading>
            </Alert>
            {result?.cost_km != null && <p><strong>Route Cost:</strong> {result.cost_km.toFixed(2)} km</p>}
            {processingTime != null && <p><strong>Processing Time:</strong> {processingTime.toFixed(3)} seconds</p>}
            {result?.parameters?.algorithm && <p><strong>Algorithm:</strong> <span className="text-uppercase">{result.parameters.algorithm}</span></p>}
            {result?.waypoints && <p className="mb-0"><strong>Waypoints:</strong> {result.waypoints.length}</p>}
          </div>
        );
      case 'error':
        return (
          <Alert variant="danger" className="mb-0">
            <Alert.Heading>Computation Failed</Alert.Heading>
            <p className="mb-0">
              {result?.message || 'An error occurred while computing the route. Please check your parameters and try again.'}
            </p>
          </Alert>
        );
      default:
        return null;
    }
  };

  const renderFooter = () => {
    switch (state) {
      case 'success':
        return (
          <>
            <Button variant="secondary" onClick={onEdit}>Edit</Button>
            <Button variant="primary" onClick={onClose}>OK</Button>
          </>
        );
      case 'error':
        return (
          <>
            <Button variant="secondary" onClick={onEdit}>Edit</Button>
            <Button variant="primary" onClick={onRetry}>Retry</Button>
          </>
        );
      default:
        return null;
    }
  };

  const isLoading = state === 'loading';

  return (
    <div className={`result-modal-overlay ${isLoading ? 'loading' : ''}`} onClick={!isLoading ? onClose : undefined}>
      <div className="result-modal-content" onClick={(e) => e.stopPropagation()}>
        {isLoading ? (
          renderContent()
        ) : (
          <>
            <div className="result-modal-header">
              <h5>
                {state === 'success' && 'Results'}
                {state === 'error' && 'Error'}
              </h5>
              <button type="button" className="btn-close" aria-label="Close" onClick={onClose}></button>
            </div>
            <div className="result-modal-body">
              {renderContent()}
            </div>
            <div className="result-modal-footer">
              {renderFooter()}
            </div>
          </>
        )}
      </div>
    </div>
  );
}

export default ResultModal;
