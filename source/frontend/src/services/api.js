import axios from 'axios';

const apiLogin = axios.create({
  baseURL: 'http://localhost:3000',
});

const apiRouting = axios.create({
  baseURL: 'http://localhost:8000',
});


// Add auth token to each request
apiLogin.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

export { apiLogin, apiRouting };
