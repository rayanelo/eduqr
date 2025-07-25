import { useState, useCallback } from 'react';
import { useSnackbar } from 'notistack';
import { apiClient } from '../utils/api';

export const useRooms = () => {
  const [rooms, setRooms] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const { enqueueSnackbar } = useSnackbar();

  const fetchRooms = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.get('/api/v1/admin/rooms');
      setRooms(response.data.data || []);
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors du chargement des salles');
      enqueueSnackbar('Erreur lors du chargement des salles', { variant: 'error' });
    } finally {
      setLoading(false);
    }
  }, [enqueueSnackbar]);

  const createRoom = useCallback(async (roomData) => {
    try {
      const response = await apiClient.post('/api/v1/admin/rooms', roomData);
      enqueueSnackbar('Salle créée avec succès', { variant: 'success' });
      await fetchRooms();
      return response.data;
    } catch (err) {
      const errorMessage = err.response?.data?.error || 'Erreur lors de la création de la salle';
      enqueueSnackbar(errorMessage, { variant: 'error' });
      throw new Error(errorMessage);
    }
  }, [fetchRooms, enqueueSnackbar]);

  const updateRoom = useCallback(async (id, roomData) => {
    try {
      const response = await apiClient.put(`/api/v1/admin/rooms/${id}`, roomData);
      enqueueSnackbar('Salle modifiée avec succès', { variant: 'success' });
      await fetchRooms();
      return response.data;
    } catch (err) {
      const errorMessage = err.response?.data?.error || 'Erreur lors de la modification de la salle';
      enqueueSnackbar(errorMessage, { variant: 'error' });
      throw new Error(errorMessage);
    }
  }, [fetchRooms, enqueueSnackbar]);

  const deleteRoom = useCallback(async (id) => {
    try {
      const response = await apiClient.delete(`/api/v1/admin/rooms/${id}`);
      enqueueSnackbar('Salle supprimée avec succès', { variant: 'success' });
      await fetchRooms();
      return response.data;
    } catch (err) {
      const errorMessage = err.response?.data?.error || 'Erreur lors de la suppression de la salle';
      enqueueSnackbar(errorMessage, { variant: 'error' });
      throw new Error(errorMessage);
    }
  }, [fetchRooms, enqueueSnackbar]);

  return {
    rooms,
    loading,
    error,
    fetchRooms,
    createRoom,
    updateRoom,
    deleteRoom,
  };
}; 