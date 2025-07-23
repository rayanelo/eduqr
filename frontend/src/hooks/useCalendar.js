import { useState, useEffect, useCallback } from 'react';
// mock
import { _calendarEvents } from '../_mock/arrays';

// ----------------------------------------------------------------------

export const useCalendar = () => {
  const [events, setEvents] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);

  // Charger les événements
  const loadEvents = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      // Simuler un délai API
      await new Promise(resolve => setTimeout(resolve, 1000));
      setEvents(_calendarEvents);
    } catch (err) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Créer un événement
  const createEvent = useCallback(async (newEvent) => {
    setIsLoading(true);
    setError(null);
    
    try {
      // Simuler un délai API
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      const event = {
        ...newEvent,
        id: `event-${Date.now()}`,
      };
      
      setEvents(prev => [...prev, event]);
      return event;
    } catch (err) {
      setError(err.message);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Mettre à jour un événement
  const updateEvent = useCallback(async (eventId, updatedEvent) => {
    setIsLoading(true);
    setError(null);
    
    try {
      // Simuler un délai API
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      const event = {
        ...updatedEvent,
        id: eventId,
      };
      
      setEvents(prev => prev.map(e => e.id === eventId ? event : e));
      return event;
    } catch (err) {
      setError(err.message);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Supprimer un événement
  const deleteEvent = useCallback(async (eventId) => {
    setIsLoading(true);
    setError(null);
    
    try {
      // Simuler un délai API
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      setEvents(prev => prev.filter(e => e.id !== eventId));
    } catch (err) {
      setError(err.message);
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

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