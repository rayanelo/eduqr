import React, { useState, useEffect, useContext } from 'react';
import {
  Container,
  Typography,
  Card,
  Stack,
  Button,
  TextField,
  InputAdornment,
  Chip,
  Tooltip,
  IconButton,
  Alert,
  Box,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Grid,
  Divider
} from '@mui/material';
import { Add as AddIcon, Edit as EditIcon, Delete as DeleteIcon } from '@mui/icons-material';
import Iconify from '../../components/iconify';
import { DataTable } from '../../components/data-table';
import { ConfirmDialog } from '../../components/confirm-dialog';
import CourseFormDialog from '../../sections/courses/CourseFormDialog';
import { useCourses } from '../../hooks/useCourses';
import { usePermissions } from '../../hooks/usePermissions';
import { useSubjects } from '../../hooks/useSubjects';
import { useTeachers } from '../../hooks/useTeachers';
import { useRooms } from '../../hooks/useRooms';
import { AuthContext } from '../../auth/JwtContext';
import RepeatIcon from '@mui/icons-material/Repeat';
import InfoIcon from '@mui/icons-material/Info';

export default function CourseManagementPage() {
  const TABLE_HEAD = [
    { id: 'name', label: 'Nom du cours', align: 'left', minWidth: 200 },
    { id: 'subject', label: 'Matière', align: 'left', minWidth: 150 },
    { id: 'teacher', label: 'Enseignant', align: 'left', minWidth: 150 },
    { id: 'room', label: 'Salle', align: 'left', minWidth: 120 },
    { id: 'start_time', label: 'Date et heure', align: 'left', minWidth: 150 },
    { id: 'duration', label: 'Durée', align: 'center', minWidth: 80 },
    { id: 'recurring', label: 'Type', align: 'center', minWidth: 100 },
    { id: 'actions', label: 'Actions', align: 'center', width: 120 },
  ];

  const { user } = useContext(AuthContext);
  const { canManageCourses } = usePermissions();
  const { courses, error, setError, fetchCourses, createCourse, updateCourse, deleteCourse, checkConflicts } = useCourses();
  const { subjects, fetchSubjects } = useSubjects();
  const { teachers, fetchTeachers } = useTeachers();
  const { rooms, fetchRooms } = useRooms();

  const [openFormDialog, setOpenFormDialog] = useState(false);
  const [openDeleteDialog, setOpenDeleteDialog] = useState(false);
  const [selectedCourse, setSelectedCourse] = useState(null);
  const [filterName, setFilterName] = useState('');
  const [filterSubject, setFilterSubject] = useState('');
  const [filterTeacher, setFilterTeacher] = useState('');
  const [filterRoom, setFilterRoom] = useState('');
  const [selectedRecurringCourse, setSelectedRecurringCourse] = useState(null);
  const [showRecurringDetails, setShowRecurringDetails] = useState(false);

  useEffect(() => {
    if (user && canManageCourses) {
      fetchCourses();
      fetchSubjects();
      fetchTeachers();
      fetchRooms();
    }
  }, [fetchCourses, fetchSubjects, fetchTeachers, fetchRooms, user, canManageCourses]);

  // Filtrage des cours
  const filteredCourses = courses.filter((course) => {
    const matchesName = course.name.toLowerCase().includes(filterName.toLowerCase());
    const matchesSubject = !filterSubject || course.subject.id === parseInt(filterSubject);
    const matchesTeacher = !filterTeacher || course.teacher.id === parseInt(filterTeacher);
    const matchesRoom = !filterRoom || course.room.id === parseInt(filterRoom);
    
    return matchesName && matchesSubject && matchesTeacher && matchesRoom;
  });

  // Fonction pour regrouper les cours récurrents
  const groupRecurringCourses = (courses) => {
    const grouped = [];
    const processedRecurrenceIds = new Set();

    courses.forEach((course) => {
      if (course.is_recurring && course.recurrence_id) {
        // Si c'est un cours récurrent et qu'on ne l'a pas encore traité
        if (!processedRecurrenceIds.has(course.recurrence_id)) {
          processedRecurrenceIds.add(course.recurrence_id);
          
          // Trouver tous les cours de cette série récurrente
          const recurringSeries = courses.filter(c => 
            c.recurrence_id === course.recurrence_id
          );
          
          // Créer un cours "parent" qui représente toute la série
          const parentCourse = {
            ...course,
            is_recurring_series: true,
            series_count: recurringSeries.length,
            series_dates: recurringSeries.map(c => new Date(c.start_time)),
            series_end_date: course.recurrence_end_date,
            series_pattern: course.recurrence_pattern
          };
          
          grouped.push(parentCourse);
        }
      } else if (!course.is_recurring) {
        // Cours ponctuel - l'ajouter directement
        grouped.push(course);
      }
    });

    return grouped;
  };

  // Grouper les cours récurrents
  const groupedCourses = groupRecurringCourses(filteredCourses);

  // Fonction pour formater les informations de récurrence
  const formatRecurrenceInfo = (course) => {
    if (!course.is_recurring_series) return null;

    try {
      const pattern = JSON.parse(course.series_pattern);
      const days = pattern.days || [];
      const dayNames = {
        'Monday': 'Lun',
        'Tuesday': 'Mar', 
        'Wednesday': 'Mer',
        'Thursday': 'Jeu',
        'Friday': 'Ven',
        'Saturday': 'Sam',
        'Sunday': 'Dim'
      };
      
      const formattedDays = days.map(day => dayNames[day] || day).join(', ');
      const startDate = new Date(course.series_dates[0]).toLocaleDateString('fr-FR');
      const endDate = new Date(course.series_end_date).toLocaleDateString('fr-FR');
      
      return `${formattedDays} | ${startDate} → ${endDate} | ${course.series_count} séances`;
    } catch (e) {
      return `${course.series_count} séances récurrentes`;
    }
  };

  const dataFiltered = groupedCourses.map((course) => ({
    id: course.id,
    name: (
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <Typography 
          variant="body2" 
          sx={{ 
            fontWeight: course.is_recurring_series ? 'bold' : 'medium',
            color: course.is_recurring_series ? 'primary.main' : 'text.primary',
            '&:hover': {
              color: 'primary.main',
              textDecoration: 'underline',
              cursor: 'pointer'
            }
          }}
        >
          {course.name}
        </Typography>
        {course.is_recurring_series && (
          <RepeatIcon 
            fontSize="small" 
            sx={{ 
              color: 'primary.main',
              opacity: 0.7
            }} 
          />
        )}
      </Box>
    ),
    subject: (
      <Chip 
        label={course.subject.name} 
        size="small" 
        sx={{ 
          backgroundColor: 'secondary.light',
          color: 'secondary.contrastText',
          fontWeight: 'medium',
          fontSize: '0.75rem'
        }}
      />
    ),
    teacher: (
      <Typography variant="body2" sx={{ color: 'text.secondary', fontStyle: 'italic' }}>
        {course.teacher.first_name} {course.teacher.last_name}
      </Typography>
    ),
    room: (
      <Typography variant="body2" sx={{ color: 'text.primary', fontWeight: 'medium' }}>
        {course.room.name} - {course.room.building} {course.room.floor}
      </Typography>
    ),
    start_time: course.is_recurring_series ? (
      <Box>
        <Typography variant="body2" sx={{ fontWeight: 'medium', color: 'primary.main' }}>
          {new Date(course.start_time).toLocaleDateString('fr-FR')}
        </Typography>
        <Typography variant="caption" color="text.secondary">
          {formatRecurrenceInfo(course)}
        </Typography>
      </Box>
    ) : (
      <Typography variant="body2" sx={{ color: 'text.primary' }}>
        {new Date(course.start_time).toLocaleString('fr-FR')}
      </Typography>
    ),
    duration: (
      <Chip 
        label={`${course.duration} min`} 
        size="small" 
        sx={{ 
          backgroundColor: course.is_recurring_series ? 'primary.light' : 'grey.100',
          color: course.is_recurring_series ? 'primary.contrastText' : 'text.primary',
          fontWeight: 'medium'
        }}
      />
    ),
    recurring: course.is_recurring_series ? (
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <Chip 
          label={`${course.series_count} séances`} 
          color="primary" 
          size="small" 
          icon={<RepeatIcon />}
          sx={{ 
            backgroundColor: 'primary.main',
            color: 'primary.contrastText',
            fontWeight: 'bold'
          }}
        />
        <Tooltip title="Cours récurrent - Cliquez pour voir les détails">
          <IconButton
            onClick={() => handleShowRecurringDetails(course)}
            size="small"
            sx={{ 
              p: 0.5,
              color: 'info.main',
              '&:hover': {
                backgroundColor: 'info.light',
                color: 'info.contrastText'
              }
            }}
          >
            <InfoIcon fontSize="small" />
          </IconButton>
        </Tooltip>
      </Box>
    ) : (
      <Chip 
        label="Ponctuel" 
        size="small" 
        sx={{ 
          backgroundColor: 'success.light',
          color: 'success.contrastText',
          fontWeight: 'medium'
        }}
      />
    ),
    actions: (
      <Stack direction="row" spacing={1}>
        <Tooltip title="Modifier">
          <IconButton
            onClick={() => handleOpenFormDialog(course)}
            color="primary"
            size="small"
            sx={{
              '&:hover': {
                backgroundColor: 'primary.light',
                transform: 'scale(1.1)'
              },
              transition: 'all 0.2s ease-in-out'
            }}
          >
            <EditIcon />
          </IconButton>
        </Tooltip>
        <Tooltip title="Supprimer">
          <IconButton
            onClick={() => handleOpenDeleteDialog(course)}
            color="error"
            size="small"
            sx={{
              '&:hover': {
                backgroundColor: 'error.light',
                transform: 'scale(1.1)'
              },
              transition: 'all 0.2s ease-in-out'
            }}
          >
            <DeleteIcon />
          </IconButton>
        </Tooltip>
      </Stack>
    ),
  }));

  const handleOpenFormDialog = (course = null) => {
    setSelectedCourse(course);
    setOpenFormDialog(true);
  };

  const handleCloseFormDialog = () => {
    setOpenFormDialog(false);
    setSelectedCourse(null);
  };

  const handleOpenDeleteDialog = (course) => {
    setSelectedCourse(course);
    setOpenDeleteDialog(true);
  };

  const handleCloseDeleteDialog = () => {
    setOpenDeleteDialog(false);
    setSelectedCourse(null);
  };

  const handleShowRecurringDetails = (course) => {
    setSelectedRecurringCourse(course);
    setShowRecurringDetails(true);
  };

  const handleCloseRecurringDetails = () => {
    setShowRecurringDetails(false);
    setSelectedRecurringCourse(null);
  };

  const handleSubmitCourse = async (data) => {
    try {
      // Vérifier les conflits d'abord
      const conflictsData = await checkConflicts(data);
      
      if (conflictsData.has_conflicts) {
        return conflictsData; // Retourner les conflits pour affichage
      }

      // Créer ou mettre à jour le cours
      if (selectedCourse) {
        await updateCourse(selectedCourse.id, data);
      } else {
        await createCourse(data);
      }

      return null; // Pas de conflits
    } catch (error) {
      setError(error.message);
      throw error;
    }
  };

  const handleDeleteCourse = async () => {
    try {
      await deleteCourse(selectedCourse.id);
      handleCloseDeleteDialog();
    } catch (error) {
      setError(error.message);
    }
  };



  if (!user) {
    return (
      <Container maxWidth="xl">
        <Typography variant="h4" sx={{ mb: 5 }}>
          Chargement...
        </Typography>
      </Container>
    );
  }

  if (canManageCourses === false) {
    return (
      <Container maxWidth="xl">
        <Typography variant="h4" sx={{ mb: 5 }}>
          Accès refusé
        </Typography>
        <Alert severity="error">
          Vous n'avez pas les permissions nécessaires pour accéder à cette page.
        </Alert>
      </Container>
    );
  }

  return (
    <Container maxWidth="xl">
      <Stack direction="row" alignItems="center" justifyContent="space-between" mb={5}>
        <Typography variant="h4">Gestion des Cours</Typography>
        <Button
          variant="contained"
          color="inherit"
          startIcon={<AddIcon />}
          onClick={() => handleOpenFormDialog()}
        >
          Nouveau Cours
        </Button>
      </Stack>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }} onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      <Card>
        <Stack
          spacing={2.5}
          direction={{ xs: 'column', md: 'row' }}
          alignItems={{ xs: 'flex-end', md: 'center' }}
          justifyContent="space-between"
          sx={{
            p: 2.5,
            pr: { xs: 2.5, md: 1 },
          }}
        >
          <Stack direction="row" alignItems="center" spacing={2} flexGrow={1}>
            <TextField
              fullWidth
              value={filterName}
              onChange={(event) => setFilterName(event.target.value)}
              placeholder="Rechercher par nom..."
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <Iconify icon="eva:search-fill" sx={{ color: 'text.disabled' }} />
                  </InputAdornment>
                ),
              }}
            />
          </Stack>
        </Stack>

        <Stack
          spacing={2}
          direction={{ xs: 'column', md: 'row' }}
          sx={{ px: 2.5, pb: 2.5 }}
        >
          <FormControl sx={{ minWidth: 150 }}>
            <InputLabel>Matière</InputLabel>
            <Select
              value={filterSubject}
              onChange={(e) => setFilterSubject(e.target.value)}
              label="Matière"
            >
              <MenuItem value="">Toutes</MenuItem>
              {subjects.map((subject) => (
                <MenuItem key={subject.id} value={subject.id}>
                  {subject.name}
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          <FormControl sx={{ minWidth: 150 }}>
            <InputLabel>Enseignant</InputLabel>
            <Select
              value={filterTeacher}
              onChange={(e) => setFilterTeacher(e.target.value)}
              label="Enseignant"
            >
              <MenuItem value="">Tous</MenuItem>
              {teachers.map((teacher) => (
                <MenuItem key={teacher.id} value={teacher.id}>
                  {teacher.first_name} {teacher.last_name}
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          <FormControl sx={{ minWidth: 150 }}>
            <InputLabel>Salle</InputLabel>
            <Select
              value={filterRoom}
              onChange={(e) => setFilterRoom(e.target.value)}
              label="Salle"
            >
              <MenuItem value="">Toutes</MenuItem>
              {rooms.map((room) => (
                <MenuItem key={room.id} value={room.id}>
                  {room.name} - {room.building}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Stack>

        <DataTable
          data={dataFiltered}
          columns={TABLE_HEAD}
          onAddNew={() => handleOpenFormDialog()}
          isFiltered={!!filterName || !!filterSubject || !!filterTeacher || !!filterRoom}
        />
      </Card>

      {/* Dialog de formulaire */}
      <CourseFormDialog
        open={openFormDialog}
        onClose={handleCloseFormDialog}
        course={selectedCourse}
        onSubmit={handleSubmitCourse}
      />

      {/* Dialog de confirmation de suppression */}
      <ConfirmDialog
        open={openDeleteDialog}
        onClose={handleCloseDeleteDialog}
        onConfirm={handleDeleteCourse}
        title="Supprimer le cours"
        content={
          selectedCourse ? (
            <>
              Êtes-vous sûr de vouloir supprimer le cours "{selectedCourse.name}" ?
              {selectedCourse.is_recurring && (
                <Box sx={{ mt: 1 }}>
                  <Alert severity="warning">
                    Ce cours est récurrent. Toute la série sera supprimée.
                  </Alert>
                </Box>
              )}
            </>
          ) : (
            'Êtes-vous sûr de vouloir supprimer ce cours ?'
          )
        }
        confirmText="Supprimer"
        cancelText="Annuler"
      />

      {/* Dialog de détails des cours récurrents */}
      <Dialog
        open={showRecurringDetails}
        onClose={handleCloseRecurringDetails}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>
          Détails du cours récurrent
        </DialogTitle>
        <DialogContent>
          {selectedRecurringCourse && (
            <Box>
              <Typography variant="h6" gutterBottom>
                {selectedRecurringCourse.name}
              </Typography>
              
              <Grid container spacing={2} sx={{ mb: 3 }}>
                <Grid item xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    Matière
                  </Typography>
                  <Typography variant="body1">
                    {selectedRecurringCourse.subject.name}
                  </Typography>
                </Grid>
                <Grid item xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    Enseignant
                  </Typography>
                  <Typography variant="body1">
                    {selectedRecurringCourse.teacher.first_name} {selectedRecurringCourse.teacher.last_name}
                  </Typography>
                </Grid>
                <Grid item xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    Salle
                  </Typography>
                  <Typography variant="body1">
                    {selectedRecurringCourse.room.name} - {selectedRecurringCourse.room.building} {selectedRecurringCourse.room.floor}
                  </Typography>
                </Grid>
                <Grid item xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    Durée
                  </Typography>
                  <Typography variant="body1">
                    {selectedRecurringCourse.duration} minutes
                  </Typography>
                </Grid>
              </Grid>

              <Divider sx={{ my: 2 }} />

              <Typography variant="h6" gutterBottom>
                Informations de récurrence
              </Typography>
              
              <Grid container spacing={2} sx={{ mb: 3 }}>
                <Grid item xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    Nombre de séances
                  </Typography>
                  <Typography variant="body1">
                    {selectedRecurringCourse.series_count}
                  </Typography>
                </Grid>
                <Grid item xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    Période
                  </Typography>
                  <Typography variant="body1">
                    {new Date(selectedRecurringCourse.series_dates[0]).toLocaleDateString('fr-FR')} → {new Date(selectedRecurringCourse.series_end_date).toLocaleDateString('fr-FR')}
                  </Typography>
                </Grid>
              </Grid>

              <Typography variant="body2" color="text.secondary" gutterBottom>
                Jours de la semaine :
              </Typography>
              <Box sx={{ mb: 2 }}>
                {(() => {
                  try {
                    const pattern = JSON.parse(selectedRecurringCourse.series_pattern);
                    const days = pattern.days || [];
                    const dayNames = {
                      'Monday': 'Lundi',
                      'Tuesday': 'Mardi', 
                      'Wednesday': 'Mercredi',
                      'Thursday': 'Jeudi',
                      'Friday': 'Vendredi',
                      'Saturday': 'Samedi',
                      'Sunday': 'Dimanche'
                    };
                    
                    return days.map(day => (
                      <Chip 
                        key={day} 
                        label={dayNames[day] || day} 
                        size="small" 
                        sx={{ mr: 1, mb: 1 }}
                      />
                    ));
                  } catch (e) {
                    return <Typography variant="body2">Informations non disponibles</Typography>;
                  }
                })()}
              </Box>

              <Typography variant="body2" color="text.secondary" gutterBottom>
                Toutes les séances :
              </Typography>
              <Box sx={{ maxHeight: 200, overflow: 'auto' }}>
                {selectedRecurringCourse.series_dates.map((date, index) => (
                  <Typography key={index} variant="body2" sx={{ py: 0.5 }}>
                    Séance {index + 1} : {new Date(date).toLocaleString('fr-FR')}
                  </Typography>
                ))}
              </Box>
            </Box>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseRecurringDetails}>
            Fermer
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
} 