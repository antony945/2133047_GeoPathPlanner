import React, { useState } from 'react';
import { useAuth } from '../../context/AuthContext';

function LoginPage() {
  const [isLoginView, setIsLoginView] = useState(true);
  const { login, register } = useAuth();
  const [error, setError] = useState('');

  const [formData, setFormData] = useState({
    email: '',
    password: '',
    username: '',
    nome: '',
    cognome: '',
    country: ''
  });

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    try {
      if (isLoginView) {
        await login(formData.email, formData.password);
      } else {
        const { email, password, username, nome, cognome, country } = formData;
        await register({ email, password, username, nome, cognome, country });
      }
    } catch (err) {
      const errorMsg = err.response?.data?.message || 'An error occurred.';
      setError(Array.isArray(errorMsg) ? errorMsg.join(', ') : errorMsg);
    }
  };

  return (
    <div className="container mt-5">
      <div className="row justify-content-center">
        <div className="col-md-6">
          <div className="card">
            <div className="card-header">
              <h3 className="text-center">{isLoginView ? 'Login' : 'Register'}</h3>
            </div>
            <div className="card-body">
              <form onSubmit={handleSubmit}>
                {error && <div className="alert alert-danger">{error}</div>}

                {!isLoginView && (
                  <>
                    <div className="mb-3">
                      <label className="form-label">Username</label>
                      <input type="text" name="username" className="form-control" onChange={handleChange} required />
                    </div>
                     <div className="row">
                        <div className="col-md-6 mb-3">
                            <label className="form-label">Nome</label>
                            <input type="text" name="nome" className="form-control" onChange={handleChange} required />
                        </div>
                        <div className="col-md-6 mb-3">
                            <label className="form-label">Cognome</label>
                            <input type="text" name="cognome" className="form-control" onChange={handleChange} required />
                        </div>
                    </div>
                    <div className="mb-3">
                      <label className="form-label">Country</label>
                      <input type="text" name="country" className="form-control" onChange={handleChange} required />
                    </div>
                  </>
                )}

                <div className="mb-3">
                  <label className="form-label">Username</label>
                  <input type="text" name="email" className="form-control" onChange={handleChange} required />
                </div>

                <div className="mb-3">
                  <label className="form-label">Password</label>
                  <input type="password" name="password" className="form-control" onChange={handleChange} required />
                </div>

                <div className="d-grid">
                  <button type="submit" className="btn btn-primary">
                    {isLoginView ? 'Login' : 'Register'}
                  </button>
                </div>
              </form>
            </div>
            <div className="card-footer text-center">
              <button
                className="btn btn-link"
                onClick={() => setIsLoginView(!isLoginView)}
              >
                {isLoginView ? 'Need an account? Register' : 'Have an account? Login'}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default LoginPage;
