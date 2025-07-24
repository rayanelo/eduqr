import { useState, useCallback, useEffect } from 'react';
import apiClient from '../utils/api';

export const useTeachers = () => {
  const [teachers, setTeachers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Récupérer tous les enseignants (utilisateurs avec le rôle 'professeur')
  const fetchTeachers = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.get('/api/v1/users/all');
      const allUsers = response.data.users || [];
      // Filtrer uniquement les utilisateurs avec le rôle 'professeur'
      const teachersOnly = allUsers.filter(user => user.role === 'professeur');
      setTeachers(teachersOnly);
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la récupération des enseignants');
    } finally {
      setLoading(false);
    }
  }, []);

  // Déclencher automatiquement fetchTeachers au montage du composant
  useEffect(() => {
    fetchTeachers();
  }, [fetchTeachers]);

  return {
    teachers,
    loading,
    error,
    setError,
    fetchTeachers
  };
}; 