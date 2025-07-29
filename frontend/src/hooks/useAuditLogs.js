import { useState, useCallback } from 'react';
import apiClient from '../utils/api';
import { EDUQR_API } from '../config-global';

export const useAuditLogs = () => {
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [pagination, setPagination] = useState({
    total: 0,
    page: 1,
    limit: 20,
    totalPages: 0,
  });

  // Récupérer les logs d'audit avec filtres et pagination
  const fetchAuditLogs = useCallback(async (filters = {}) => {
    setLoading(true);
    setError(null);

    try {
      const params = new URLSearchParams();
      
      // Pagination
      params.append('page', filters.page || 1);
      params.append('limit', filters.limit || 20);
      
      // Filtres
      if (filters.action) params.append('action', filters.action);
      if (filters.resource_type) params.append('resource_type', filters.resource_type);
      if (filters.resource_id) params.append('resource_id', filters.resource_id);
      if (filters.user_id) params.append('user_id', filters.user_id);
      if (filters.start_date) params.append('start_date', filters.start_date);
      if (filters.end_date) params.append('end_date', filters.end_date);
      if (filters.search) params.append('search', filters.search);

      const response = await apiClient.get(`${EDUQR_API.endpoints.auditLogs.list}?${params.toString()}`);
      
      setLogs(response.data.logs || []);
      setPagination({
        total: response.data.total || 0,
        page: response.data.page || 1,
        limit: response.data.limit || 20,
        totalPages: response.data.total_pages || 0,
      });
      
      return response.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la récupération des logs d\'audit');
      setLogs([]); // S'assurer que logs reste un tableau vide en cas d'erreur
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // Récupérer les statistiques des logs d'audit
  const fetchAuditStats = useCallback(async (startDate, endDate) => {
    setLoading(true);
    setError(null);

    try {
      const params = new URLSearchParams();
      if (startDate) params.append('start_date', startDate);
      if (endDate) params.append('end_date', endDate);

      const response = await apiClient.get(`${EDUQR_API.endpoints.auditLogs.stats}?${params.toString()}`);
      return response.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la récupération des statistiques');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // Récupérer les logs récents
  const fetchRecentLogs = useCallback(async (limit = 10) => {
    setLoading(true);
    setError(null);

    try {
      const response = await apiClient.get(`${EDUQR_API.endpoints.auditLogs.recent}?limit=${limit}`);
      return response.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la récupération des logs récents');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // Récupérer l'activité d'un utilisateur
  const fetchUserActivity = useCallback(async (userId, days = 30) => {
    setLoading(true);
    setError(null);

    try {
      const response = await apiClient.get(`${EDUQR_API.endpoints.auditLogs.userActivity(userId)}?days=${days}`);
      return response.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la récupération de l\'activité utilisateur');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // Récupérer l'historique d'une ressource
  const fetchResourceHistory = useCallback(async (resourceType, resourceId, limit = 20) => {
    setLoading(true);
    setError(null);

    try {
      const response = await apiClient.get(`${EDUQR_API.endpoints.auditLogs.resourceHistory(resourceType, resourceId)}?limit=${limit}`);
      return response.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la récupération de l\'historique de la ressource');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // Nettoyer les anciens logs
  const cleanOldLogs = useCallback(async (days) => {
    setLoading(true);
    setError(null);

    try {
      const response = await apiClient.delete(`${EDUQR_API.endpoints.auditLogs.clean}?days=${days}`);
      return response.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors du nettoyage des anciens logs');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // Récupérer un log spécifique par ID
  const fetchAuditLogById = useCallback(async (id) => {
    setLoading(true);
    setError(null);

    try {
      const response = await apiClient.get(`${EDUQR_API.endpoints.auditLogs.getById(id)}`);
      return response.data;
    } catch (err) {
      setError(err.response?.data?.error || 'Erreur lors de la récupération du log');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    logs,
    loading,
    error,
    pagination,
    fetchAuditLogs,
    fetchAuditStats,
    fetchRecentLogs,
    fetchUserActivity,
    fetchResourceHistory,
    cleanOldLogs,
    fetchAuditLogById,
  };
}; 