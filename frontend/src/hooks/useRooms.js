import { useState, useCallback } from 'react';
import apiClient from '../utils/api';

export const useRooms = () => {
  const [rooms, setRooms] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);

  // Récupérer toutes les salles
  const getAllRooms = useCallback(async (filters = {}) => {
    setIsLoading(true);
    setError(null);
    
    try {
      const params = new URLSearchParams();
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          params.append(key, value);
        }
      });

      const response = await apiClient.get(`/api/v1/admin/rooms?${params.toString()}`);
      setRooms(response.data.data || []);
      return response.data.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la récupération des salles');
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Récupérer une salle par ID
  const getRoomById = useCallback(async (id) => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await apiClient.get(`/api/v1/admin/rooms/${id}`);
      return response.data.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la récupération de la salle');
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Créer une nouvelle salle
  const createRoom = useCallback(async (roomData) => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await apiClient.post('/api/v1/admin/rooms', roomData);
      const newRoom = response.data.data;
      setRooms(prevRooms => [...prevRooms, newRoom]);
      return newRoom;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la création de la salle');
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Mettre à jour une salle
  const updateRoom = useCallback(async (id, roomData) => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await apiClient.put(`/api/v1/admin/rooms/${id}`, roomData);
      const updatedRoom = response.data.data;
      setRooms(prevRooms => prevRooms.map(room => room.id === id ? updatedRoom : room));
      return updatedRoom;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la mise à jour de la salle');
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Supprimer une salle
  const deleteRoom = useCallback(async (id) => {
    setIsLoading(true);
    setError(null);
    
    try {
      await apiClient.delete(`/api/v1/admin/rooms/${id}`);
      setRooms(prevRooms => prevRooms.filter(room => room.id !== id));
      return true;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la suppression de la salle');
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Récupérer les salles modulables
  const getModularRooms = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await apiClient.get('/api/v1/admin/rooms/modular');
      return response.data.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la récupération des salles modulables');
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  return {
    rooms,
    setRooms,
    isLoading,
    error,
    getAllRooms,
    getRoomById,
    createRoom,
    updateRoom,
    deleteRoom,
    getModularRooms,
  };
}; 