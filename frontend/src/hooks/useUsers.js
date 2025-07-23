import { useState, useEffect, useCallback } from 'react';
import { userAPI } from '../utils/api';

// ----------------------------------------------------------------------

export const useUsers = () => {
  const [users, setUsers] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const [selectedUser, setSelectedUser] = useState(null);

  // Charger tous les utilisateurs
  const loadUsers = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      console.log('Loading users...');
      const response = await userAPI.getAllUsers();
      console.log('Users response:', response);
      setUsers(response.data.users || response.data || []);
      console.log('Users set:', response.data.users || response.data || []);
    } catch (err) {
      console.error('Error loading users:', err);
      setError(err.response?.data?.message || 'Erreur lors du chargement des utilisateurs');
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Créer un utilisateur
  const createUser = useCallback(async (userData) => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await userAPI.createUser(userData);
      const newUser = response.data.user || response.data;
      setUsers(prev => [...prev, newUser]);
      return newUser;
    } catch (err) {
      setError(err.response?.data?.message || 'Erreur lors de la création de l\'utilisateur');
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Mettre à jour un utilisateur
  const updateUser = useCallback(async (id, userData) => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await userAPI.updateUser(id, userData);
      const updatedUser = response.data.user || response.data;
      setUsers(prev => prev.map(user => user.id === id ? updatedUser : user));
      return updatedUser;
    } catch (err) {
      setError(err.response?.data?.message || 'Erreur lors de la mise à jour de l\'utilisateur');
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Supprimer un utilisateur
  const deleteUser = useCallback(async (id) => {
    setIsLoading(true);
    setError(null);
    
    try {
      await userAPI.deleteUser(id);
      setUsers(prev => prev.filter(user => user.id !== id));
    } catch (err) {
      setError(err.response?.data?.message || 'Erreur lors de la suppression de l\'utilisateur');
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Mettre à jour le rôle d'un utilisateur
  const updateUserRole = useCallback(async (id, role) => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await userAPI.updateUserRole(id, role);
      const updatedUser = response.data.user || response.data;
      setUsers(prev => prev.map(user => user.id === id ? updatedUser : user));
      return updatedUser;
    } catch (err) {
      setError(err.response?.data?.message || 'Erreur lors de la mise à jour du rôle');
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Obtenir un utilisateur par ID
  const getUserById = useCallback(async (id) => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await userAPI.getUserById(id);
      const user = response.data.user || response.data;
      setSelectedUser(user);
      return user;
    } catch (err) {
      setError(err.response?.data?.message || 'Erreur lors de la récupération de l\'utilisateur');
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Charger les utilisateurs au montage
  useEffect(() => {
    loadUsers();
  }, [loadUsers]);

  return {
    users,
    selectedUser,
    isLoading,
    error,
    loadUsers,
    createUser,
    updateUser,
    deleteUser,
    updateUserRole,
    getUserById,
    setSelectedUser,
    setError,
  };
}; 