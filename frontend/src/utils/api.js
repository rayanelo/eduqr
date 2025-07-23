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
    return Promise.reject(error.response?.data || error);
  }
);

// Auth API
export const authAPI = {
  login: async (email, password) => {
    const response = await apiClient.post(EDUQR_API.endpoints.auth.login, {
      email,
      password,
    });
    return response.data;
  },

  register: async (email, password, firstName, lastName) => {
    const response = await apiClient.post(EDUQR_API.endpoints.auth.register, {
      email,
      password,
      first_name: firstName,
      last_name: lastName,
    });
    return response.data;
  },
};

// User API
export const userAPI = {
  getProfile: async () => {
    const response = await apiClient.get(EDUQR_API.endpoints.users.profile);
    return response.data;
  },

  updateProfile: async (userData) => {
    const response = await apiClient.put(EDUQR_API.endpoints.users.updateProfile, userData);
    return response.data;
  },
};

// Event API
export const eventAPI = {
  getEvents: async () => {
    const response = await apiClient.get(EDUQR_API.endpoints.events.list);
    return response.data;
  },

  createEvent: async (eventData) => {
    const response = await apiClient.post(EDUQR_API.endpoints.events.create, eventData);
    return response.data;
  },

  getEventById: async (id) => {
    const response = await apiClient.get(EDUQR_API.endpoints.events.getById(id));
    return response.data;
  },

  updateEvent: async (id, eventData) => {
    const response = await apiClient.put(EDUQR_API.endpoints.events.update(id), eventData);
    return response.data;
  },

  deleteEvent: async (id) => {
    const response = await apiClient.delete(EDUQR_API.endpoints.events.delete(id));
    return response.data;
  },
};

export default apiClient; 