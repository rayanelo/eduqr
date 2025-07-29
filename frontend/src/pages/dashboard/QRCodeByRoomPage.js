import { useState, useEffect, useCallback } from 'react';
import {
  Box,
  Card,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Typography,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Grid,
  Chip,
  Alert,
  CircularProgress,
} from '@mui/material';
import { useSnackbar } from 'notistack';

import { apiClient } from '../../utils/api';
import { QRCodeSVG } from 'qrcode.react';
import Iconify from '../../components/iconify';

export default function QRCodeByRoomPage() {
  const { enqueueSnackbar } = useSnackbar();
  
  const [rooms, setRooms] = useState([]);
  const [roomsWithCourses, setRoomsWithCourses] = useState([]);
  const [selectedRoom, setSelectedRoom] = useState('');
  const [courses, setCourses] = useState([]);
  const [loading, setLoading] = useState(false);
  const [loadingRooms, setLoadingRooms] = useState(false);
  const [qrModalOpen, setQrModalOpen] = useState(false);
  const [selectedCourse, setSelectedCourse] = useState(null);
  const [qrCodeData, setQrCodeData] = useState(null);
  const [autoUpdate, setAutoUpdate] = useState(false);

  // Load rooms with future courses
  const loadRoomsWithCourses = useCallback(async () => {
    setLoadingRooms(true);
    try {
      const response = await apiClient.get('/api/v1/admin/rooms');
      const allRooms = response.data.data || [];
      setRooms(allRooms);

      // Charger les cours futurs pour chaque salle
      const today = new Date().toISOString().split('T')[0];
      const now = new Date();
      // Permettre les cours qui commencent dans les 2 heures précédentes (pour les cours en cours)
      const twoHoursAgo = new Date(now.getTime() - 2 * 60 * 60 * 1000);
      const roomsWithFutureCourses = [];

      for (const room of allRooms) {
        try {
          const coursesResponse = await apiClient.get(`/api/v1/admin/courses/by-room/${room.id}?date=${today}`);
          
          const futureCourses = (coursesResponse.data.data || [])
            .filter(course => {
              const courseStart = new Date(course.start_time);
              const courseEnd = new Date(course.end_time);
              // Inclure les cours qui commencent dans les 2 heures précédentes ou qui sont en cours
              return courseStart >= twoHoursAgo && courseEnd >= now;
            })
            .sort((a, b) => new Date(a.start_time) - new Date(b.start_time));

          if (futureCourses.length > 0) {
            roomsWithFutureCourses.push({
              ...room,
              futureCourses,
              nextCourse: futureCourses[0] // Le prochain cours
            });
          }
        } catch (error) {
          console.error(`Error loading courses for room ${room.id}:`, error);
        }
      }

      // Trier les salles par heure du prochain cours
      roomsWithFutureCourses.sort((a, b) => 
        new Date(a.nextCourse.start_time) - new Date(b.nextCourse.start_time)
      );

      setRoomsWithCourses(roomsWithFutureCourses);
    } catch (error) {
      console.error('Error loading rooms:', error);
      enqueueSnackbar('Erreur lors du chargement des salles', { variant: 'error' });
    } finally {
      setLoadingRooms(false);
    }
  }, [enqueueSnackbar]);

  // Load courses for selected room
  const loadCourses = useCallback(async () => {
    if (!selectedRoom) {
      setCourses([]);
      return;
    }

    setLoading(true);
    try {
      // Charger tous les cours futurs (à partir d'aujourd'hui)
      const today = new Date().toISOString().split('T')[0];
      const response = await apiClient.get(`/api/v1/admin/courses/by-room/${selectedRoom}?date=${today}`);
      
      // Filtrer et trier les cours futurs
      const now = new Date();
      // Permettre les cours qui commencent dans les 2 heures précédentes (pour les cours en cours)
      const twoHoursAgo = new Date(now.getTime() - 2 * 60 * 60 * 1000);
      
      const futureCourses = (response.data.data || [])
        .filter(course => {
          const courseStart = new Date(course.start_time);
          const courseEnd = new Date(course.end_time);
          // Inclure les cours qui commencent dans les 2 heures précédentes ou qui sont en cours
          return courseStart >= twoHoursAgo && courseEnd >= now;
        })
        .sort((a, b) => new Date(a.start_time) - new Date(b.start_time));
      
      setCourses(futureCourses);
    } catch (error) {
      console.error('Error loading courses:', error);
      enqueueSnackbar('Erreur lors du chargement des cours', { variant: 'error' });
    } finally {
      setLoading(false);
    }
  }, [selectedRoom, enqueueSnackbar]);

  // Load QR code for a course
  const loadQRCode = useCallback(async (courseId) => {
    try {
      const response = await apiClient.get(`/api/v1/qr-codes/course/${courseId}`);
      setQrCodeData(response.data.qr_code);
    } catch (error) {
      console.error('Error loading QR code:', error);
      enqueueSnackbar('Erreur lors du chargement du QR code', { variant: 'error' });
    }
  }, [enqueueSnackbar]);

  // Open QR modal for a course
  const openQRModal = useCallback(async (course) => {
    setSelectedCourse(course);
    setQrModalOpen(true);
    await loadQRCode(course.id);
  }, [loadQRCode]);

  // Auto-update QR code for next course
  useEffect(() => {
    if (!autoUpdate || !qrModalOpen || !selectedCourse) return;

    const interval = setInterval(async () => {
      const now = new Date();
      const currentCourse = courses.find(course => {
        const startTime = new Date(course.start_time);
        const endTime = new Date(course.end_time);
        return now >= startTime && now <= endTime;
      });

      if (currentCourse && currentCourse.id !== selectedCourse.id) {
        setSelectedCourse(currentCourse);
        await loadQRCode(currentCourse.id);
        enqueueSnackbar(`QR code mis à jour pour: ${currentCourse.name}`, { variant: 'info' });
      }
    }, 30000); // Check every 30 seconds

    return () => clearInterval(interval);
  }, [autoUpdate, qrModalOpen, selectedCourse, courses, loadQRCode, enqueueSnackbar]);

  // Load rooms on component mount
  useEffect(() => {
    loadRoomsWithCourses();
  }, [loadRoomsWithCourses]);

  // Load courses when room changes
  useEffect(() => {
    loadCourses();
  }, [loadCourses]);

  // Format time
  const formatTime = (dateString) => {
    return new Date(dateString).toLocaleTimeString('fr-FR', {
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  // Format date and time
  const formatDateTime = (dateString) => {
    const date = new Date(dateString);
    const today = new Date();
    const tomorrow = new Date(today);
    tomorrow.setDate(tomorrow.getDate() + 1);
    
    if (date.toDateString() === today.toDateString()) {
      return `Aujourd'hui à ${date.toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' })}`;
    } else if (date.toDateString() === tomorrow.toDateString()) {
      return `Demain à ${date.toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' })}`;
    } else {
      return date.toLocaleDateString('fr-FR', { 
        weekday: 'long', 
        day: 'numeric', 
        month: 'long' 
      }) + ` à ${date.toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' })}`;
    }
  };

  // Get current course status
  const getCourseStatus = (course) => {
    const now = new Date();
    const startTime = new Date(course.start_time);
    const endTime = new Date(course.end_time);

    if (now < startTime) {
      return { status: 'upcoming', label: 'À venir', color: 'default' };
    } else if (now >= startTime && now <= endTime) {
      return { status: 'current', label: 'En cours', color: 'success' };
    } else {
      return { status: 'finished', label: 'Terminé', color: 'error' };
    }
  };

  // Get time until next course
  const getTimeUntilCourse = (course) => {
    const now = new Date();
    const startTime = new Date(course.start_time);
    const diffMs = startTime - now;
    
    if (diffMs <= 0) {
      return 'En cours';
    }
    
    const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
    const diffMinutes = Math.floor((diffMs % (1000 * 60 * 60)) / (1000 * 60));
    
    if (diffHours > 0) {
      return `Dans ${diffHours}h${diffMinutes > 0 ? ` ${diffMinutes}min` : ''}`;
    } else {
      return `Dans ${diffMinutes}min`;
    }
  };

  return (
    <Box>
      <Typography variant="h4" sx={{ mb: 3 }}>
        Affichage QR Codes par Salle
      </Typography>

      {/* Rooms Overview */}
      <Card sx={{ p: 3, mb: 3 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Typography variant="h6">
            Salles avec cours futurs
          </Typography>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Button
              variant="outlined"
              size="small"
              onClick={loadRoomsWithCourses}
              disabled={loadingRooms}
              startIcon={<Iconify icon="mdi:refresh" />}
            >
              Actualiser
            </Button>
            {loadingRooms && <CircularProgress size={20} />}
          </Box>
        </Box>
        
        {loadingRooms ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        ) : roomsWithCourses.length === 0 ? (
          <Alert severity="info">
            Aucune salle n'a de cours futurs programmés.
          </Alert>
        ) : (
          <Grid container spacing={2}>
            {roomsWithCourses.map((room) => {
              const nextCourseStatus = getCourseStatus(room.nextCourse);
              return (
                <Grid item xs={12} sm={6} md={4} key={room.id}>
                  <Card
                    sx={{
                      p: 2,
                      cursor: 'pointer',
                      border: selectedRoom === room.id ? 2 : 1,
                      borderColor: selectedRoom === room.id ? 'primary.main' : 'divider',
                      '&:hover': {
                        borderColor: 'primary.main',
                        boxShadow: 2,
                      },
                    }}
                    onClick={() => setSelectedRoom(room.id)}
                  >
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 1 }}>
                      <Typography variant="h6" sx={{ fontWeight: 'bold' }}>
                        {room.name}
                      </Typography>
                      <Chip
                        label={`${room.futureCourses.length} cours`}
                        size="small"
                        color="primary"
                        variant="outlined"
                      />
                    </Box>
                    
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                      {room.capacity} places
                    </Typography>
                    
                    <Box sx={{ mb: 1 }}>
                      <Typography variant="subtitle2" color="text.secondary">
                        Prochain cours :
                      </Typography>
                      <Typography variant="body2" fontWeight="medium">
                        {room.nextCourse.subject.name}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        {formatDateTime(room.nextCourse.start_time)}
                      </Typography>
                      <Typography variant="caption" color="primary.main" fontWeight="medium">
                        {getTimeUntilCourse(room.nextCourse)}
                      </Typography>
                    </Box>
                    
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                      <Chip
                        label={nextCourseStatus.label}
                        color={nextCourseStatus.color}
                        size="small"
                        variant={nextCourseStatus.status === 'current' ? 'filled' : 'outlined'}
                      />
                      <Button
                        variant="outlined"
                        size="small"
                        onClick={(e) => {
                          e.stopPropagation();
                          setSelectedRoom(room.id);
                        }}
                        startIcon={<Iconify icon="mdi:qrcode" />}
                      >
                        Voir cours
                      </Button>
                    </Box>
                  </Card>
                </Grid>
              );
            })}
          </Grid>
        )}
      </Card>

      {/* Traditional Room Selection (fallback) */}
      <Card sx={{ p: 3, mb: 3 }}>
        <Typography variant="h6" sx={{ mb: 2 }}>
          Ou sélectionner une salle spécifique
        </Typography>
        <FormControl fullWidth>
          <InputLabel>Sélectionner une salle</InputLabel>
          <Select
            value={selectedRoom}
            onChange={(e) => setSelectedRoom(e.target.value)}
            label="Sélectionner une salle"
          >
            <MenuItem value="">
              <em>Sélectionner une salle</em>
            </MenuItem>
            {rooms.map((room) => (
              <MenuItem key={room.id} value={room.id}>
                {room.name} - {room.capacity} places
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </Card>

      {/* Courses Table */}
      {selectedRoom && (
        <Card>
          <Box sx={{ p: 3, borderBottom: 1, borderColor: 'divider' }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <Box>
                <Typography variant="h6">
                  Prochains cours - {rooms.find(r => r.id === selectedRoom)?.name}
                </Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                  Cours programmés dans l'ordre chronologique
                </Typography>
              </Box>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <Button
                  variant="outlined"
                  size="small"
                  onClick={loadCourses}
                  disabled={loading}
                  startIcon={<Iconify icon="mdi:refresh" />}
                >
                  Actualiser
                </Button>
                <Chip
                  label={`${courses.length} cours`}
                  color="primary"
                  variant="outlined"
                  size="small"
                />
              </Box>
            </Box>
          </Box>

          {loading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
              <CircularProgress />
            </Box>
          ) : courses.length === 0 ? (
            <Box sx={{ p: 3 }}>
              <Alert severity="info">
                Aucun cours futur programmé dans cette salle.
              </Alert>
            </Box>
          ) : (
            <TableContainer component={Paper}>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Date et heure</TableCell>
                    <TableCell>Durée</TableCell>
                    <TableCell>Matière</TableCell>
                    <TableCell>Professeur</TableCell>
                    <TableCell>Statut</TableCell>
                    <TableCell>Actions</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {courses.map((course) => {
                    const status = getCourseStatus(course);
                    const startTime = new Date(course.start_time);
                    const endTime = new Date(course.end_time);
                    const duration = Math.round((endTime - startTime) / (1000 * 60)); // Durée en minutes
                    
                    return (
                      <TableRow 
                        key={course.id}
                        sx={{
                          backgroundColor: status.status === 'current' ? 'action.hover' : 'inherit',
                          '&:hover': {
                            backgroundColor: 'action.hover',
                          },
                        }}
                      >
                        <TableCell>
                          <Box>
                            <Typography variant="body2" fontWeight="medium">
                              {formatDateTime(course.start_time)}
                            </Typography>
                            <Typography variant="caption" color="text.secondary">
                              {formatTime(course.start_time)} - {formatTime(course.end_time)}
                            </Typography>
                          </Box>
                        </TableCell>
                        <TableCell>
                          <Typography variant="body2">
                            {duration} min
                          </Typography>
                        </TableCell>
                        <TableCell>
                          <Typography variant="body2" fontWeight="medium">
                            {course.subject.name}
                          </Typography>
                        </TableCell>
                        <TableCell>
                          <Typography variant="body2">
                            {course.teacher.first_name} {course.teacher.last_name}
                          </Typography>
                        </TableCell>
                        <TableCell>
                          <Chip
                            label={status.label}
                            color={status.color}
                            size="small"
                            variant={status.status === 'current' ? 'filled' : 'outlined'}
                          />
                        </TableCell>
                        <TableCell>
                          <Button
                            variant="contained"
                            size="small"
                            onClick={() => openQRModal(course)}
                            disabled={status.status === 'finished'}
                            startIcon={<Iconify icon="mdi:qrcode" />}
                          >
                            Voir QR
                          </Button>
                        </TableCell>
                      </TableRow>
                    );
                  })}
                </TableBody>
              </Table>
            </TableContainer>
          )}
        </Card>
      )}

      {/* QR Code Modal */}
      <Dialog
        open={qrModalOpen}
        onClose={() => setQrModalOpen(false)}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>
          <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
            <Box>
              <Typography variant="h6">
                QR Code - {selectedCourse?.subject?.name}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {selectedCourse?.name}
              </Typography>
            </Box>
            <Button
              onClick={() => setAutoUpdate(!autoUpdate)}
              variant={autoUpdate ? "contained" : "outlined"}
              color={autoUpdate ? "success" : "primary"}
              size="small"
              startIcon={<Iconify icon={autoUpdate ? "mdi:refresh" : "mdi:refresh-off"} />}
            >
              {autoUpdate ? "Auto-update activé" : "Activer auto-update"}
            </Button>
          </Box>
        </DialogTitle>
        <DialogContent>
          {selectedCourse && qrCodeData && (
            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <Box sx={{ textAlign: 'center', p: 2 }}>
                  <Typography variant="h6" sx={{ mb: 2 }}>
                    QR Code à scanner
                  </Typography>
                  <Box sx={{ 
                    display: 'flex', 
                    justifyContent: 'center', 
                    alignItems: 'center',
                    minHeight: 300,
                    border: '2px solid #e0e0e0',
                    borderRadius: 2,
                    p: 2
                  }}>
                    <QRCodeSVG
                      value={qrCodeData.qr_code_data}
                      size={250}
                      level="H"
                      includeMargin={true}
                    />
                  </Box>
                </Box>
              </Grid>
              <Grid item xs={12} md={6}>
                <Box sx={{ p: 2 }}>
                  <Typography variant="h6" sx={{ mb: 2 }}>
                    Informations du cours
                  </Typography>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="subtitle2" color="text.secondary">
                      Salle
                    </Typography>
                    <Typography variant="body1">
                      {rooms.find(r => r.id === selectedRoom)?.name}
                    </Typography>
                  </Box>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="subtitle2" color="text.secondary">
                      Date et heure
                    </Typography>
                    <Typography variant="body1">
                      {formatDateTime(selectedCourse.start_time)}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {formatTime(selectedCourse.start_time)} - {formatTime(selectedCourse.end_time)}
                    </Typography>
                  </Box>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="subtitle2" color="text.secondary">
                      Nom du cours
                    </Typography>
                    <Typography variant="body1">
                      {selectedCourse.name}
                    </Typography>
                  </Box>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="subtitle2" color="text.secondary">
                      Matière
                    </Typography>
                    <Typography variant="body1">
                      {selectedCourse.subject.name}
                    </Typography>
                  </Box>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="subtitle2" color="text.secondary">
                      Enseignant
                    </Typography>
                    <Typography variant="body1">
                      {selectedCourse.teacher.first_name} {selectedCourse.teacher.last_name}
                    </Typography>
                  </Box>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="subtitle2" color="text.secondary">
                      Statut
                    </Typography>
                    <Chip
                      label={getCourseStatus(selectedCourse).label}
                      color={getCourseStatus(selectedCourse).color}
                      size="small"
                      variant={getCourseStatus(selectedCourse).status === 'current' ? 'filled' : 'outlined'}
                    />
                  </Box>
                  {autoUpdate && (
                    <Alert severity="info" sx={{ mt: 2 }}>
                      Le QR code se mettra automatiquement à jour pour le prochain cours.
                    </Alert>
                  )}
                </Box>
              </Grid>
            </Grid>
          )}
        </DialogContent>
        <DialogActions>
          <Button 
            onClick={() => setQrModalOpen(false)}
            variant="outlined"
            startIcon={<Iconify icon="mdi:close" />}
          >
            Fermer
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
} 