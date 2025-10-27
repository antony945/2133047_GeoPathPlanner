import React, { useState, useEffect } from 'react';
import api from '../../services/api';

function ProfilePage() {
  const [loading, setLoading] = useState(true);
  const [profileError, setProfileError] = useState('');
  const [profileSuccess, setProfileSuccess] = useState('');

  // state for user profile data
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    nome: '',
    cognome: '',
    country: ''
  });  
  // State for password change
  const [passwordData, setPasswordData] = useState({ oldPassword: '', newPassword: '' });
  const [passwordError, setPasswordError] = useState('');
  const [passwordSuccess, setPasswordSuccess] = useState('');

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const response = await api.get('/users/me');
        setFormData(response.data);
      } catch (err) {
        console.error(err)
        setProfileError('Failed to fetch user data. Please try again later.');
      } finally {
        setLoading(false);
      }
    };
    fetchUserData();
  }, []);

  const handleProfileChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handlePasswordChange = (e) => {
    setPasswordData({ ...passwordData, [e.target.name]: e.target.value });
  };

  const handleProfileSubmit = async (e) => {
    e.preventDefault();
    setProfileError('');
    setProfileSuccess('');
    
    const updateData = {
      nome: formData.nome,
      cognome: formData.cognome,
      country: formData.country,
      email: formData.email
    };

    try {
      await api.put('/users/me', updateData);
      setProfileSuccess('Profile updated successfully!');
    } catch (err) {
      const errorMsg = err.response?.data?.message || 'An error occurred while updating.';
      setProfileError(Array.isArray(errorMsg) ? errorMsg.join(', ') : errorMsg);
    }
  };

  const handlePasswordSubmit = async (e) => {
    e.preventDefault();
    setPasswordError('');
    setPasswordSuccess('');
    if (!passwordData.oldPassword || !passwordData.newPassword) {
      setPasswordError('Please fill in both password fields.');
      return;
    }
    try {
      await api.patch('/users/me/password', passwordData);
      setPasswordSuccess('Password changed successfully!');
      setPasswordData({ oldPassword: '', newPassword: '' }); 
    } catch (err) {
      const errorMsg = err.response?.data?.message || 'An error occurred.';
      setPasswordError(Array.isArray(errorMsg) ? errorMsg.join(', ') : errorMsg);
    }
  };

  if (loading) {
    return <div className="container mt-4">Loading profile...</div>;
  }

  return (
    <div className="container mt-5 mb-5">
      <div className="row justify-content-center">
        <div className="col-md-8">
          {/* Profile Card */}
          <div className="card">
            <div className="card-header">
              <h3 className="text-center">Your Profile</h3>
            </div>
            <div className="card-body">
              <form onSubmit={handleProfileSubmit}>
                {profileError && <div className="alert alert-danger">{profileError}</div>}
                {profileSuccess && <div className="alert alert-success">{profileSuccess}</div>}
                
                <div className="mb-3">
                  <label className="form-label">Username</label>
                  <input type="text" name="username" className="form-control" value={formData.username} disabled readOnly />
                  <div className="form-text">Username cannot be changed.</div>
                </div>

                <div className="mb-3">
                  <label className="form-label">Email</label>
                  <input type="email" name="email" className="form-control" value={formData.email} onChange={handleProfileChange} required />
                </div>

                <div className="row">
                  <div className="col-md-6 mb-3">
                    <label className="form-label">Nome</label>
                    <input type="text" name="nome" className="form-control" value={formData.nome} onChange={handleProfileChange} required />
                  </div>
                  <div className="col-md-6 mb-3">
                    <label className="form-label">Cognome</label>
                    <input type="text" name="cognome" className="form-control" value={formData.cognome} onChange={handleProfileChange} required />
                  </div>
                </div>
                <div className="mb-3">
                  <label className="form-label">Country</label>
                  <input type="text" name="country" className="form-control" value={formData.country} onChange={handleProfileChange} required />
                </div>
                <div className="d-grid">
                  <button type="submit" className="btn btn-primary">Save Profile Changes</button>
                </div>
              </form>
            </div>
          </div>

          {/* Change Password Card */}
          <div className="card mt-4">
            <div className="card-header">
              <h4 className="text-center">Change Password</h4>
            </div>
            <div className="card-body">
              <form onSubmit={handlePasswordSubmit}>
                {passwordError && <div className="alert alert-danger">{passwordError}</div>}
                {passwordSuccess && <div className="alert alert-success">{passwordSuccess}</div>}
                <div className="mb-3">
                  <label className="form-label">Current Password</label>
                  <input type="password" name="oldPassword" className="form-control" value={passwordData.oldPassword} onChange={handlePasswordChange} required />
                </div>
                <div className="mb-3">
                  <label className="form-label">New Password</label>
                  <input type="password" name="newPassword" className="form-control" value={passwordData.newPassword} onChange={handlePasswordChange} required />
                </div>
                <div className="d-grid">
                  <button type="submit" className="btn btn-secondary">Change Password</button>
                </div>
              </form>
            </div>
          </div>

        </div>
      </div>
    </div>
  );
}

export default ProfilePage;
