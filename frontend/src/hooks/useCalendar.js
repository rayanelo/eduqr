import { useState, useEffect, useCallback } from 'react';
import { useSnackbar } from '../components/snackbar';
import apiClient from '../utils/api';

// ----------------------------------------------------------------------

export const useCalendar = () => {
  const [events, setEvents] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const { enqueueSnackbar } = useSnackbar();

  // Transformer les cours en événements de calendrier
  const transformCoursesToEvents = useCallback((courses) => {
    return courses.map((course) => {
      return {
        id: course.id,
        title: course.name,
        start: course.start_time,
        end: course.end_time,
        allDay: false,
        color: '#2196F3', // Bleu plus moderne pour les cours
        textColor: '#FFFFFF',
        extendedProps: {
          type: 'course',
          course: course,
          subject: course.subject?.name || 'Matière non définie',
          teacher: course.teacher ? `${course.teacher.first_name} ${course.teacher.last_name}` : 'Professeur non défini',
          room: course.room?.name || 'Salle non définie',
          description: course.description || '',
          isRecurring: course.is_recurring,
          recurrencePattern: course.recurrence_pattern,
          recurrenceEndDate: course.recurrence_end_date,
        },
      };
    });
  }, []);

  // Charger les événements (cours)
  const loadEvents = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await apiClient.get('/api/v1/admin/courses');
      const courses = response.data.data || [];
      const courseEvents = transformCoursesToEvents(courses);
      setEvents(courseEvents);
    } catch (err) {
      setError(err.message);
      enqueueSnackbar('Erreur lors du chargement des cours', { variant: 'error' });
    } finally {
      setIsLoading(false);
    }
  }, [transformCoursesToEvents, enqueueSnackbar]);

  // Créer un cours (redirige vers la création de cours)
  const createEvent = useCallback(async (newEvent) => {
    // Cette fonction redirigera vers la création de cours
    // Pour l'instant, on recharge les événements
    await loadEvents();
  }, [loadEvents]);

  // Mettre à jour un cours
  const updateEvent = useCallback(async (eventId, updatedEvent) => {
    setIsLoading(true);
    setError(null);
    
    try {
      const courseData = {
        name: updatedEvent.title,
        start_time: updatedEvent.start,
        end_time: updatedEvent.end,
        description: updatedEvent.description || '',
      };

      await apiClient.put(`/api/v1/admin/courses/${eventId}`, courseData);
      await loadEvents(); // Recharger les événements
      enqueueSnackbar('Cours mis à jour avec succès', { variant: 'success' });
    } catch (err) {
      setError(err.message);
      enqueueSnackbar('Erreur lors de la mise à jour du cours', { variant: 'error' });
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, [loadEvents, enqueueSnackbar]);

  // Supprimer un cours
  const deleteEvent = useCallback(async (eventId) => {
    setIsLoading(true);
    setError(null);
    
    try {
      await apiClient.delete(`/api/v1/admin/courses/${eventId}`);
      await loadEvents(); // Recharger les événements
      enqueueSnackbar('Cours supprimé avec succès', { variant: 'success' });
    } catch (err) {
      setError(err.message);
      enqueueSnackbar('Erreur lors de la suppression du cours', { variant: 'error' });
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, [loadEvents, enqueueSnackbar]);

  // Charger les événements au montage
  useEffect(() => {
    loadEvents();
  }, [loadEvents]);

  return {
    events,
    isLoading,
    error,
    createEvent,
    updateEvent,
    deleteEvent,
    loadEvents,
  };
}; 