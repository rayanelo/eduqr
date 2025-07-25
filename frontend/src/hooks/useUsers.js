import { useState, useCallback } from 'react';
import apiClient from '../utils/api';

export const useUsers = () => {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Récupérer tous les utilisateurs
  const fetchUsers = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.get('/api/v1/users/all');
      setUsers(response.data.users || []);
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la récupération des utilisateurs');
    } finally {
      setLoading(false);
    }
  }, []);

  // Récupérer un utilisateur par ID
  const getUserById = useCallback(async (id) => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.get(`/api/v1/users/${id}`);
      return response.data.user || response.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la récupération de l\'utilisateur');
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  // Créer un utilisateur
  const createUser = useCallback(async (userData) => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.post('/api/v1/users/create', userData);
      await fetchUsers(); // Rafraîchir la liste
      return response.data.user || response.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la création de l\'utilisateur');
      throw err;
    } finally {
      setLoading(false);
    }
  }, [fetchUsers]);

  // Mettre à jour un utilisateur
  const updateUser = useCallback(async (id, userData) => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.put(`/api/v1/users/${id}`, userData);
      await fetchUsers(); // Rafraîchir la liste
      return response.data.user || response.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la mise à jour de l\'utilisateur');
      throw err;
    } finally {
      setLoading(false);
    }
  }, [fetchUsers]);

  // Supprimer un utilisateur
  const deleteUser = useCallback(async (id) => {
    setLoading(true);
    setError(null);
    try {
      await apiClient.delete(`/api/v1/users/${id}`);
      await fetchUsers(); // Rafraîchir la liste
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la suppression de l\'utilisateur');
      throw err;
    } finally {
      setLoading(false);
    }
  }, [fetchUsers]);

  // Mettre à jour le rôle d'un utilisateur
  const updateUserRole = useCallback(async (id, role) => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.patch(`/api/v1/users/${id}/role`, { role });
      await fetchUsers(); // Rafraîchir la liste
      return response.data.user || response.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la mise à jour du rôle');
      throw err;
    } finally {
      setLoading(false);
    }
  }, [fetchUsers]);

  return {
    users,
    loading,
    error,
    setError,
    fetchUsers,
    getUserById,
    createUser,
    updateUser,
    deleteUser,
    updateUserRole
  };
}; 