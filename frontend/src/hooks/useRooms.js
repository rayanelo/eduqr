import { useState, useCallback } from 'react';
import apiClient from '../utils/api';

export const useRooms = () => {
  const [rooms, setRooms] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Récupérer toutes les salles
  const fetchRooms = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.get('/api/v1/admin/rooms');
      setRooms(response.data.data || []);
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la récupération des salles');
    } finally {
      setLoading(false);
    }
  }, []);

  // Récupérer une salle par ID
  const getRoomById = useCallback(async (id) => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.get(`/api/v1/admin/rooms/${id}`);
      return response.data.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la récupération de la salle');
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  // Récupérer les salles modulables
  const getModularRooms = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.get('/api/v1/admin/rooms/modular');
      return response.data.data || [];
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la récupération des salles modulables');
      return [];
    } finally {
      setLoading(false);
    }
  }, []);

  // Créer une salle
  const createRoom = useCallback(async (roomData) => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.post('/api/v1/admin/rooms', roomData);
      await fetchRooms(); // Rafraîchir la liste
      return response.data.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la création de la salle');
      throw err;
    } finally {
      setLoading(false);
    }
  }, [fetchRooms]);

  // Mettre à jour une salle
  const updateRoom = useCallback(async (id, roomData) => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.put(`/api/v1/admin/rooms/${id}`, roomData);
      await fetchRooms(); // Rafraîchir la liste
      return response.data.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la mise à jour de la salle');
      throw err;
    } finally {
      setLoading(false);
    }
  }, [fetchRooms]);

  // Supprimer une salle
  const deleteRoom = useCallback(async (id) => {
    setLoading(true);
    setError(null);
    try {
      await apiClient.delete(`/api/v1/admin/rooms/${id}`);
      await fetchRooms(); // Rafraîchir la liste
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la suppression de la salle');
      throw err;
    } finally {
      setLoading(false);
    }
  }, [fetchRooms]);

  return {
    rooms,
    loading,
    error,
    setError,
    fetchRooms,
    getRoomById,
    getModularRooms,
    createRoom,
    updateRoom,
    deleteRoom
  };
}; 