import axios from 'axios';
import { EDUQR_API } from '../config-global';

// Create axios instance for EduQR API
const apiClient = axios.create({
  baseURL: EDUQR_API.baseURL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('accessToken');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle errors
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Token expired or invalid
      localStorage.removeItem('accessToken');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Auth API
export const authAPI = {
  login: (credentials) => apiClient.post(EDUQR_API.endpoints.auth.login, credentials),
  register: (userData) => apiClient.post(EDUQR_API.endpoints.auth.register, userData),
};

// User API
export const userAPI = {
  getProfile: () => apiClient.get(EDUQR_API.endpoints.users.profile),
  updateProfile: (data) => apiClient.put(EDUQR_API.endpoints.users.updateProfile, data),
  getAllUsers: () => apiClient.get(EDUQR_API.endpoints.users.list),
  getUserById: (id) => apiClient.get(EDUQR_API.endpoints.users.getById(id)),
  createUser: (userData) => apiClient.post(EDUQR_API.endpoints.users.create, userData),
  updateUser: (id, userData) => apiClient.put(EDUQR_API.endpoints.users.update(id), userData),
  deleteUser: (id) => apiClient.delete(EDUQR_API.endpoints.users.delete(id)),
  updateUserRole: (id, role) => apiClient.patch(EDUQR_API.endpoints.users.updateRole(id), { role }),
};

// Event API
export const eventAPI = {
  getEvents: () => apiClient.get(EDUQR_API.endpoints.events.list),
  createEvent: (eventData) => apiClient.post(EDUQR_API.endpoints.events.create, eventData),
  getEventById: (id) => apiClient.get(EDUQR_API.endpoints.events.getById(id)),
  updateEvent: (id, eventData) => apiClient.put(EDUQR_API.endpoints.events.update(id), eventData),
  deleteEvent: (id) => apiClient.delete(EDUQR_API.endpoints.events.delete(id)),
  getEventsByDateRange: (startDate, endDate) => 
    apiClient.get(EDUQR_API.endpoints.events.range, { params: { start_date: startDate, end_date: endDate } }),
};

export { apiClient };
export default apiClient; 