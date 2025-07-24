import { useState, useCallback } from 'react';
import apiClient from '../utils/api';

export const useCourses = () => {
  const [courses, setCourses] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Récupérer tous les cours
  const fetchCourses = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.get('/api/v1/admin/courses');
      setCourses(response.data.data || []);
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la récupération des cours');
    } finally {
      setLoading(false);
    }
  }, []);

  // Récupérer un cours par ID
  const getCourseById = useCallback(async (id) => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.get(`/api/v1/admin/courses/${id}`);
      return response.data.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la récupération du cours');
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  // Créer un cours
  const createCourse = useCallback(async (courseData) => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.post('/api/v1/admin/courses', courseData);
      await fetchCourses(); // Rafraîchir la liste
      return response.data.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la création du cours');
      throw err;
    } finally {
      setLoading(false);
    }
  }, [fetchCourses]);

  // Mettre à jour un cours
  const updateCourse = useCallback(async (id, courseData) => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.put(`/api/v1/admin/courses/${id}`, courseData);
      await fetchCourses(); // Rafraîchir la liste
      return response.data.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la mise à jour du cours');
      throw err;
    } finally {
      setLoading(false);
    }
  }, [fetchCourses]);

  // Supprimer un cours
  const deleteCourse = useCallback(async (id) => {
    setLoading(true);
    setError(null);
    try {
      await apiClient.delete(`/api/v1/admin/courses/${id}`);
      await fetchCourses(); // Rafraîchir la liste
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la suppression du cours');
      throw err;
    } finally {
      setLoading(false);
    }
  }, [fetchCourses]);

  // Récupérer les cours par plage de dates
  const getCoursesByDateRange = useCallback(async (startDate, endDate) => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.get('/api/v1/admin/courses/by-date-range', {
        params: {
          start_date: startDate,
          end_date: endDate
        }
      });
      return response.data.data || [];
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la récupération des cours');
      return [];
    } finally {
      setLoading(false);
    }
  }, []);

  // Récupérer les cours d'une salle
  const getCoursesByRoom = useCallback(async (roomId) => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.get(`/api/v1/admin/courses/by-room/${roomId}`);
      return response.data.data || [];
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la récupération des cours');
      return [];
    } finally {
      setLoading(false);
    }
  }, []);

  // Récupérer les cours d'un enseignant
  const getCoursesByTeacher = useCallback(async (teacherId) => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.get(`/api/v1/admin/courses/by-teacher/${teacherId}`);
      return response.data.data || [];
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la récupération des cours');
      return [];
    } finally {
      setLoading(false);
    }
  }, []);

  // Vérifier les conflits
  const checkConflicts = useCallback(async (courseData) => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.post('/api/v1/admin/courses/check-conflicts', courseData);
      return response.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la vérification des conflits');
      return { has_conflicts: false, data: [] };
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    courses,
    loading,
    error,
    setError,
    fetchCourses,
    getCourseById,
    createCourse,
    updateCourse,
    deleteCourse,
    getCoursesByDateRange,
    getCoursesByRoom,
    getCoursesByTeacher,
    checkConflicts
  };
}; 