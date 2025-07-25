import { useState } from 'react';
import { useSnackbar } from 'notistack';
import { apiClient } from '../utils/api';

export const useDeletion = () => {
  const [isDeleting, setIsDeleting] = useState(false);
  const { enqueueSnackbar } = useSnackbar();

  // Suppression d'un utilisateur
  const deleteUser = async (userId, confirmWithCourses = false) => {
    setIsDeleting(true);
    try {
      const url = confirmWithCourses 
        ? `/api/v1/admin/users/${userId}?confirm_with_courses=true`
        : `/api/v1/admin/users/${userId}`;
      
      const response = await apiClient.delete(url);
      const result = response.data;
      
      if (result.success) {
        enqueueSnackbar(result.message, { variant: 'success' });
        if (result.warnings && result.warnings.length > 0) {
          result.warnings.forEach(warning => {
            enqueueSnackbar(warning, { variant: 'warning' });
          });
        }
        return { success: true, data: result };
      } else {
        // Si la suppression échoue à cause de cours liés, retourner les détails
        return { success: false, data: result, hasLinkedCourses: true };
      }
    } catch (error) {
      const message = error.response?.data?.error || 'Erreur lors de la suppression';
      enqueueSnackbar(message, { variant: 'error' });
      return { success: false, error };
    } finally {
      setIsDeleting(false);
    }
  };

  // Suppression d'une salle
  const deleteRoom = async (roomId) => {
    setIsDeleting(true);
    try {
      const response = await apiClient.delete(`/api/v1/admin/rooms/${roomId}`);
      const result = response.data;
      
      if (result.success) {
        enqueueSnackbar(result.message, { variant: 'success' });
        if (result.warnings && result.warnings.length > 0) {
          result.warnings.forEach(warning => {
            enqueueSnackbar(warning, { variant: 'warning' });
          });
        }
        return { success: true, data: result };
      } else {
        enqueueSnackbar(result.message, { variant: 'error' });
        return { success: false, data: result };
      }
    } catch (error) {
      const message = error.response?.data?.error || 'Erreur lors de la suppression';
      enqueueSnackbar(message, { variant: 'error' });
      return { success: false, error };
    } finally {
      setIsDeleting(false);
    }
  };

  // Suppression d'une matière
  const deleteSubject = async (subjectId) => {
    setIsDeleting(true);
    try {
      const response = await apiClient.delete(`/api/v1/admin/subjects/${subjectId}`);
      const result = response.data;
      
      if (result.success) {
        enqueueSnackbar(result.message, { variant: 'success' });
        return { success: true, data: result };
      } else {
        enqueueSnackbar(result.message, { variant: 'error' });
        return { success: false, data: result };
      }
    } catch (error) {
      const message = error.response?.data?.error || 'Erreur lors de la suppression';
      enqueueSnackbar(message, { variant: 'error' });
      return { success: false, error };
    } finally {
      setIsDeleting(false);
    }
  };

  // Suppression d'un cours
  const deleteCourse = async (courseId, deleteRecurring = false) => {
    setIsDeleting(true);
    try {
      const url = deleteRecurring 
        ? `/api/v1/admin/courses/${courseId}?delete_recurring=true`
        : `/api/v1/admin/courses/${courseId}`;
      
      const response = await apiClient.delete(url);
      const result = response.data;
      
      if (result.success) {
        enqueueSnackbar(result.message, { variant: 'success' });
        if (result.warnings && result.warnings.length > 0) {
          result.warnings.forEach(warning => {
            enqueueSnackbar(warning, { variant: 'warning' });
          });
        }
        return { success: true, data: result };
      } else {
        enqueueSnackbar(result.message, { variant: 'error' });
        return { success: false, data: result };
      }
    } catch (error) {
      const message = error.response?.data?.error || 'Erreur lors de la suppression';
      enqueueSnackbar(message, { variant: 'error' });
      return { success: false, error };
    } finally {
      setIsDeleting(false);
    }
  };

  return {
    isDeleting,
    deleteUser,
    deleteRoom,
    deleteSubject,
    deleteCourse,
  };
}; 