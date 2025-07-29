// routes
import { PATH_DASHBOARD } from './routes/paths';

// API
// ----------------------------------------------------------------------

export const HOST_API_KEY = process.env.REACT_APP_HOST_API_KEY || 'http://localhost:8081';

export const EDUQR_API = {
  baseURL: process.env.REACT_APP_EDUQR_API_URL || 'http://localhost:8081',
  endpoints: {
    auth: {
      login: '/api/v1/auth/login',
      register: '/api/v1/auth/register',
    },
    users: {
      profile: '/api/v1/users/profile',
      updateProfile: '/api/v1/users/profile',
      changePassword: '/api/v1/users/profile/password',
      validatePassword: '/api/v1/users/profile/validate-password',
      list: '/api/v1/users/all',
      create: '/api/v1/users/create',
      getById: (id) => `/api/v1/users/${id}`,
      update: (id) => `/api/v1/users/${id}`,
      delete: (id) => `/api/v1/users/${id}`,
      updateRole: (id) => `/api/v1/users/${id}/role`,
    },
    events: {
      list: '/api/v1/events',
      create: '/api/v1/events',
      getById: (id) => `/api/v1/events/${id}`,
      update: (id) => `/api/v1/events/${id}`,
      delete: (id) => `/api/v1/events/${id}`,
      range: '/api/v1/events/range',
    },
    absences: {
      list: '/api/v1/absences/my',
      create: '/api/v1/absences',
      getById: (id) => `/api/v1/absences/${id}`,
      review: (id) => `/api/v1/absences/${id}/review`,
      delete: (id) => `/api/v1/absences/${id}`,
      stats: '/api/v1/absences/stats',
      filter: '/api/v1/absences/filter',
      teacher: '/api/v1/absences/teacher',
    },
    presences: {
      my: '/api/v1/presences/my',
      scan: '/api/v1/presences/scan',
      byCourse: (courseId) => `/api/v1/presences/course/${courseId}`,
      statsByCourse: (courseId) => `/api/v1/presences/course/${courseId}/stats`,
      createForAll: (courseId) => `/api/v1/presences/course/${courseId}/create-all`,
    },
    qrCodes: {
      getInfo: (courseId) => `/api/v1/qr-codes/course/${courseId}`,
      regenerate: (courseId) => `/api/v1/qr-codes/course/${courseId}/regenerate`,
    },
    rooms: {
      list: '/api/v1/admin/rooms',
      create: '/api/v1/admin/rooms',
      getById: (id) => `/api/v1/admin/rooms/${id}`,
      update: (id) => `/api/v1/admin/rooms/${id}`,
      delete: (id) => `/api/v1/admin/rooms/${id}`,
      modular: '/api/v1/admin/rooms/modular',
    },
    subjects: {
      list: '/api/v1/admin/subjects',
      create: '/api/v1/admin/subjects',
      getById: (id) => `/api/v1/admin/subjects/${id}`,
      update: (id) => `/api/v1/admin/subjects/${id}`,
      delete: (id) => `/api/v1/admin/subjects/${id}`,
    },
    courses: {
      list: '/api/v1/admin/courses',
      create: '/api/v1/admin/courses',
      getById: (id) => `/api/v1/admin/courses/${id}`,
      update: (id) => `/api/v1/admin/courses/${id}`,
      delete: (id) => `/api/v1/admin/courses/${id}`,
      byDateRange: '/api/v1/admin/courses/by-date-range',
      byRoom: (roomId) => `/api/v1/admin/courses/by-room/${roomId}`,
      byTeacher: (teacherId) => `/api/v1/admin/courses/by-teacher/${teacherId}`,
      checkConflicts: '/api/v1/admin/courses/check-conflicts',
      checkConflictsForUpdate: (id) => `/api/v1/admin/courses/${id}/check-conflicts`,
    },
    adminAbsences: {
      list: '/api/v1/admin/absences',
    },
    adminPresences: {
      list: '/api/v1/admin/presences',
    },
    auditLogs: {
      list: '/api/v1/admin/audit-logs',
      stats: '/api/v1/admin/audit-logs/stats',
      recent: '/api/v1/admin/audit-logs/recent',
      getById: (id) => `/api/v1/admin/audit-logs/${id}`,
      userActivity: (userId) => `/api/v1/admin/audit-logs/user/${userId}/activity`,
      resourceHistory: (resourceType, resourceId) => `/api/v1/admin/audit-logs/resource/${resourceType}/${resourceId}`,
      clean: '/api/v1/admin/audit-logs/clean',
    },
  },
};

export const FIREBASE_API = {
  apiKey: process.env.REACT_APP_FIREBASE_API_KEY,
  authDomain: process.env.REACT_APP_FIREBASE_AUTH_DOMAIN,
  projectId: process.env.REACT_APP_FIREBASE_PROJECT_ID,
  storageBucket: process.env.REACT_APP_FIREBASE_STORAGE_BUCKET,
  messagingSenderId: process.env.REACT_APP_FIREBASE_MESSAGING_SENDER_ID,
  appId: process.env.REACT_APP_FIREBASE_APPID,
  measurementId: process.env.REACT_APP_FIREBASE_MEASUREMENT_ID,
};

export const COGNITO_API = {
  userPoolId: process.env.REACT_APP_AWS_COGNITO_USER_POOL_ID,
  clientId: process.env.REACT_APP_AWS_COGNITO_CLIENT_ID,
};

export const AUTH0_API = {
  clientId: process.env.REACT_APP_AUTH0_CLIENT_ID,
  domain: process.env.REACT_APP_AUTH0_DOMAIN,
};

export const MAP_API = process.env.REACT_APP_MAPBOX_API;

// ROOT PATH AFTER LOGIN SUCCESSFUL
export const PATH_AFTER_LOGIN = PATH_DASHBOARD.one;

// LAYOUT
// ----------------------------------------------------------------------

export const HEADER = {
  H_MOBILE: 64,
  H_MAIN_DESKTOP: 88,
  H_DASHBOARD_DESKTOP: 92,
  H_DASHBOARD_DESKTOP_OFFSET: 92 - 32,
};

export const NAV = {
  W_BASE: 260,
  W_LARGE: 320,
  W_DASHBOARD: 280,
  W_DASHBOARD_MINI: 88,
  //
  H_DASHBOARD_ITEM: 48,
  H_DASHBOARD_ITEM_SUB: 36,
  //
  H_DASHBOARD_ITEM_HORIZONTAL: 32,
};

export const ICON = {
  NAV_ITEM: 24,
  NAV_ITEM_HORIZONTAL: 22,
  NAV_ITEM_MINI: 22,
};
