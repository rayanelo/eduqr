import { useState, useCallback } from 'react';
import { useAuthContext } from '../auth/JwtContext';
import apiClient from '../utils/api';

export const useProfile = () => {
  const { user, updateUser } = useAuthContext();
  const [isLoading, setIsLoading] = useState(false);

  const updateProfile = useCallback(async (data) => {
    setIsLoading(true);
    try {
      const response = await apiClient.put('/users/profile', data);
      // Mettre à jour le contexte utilisateur avec les nouvelles données
      updateUser(response.data);
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.error || 'Erreur lors de la mise à jour du profil');
    } finally {
      setIsLoading(false);
    }
  }, [updateUser]);

  const changePassword = useCallback(async (data) => {
    setIsLoading(true);
    try {
      const response = await apiClient.put('/users/profile/password', data);
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.error || 'Erreur lors du changement de mot de passe');
    } finally {
      setIsLoading(false);
    }
  }, []);

  const validatePassword = useCallback(async (password) => {
    try {
      const response = await apiClient.post('/users/profile/validate-password', { password });
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.error || 'Erreur lors de la validation du mot de passe');
    }
  }, []);

  return {
    user,
    isLoading,
    updateProfile,
    changePassword,
    validatePassword,
  };
}; 