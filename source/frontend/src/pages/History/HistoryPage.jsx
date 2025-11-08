import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { apiRouting } from '../../services/api';
import { useAuth } from '../../context/AuthContext';
import HistoryCard from './HistoryCard';

function HistoryPage() {
  const [routes, setRoutes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [modalState, setModalState] = useState({ type: 'closed', data: null }); // type: 'delete'
  const navigate = useNavigate();
  const { user } = useAuth();

  const getHistory = (userId) => apiRouting.get('/routes', { params: { user_id: userId } });
  const deleteHistory = (routeId, userId) => apiRouting.delete(`/routes/${routeId}`, { params: { user_id: userId } });

  useEffect(() => {
    const fetchHistory = async () => {
      if (!user?.id) {
        setLoading(false);
        return;
      }
      try {
        setLoading(true);
        const response = await getHistory(user.id);
        setRoutes(response.data.data || []);
      } catch (err) {
        setError('Failed to fetch route history.');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchHistory();
  }, [user]);

  const handleDelete = (routeId) => {
    setModalState({ type: 'delete', data: { id: routeId } });
  };

  const confirmDelete = async () => {
    if (!modalState.data?.id || !user?.id) return;
    try {
      await deleteHistory(modalState.data.id, user.id);
      setRoutes(prev => prev.filter(r => r.request_id !== modalState.data.id));
    } catch (err) {
      console.error('Failed to delete route:', err);
      // TODO: show error modal
    } finally {
      setModalState({ type: 'closed', data: null });
    }
  };

  const handleViewOnMap = (requestData) => {
    sessionStorage.setItem('routeToEdit', JSON.stringify(requestData));
    navigate('/');
  };

  const closeModal = () => {
    setModalState({ type: 'closed', data: null });
  };

  if (loading) {
    return (
      <div className="container mt-4 text-center">
        <div className="spinner-border text-primary" role="status">
          <span className="visually-hidden">Loading...</span>
        </div>
        <p className="mt-2">Loading your route history...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mt-4">
        <div className="alert alert-danger" role="alert">
          <strong>Error:</strong> {error} Please try again later.
        </div>
      </div>
    );
  }

  return (
    <div className="container mt-4">
      <h1 className="mb-4">Route History</h1>
      <div className="row">
        {routes.map(route => (
          <HistoryCard 
            key={route.request_id}
            route={route}
            onViewOnMap={handleViewOnMap}
            onDelete={handleDelete}
          />
        ))}
        {routes.length === 0 && !loading && (
          <div className="col-12">
            <div className="alert alert-info" role="alert">
              You have no saved routes yet. Go to the <a href="/" className="alert-link">homepage</a> to create one!
            </div>
          </div>
        )}
      </div>

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
