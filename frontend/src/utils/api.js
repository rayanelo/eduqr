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

// Absence API
export const absenceAPI = {
  getMyAbsences: () => apiClient.get(EDUQR_API.endpoints.absences.list),
  createAbsence: (absenceData) => apiClient.post(EDUQR_API.endpoints.absences.create, absenceData),
  getAbsenceById: (id) => apiClient.get(EDUQR_API.endpoints.absences.getById(id)),
  reviewAbsence: (id, reviewData) => apiClient.post(EDUQR_API.endpoints.absences.review(id), reviewData),
  deleteAbsence: (id) => apiClient.delete(EDUQR_API.endpoints.absences.delete(id)),
  getAbsenceStats: () => apiClient.get(EDUQR_API.endpoints.absences.stats),
  getAbsencesWithFilters: (params) => apiClient.get(EDUQR_API.endpoints.absences.filter, { params }),
  getTeacherAbsences: () => apiClient.get(EDUQR_API.endpoints.absences.teacher),
};

// Presence API
export const presenceAPI = {
  getMyPresences: () => apiClient.get(EDUQR_API.endpoints.presences.my),
  scanQRCode: (qrData) => apiClient.post(EDUQR_API.endpoints.presences.scan, qrData),
  getPresencesByCourse: (courseId) => apiClient.get(EDUQR_API.endpoints.presences.byCourse(courseId)),
  getPresenceStats: (courseId) => apiClient.get(EDUQR_API.endpoints.presences.statsByCourse(courseId)),
  createPresenceForAllStudents: (courseId) => apiClient.post(EDUQR_API.endpoints.presences.createForAll(courseId)),
};

// QR Code API
export const qrCodeAPI = {
  getQRCodeInfo: (courseId) => apiClient.get(EDUQR_API.endpoints.qrCodes.getInfo(courseId)),
  regenerateQRCode: (courseId) => apiClient.post(EDUQR_API.endpoints.qrCodes.regenerate(courseId)),
};

// Room API
export const roomAPI = {
  getAllRooms: () => apiClient.get(EDUQR_API.endpoints.rooms.list),
  createRoom: (roomData) => apiClient.post(EDUQR_API.endpoints.rooms.create, roomData),
  getRoomById: (id) => apiClient.get(EDUQR_API.endpoints.rooms.getById(id)),
  updateRoom: (id, roomData) => apiClient.put(EDUQR_API.endpoints.rooms.update(id), roomData),
  deleteRoom: (id) => apiClient.delete(EDUQR_API.endpoints.rooms.delete(id)),
  getModularRooms: () => apiClient.get(EDUQR_API.endpoints.rooms.modular),
};

// Subject API
export const subjectAPI = {
  getAllSubjects: () => apiClient.get(EDUQR_API.endpoints.subjects.list),
  createSubject: (subjectData) => apiClient.post(EDUQR_API.endpoints.subjects.create, subjectData),
  getSubjectById: (id) => apiClient.get(EDUQR_API.endpoints.subjects.getById(id)),
  updateSubject: (id, subjectData) => apiClient.put(EDUQR_API.endpoints.subjects.update(id), subjectData),
  deleteSubject: (id) => apiClient.delete(EDUQR_API.endpoints.subjects.delete(id)),
};

// Course API
export const courseAPI = {
  getAllCourses: () => apiClient.get(EDUQR_API.endpoints.courses.list),
  createCourse: (courseData) => apiClient.post(EDUQR_API.endpoints.courses.create, courseData),
  getCourseById: (id) => apiClient.get(EDUQR_API.endpoints.courses.getById(id)),
  updateCourse: (id, courseData) => apiClient.put(EDUQR_API.endpoints.courses.update(id), courseData),
  deleteCourse: (id) => apiClient.delete(EDUQR_API.endpoints.courses.delete(id)),
  getCoursesByDateRange: (startDate, endDate) => 
    apiClient.get(EDUQR_API.endpoints.courses.byDateRange, { params: { start_date: startDate, end_date: endDate } }),
  getCoursesByRoom: (roomId) => apiClient.get(EDUQR_API.endpoints.courses.byRoom(roomId)),
  getCoursesByTeacher: (teacherId) => apiClient.get(EDUQR_API.endpoints.courses.byTeacher(teacherId)),
  checkConflicts: (conflictData) => apiClient.post(EDUQR_API.endpoints.courses.checkConflicts, conflictData),
  checkConflictsForUpdate: (id, conflictData) => apiClient.post(EDUQR_API.endpoints.courses.checkConflictsForUpdate(id), conflictData),
};

// Admin Absence API
export const adminAbsenceAPI = {
  getAllAbsences: (params) => apiClient.get(EDUQR_API.endpoints.adminAbsences.list, { params }),
};

// Admin Presence API
export const adminPresenceAPI = {
  getPresencesWithFilters: (params) => apiClient.get(EDUQR_API.endpoints.adminPresences.list, { params }),
};

// Audit Log API
export const auditLogAPI = {
  getAuditLogs: (params) => apiClient.get(EDUQR_API.endpoints.auditLogs.list, { params }),
  getAuditLogStats: (params) => apiClient.get(EDUQR_API.endpoints.auditLogs.stats, { params }),
  getRecentAuditLogs: (params) => apiClient.get(EDUQR_API.endpoints.auditLogs.recent, { params }),
  getAuditLogById: (id) => apiClient.get(EDUQR_API.endpoints.auditLogs.getById(id)),
  getUserActivity: (userId) => apiClient.get(EDUQR_API.endpoints.auditLogs.userActivity(userId)),
  getResourceHistory: (resourceType, resourceId) => apiClient.get(EDUQR_API.endpoints.auditLogs.resourceHistory(resourceType, resourceId)),
  cleanOldLogs: () => apiClient.delete(EDUQR_API.endpoints.auditLogs.clean),
};

export { apiClient };
export default apiClient; 