import { useState, useCallback } from 'react';
import apiClient from '../utils/api';

export const useSubjects = () => {
  const [subjects, setSubjects] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Récupérer toutes les matières
  const fetchSubjects = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      console.log('Fetching subjects...');
      const response = await apiClient.get('/api/v1/admin/subjects');
      console.log('Subjects response:', response.data);
      setSubjects(response.data.data || []);
    } catch (err) {
      console.error('Error fetching subjects:', err);
      console.error('Error response:', err.response);
      setError(err.response?.data?.error || 'Erreur lors de la récupération des matières');
    } finally {
      setLoading(false);
    }
  }, []);

  // Créer une nouvelle matière
  const createSubject = useCallback(async (subjectData) => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.post('/api/v1/admin/subjects', subjectData);
      setSubjects(prev => [...prev, response.data.data]);
      return response.data.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la création de la matière');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // Récupérer une matière par ID
  const getSubjectById = useCallback(async (id) => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.get(`/api/v1/admin/subjects/${id}`);
      return response.data.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la récupération de la matière');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // Mettre à jour une matière
  const updateSubject = useCallback(async (id, subjectData) => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.put(`/api/v1/admin/subjects/${id}`, subjectData);
      setSubjects(prev => prev.map(subject => 
        subject.id === id ? response.data.data : subject
      ));
      return response.data.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la mise à jour de la matière');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // Supprimer une matière
  const deleteSubject = useCallback(async (id) => {
    setLoading(true);
    setError(null);
    try {
      await apiClient.delete(`/api/v1/admin/subjects/${id}`);
      setSubjects(prev => prev.filter(subject => subject.id !== id));
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la suppression de la matière');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    subjects,
    loading,
    error,
    setError,
    fetchSubjects,
    createSubject,
    getSubjectById,
    updateSubject,
    deleteSubject,
  };
}; 