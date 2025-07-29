import { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  FormControlLabel,
  Switch,
  Box,
  Chip,
  Alert,
  Grid,
  Typography,
  Divider
} from '@mui/material';
import { DateTimePicker } from '@mui/x-date-pickers/DateTimePicker';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { fr } from 'date-fns/locale';
import { useForm, Controller } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import { useSubjects } from '../../hooks/useSubjects';
import { useTeachers } from '../../hooks/useTeachers';
import { useRooms } from '../../hooks/useRooms';

const schema = yup.object().shape({
  name: yup.string().required('Le nom du cours est requis'),
  subject_id: yup.number().required('La matière est requise'),
  teacher_id: yup.number().required('L\'enseignant est requis'),
  room_id: yup.number().required('La salle est requise'),
  start_time: yup.date().required('La date et heure de début sont requises'),
  duration: yup.number().min(15, 'Durée minimum 15 minutes').max(480, 'Durée maximum 8 heures').required('La durée est requise'),
  description: yup.string(),
  is_recurring: yup.boolean(),
  recurrence_pattern: yup.string().when('is_recurring', {
    is: true,
    then: yup.string().required('Le pattern de récurrence est requis')
  }),
  recurrence_end_date: yup.date().when('is_recurring', {
    is: true,
    then: yup.date()
      .required('La date de fin de récurrence est requise')
      .test('is-after-start', 'La date de fin doit être après la date de début', function(value) {
        const startTime = this.parent.start_time;
        if (!startTime || !value) return true;
        return value > startTime;
      })
  }),
  exclude_holidays: yup.boolean()
});

const DAYS_OF_WEEK = [
  { value: 'Monday', label: 'Lundi' },
  { value: 'Tuesday', label: 'Mardi' },
  { value: 'Wednesday', label: 'Mercredi' },
  { value: 'Thursday', label: 'Jeudi' },
  { value: 'Friday', label: 'Vendredi' },
  { value: 'Saturday', label: 'Samedi' },
  { value: 'Sunday', label: 'Dimanche' }
];

export default function CourseFormDialog({ open, onClose, course = null, onSubmit, initialData = null }) {
  const [selectedDays, setSelectedDays] = useState([]);
  const [conflicts, setConflicts] = useState([]);
  const [showConflicts, setShowConflicts] = useState(false);

  const { subjects, fetchSubjects } = useSubjects();
  const { teachers, fetchTeachers } = useTeachers();
  const { rooms, fetchRooms } = useRooms();

  const {
    control,
    handleSubmit,
    reset,
    watch,
    setValue,
    formState: { errors, isSubmitting }
  } = useForm({
    resolver: yupResolver(schema),
    defaultValues: {
      name: '',
      subject_id: '',
      teacher_id: '',
      room_id: '',
      start_time: new Date(),
      duration: 60,
      description: '',
      is_recurring: false,
      recurrence_pattern: '',
      recurrence_end_date: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000), // 30 jours dans le futur
      exclude_holidays: true
    }
  });

  const isRecurring = watch('is_recurring');

  useEffect(() => {
    if (open) {
      fetchSubjects();
      fetchTeachers();
      fetchRooms();
      
      if (course) {
        reset({
          name: course.name,
          subject_id: course.subject.id,
          teacher_id: course.teacher.id,
          room_id: course.room.id,
          start_time: new Date(course.start_time),
          duration: course.duration,
          description: course.description || '',
          is_recurring: course.is_recurring,
          recurrence_pattern: course.recurrence_pattern || '',
          recurrence_end_date: course.recurrence_end_date ? new Date(course.recurrence_end_date) : new Date(),
          exclude_holidays: course.exclude_holidays
        });

        if (course.recurrence_pattern) {
          try {
            const pattern = JSON.parse(course.recurrence_pattern);
            setSelectedDays(pattern.days || []);
          } catch (e) {
            setSelectedDays([]);
          }
        }
      } else {
        // Utiliser les données initiales si disponibles
        const defaultValues = {
          name: '',
          subject_id: '',
          teacher_id: '',
          room_id: '',
          start_time: initialData?.start_time || new Date(),
          duration: 60,
          description: '',
          is_recurring: false,
          recurrence_pattern: '',
          recurrence_end_date: initialData?.end_time || new Date(Date.now() + 30 * 24 * 60 * 60 * 1000), // 30 jours dans le futur
          exclude_holidays: true
        };
        
        reset(defaultValues);
        setSelectedDays([]);
      }
    }
  }, [open, course, reset, fetchSubjects, fetchTeachers, fetchRooms, initialData]);

  const handleDayToggle = (day) => {
    const newSelectedDays = selectedDays.includes(day)
      ? selectedDays.filter(d => d !== day)
      : [...selectedDays, day];
    
    setSelectedDays(newSelectedDays);
    setValue('recurrence_pattern', JSON.stringify({ days: newSelectedDays }));
  };

  const handleFormSubmit = async (data) => {
    // Éviter les soumissions multiples
    if (isSubmitting) {
      return;
    }

    try {
      // Soumettre le formulaire (inclut la vérification des conflits)
      const conflictsData = await onSubmit(data);
      if (conflictsData && conflictsData.has_conflicts) {
        setConflicts(conflictsData.data);
        setShowConflicts(true);
        return;
      }

      // Si pas de conflits, fermer le dialogue et réinitialiser
      onClose();
      reset();
      setSelectedDays([]);
      setConflicts([]);
      setShowConflicts(false);
    } catch (error) {
      console.error('Erreur lors de la soumission:', error);
    }
  };

  const handleClose = () => {
    onClose();
    reset();
    setSelectedDays([]);
    setConflicts([]);
    setShowConflicts(false);
  };

  // Auto-complétion pour le nom du cours basé sur la matière sélectionnée
  const selectedSubjectId = watch('subject_id');
  const selectedSubject = subjects.find(subject => subject.id === selectedSubjectId);
  const suggestedCourseName = selectedSubject ? `${selectedSubject.name} - Cours` : '';

  // Mettre à jour le nom du cours quand la matière change
  useEffect(() => {
    if (selectedSubject && !watch('name')) {
      setValue('name', suggestedCourseName);
    }
  }, [selectedSubject, setValue, watch, suggestedCourseName]);

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
      <DialogTitle>
        {course ? 'Modifier le cours' : 'Créer un nouveau cours'}
      </DialogTitle>
      
      <form onSubmit={handleSubmit(handleFormSubmit)}>
        <DialogContent>
          <Grid container spacing={3}>
            {/* Informations de base */}


            <Grid item xs={12} md={6}>
              <Controller
                name="name"
                control={control}
                render={({ field }) => (
                  <TextField
                    {...field}
                    fullWidth
                    label="Nom du cours"
                    error={!!errors.name}
                    helperText={errors.name?.message || (selectedSubject && `Suggestion: ${suggestedCourseName}`)}
                    placeholder={suggestedCourseName}
                  />
                )}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <Controller
                name="subject_id"
                control={control}
                render={({ field }) => (
                  <FormControl fullWidth error={!!errors.subject_id}>
                    <InputLabel>Matière</InputLabel>
                    <Select {...field} label="Matière">
                      {subjects.map((subject) => (
                        <MenuItem key={subject.id} value={subject.id}>
                          {subject.name}
                        </MenuItem>
                      ))}
                    </Select>
                    {errors.subject_id && (
                      <Typography color="error" variant="caption">
                        {errors.subject_id.message}
                      </Typography>
                    )}
                  </FormControl>
                )}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <Controller
                name="teacher_id"
                control={control}
                render={({ field }) => (
                  <FormControl fullWidth error={!!errors.teacher_id}>
                    <InputLabel>Enseignant</InputLabel>
                    <Select {...field} label="Enseignant">
                      {teachers.map((teacher) => (
                        <MenuItem key={teacher.id} value={teacher.id}>
                          {teacher.first_name} {teacher.last_name}
                        </MenuItem>
                      ))}
                    </Select>
                    {errors.teacher_id && (
                      <Typography color="error" variant="caption">
                        {errors.teacher_id.message}
                      </Typography>
                    )}
                  </FormControl>
                )}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <Controller
                name="room_id"
                control={control}
                render={({ field }) => (
                  <FormControl fullWidth error={!!errors.room_id}>
                    <InputLabel>Salle</InputLabel>
                    <Select {...field} label="Salle">
                      {rooms.map((room) => (
                        <MenuItem key={room.id} value={room.id}>
                          {room.name} - {room.building} {room.floor}
                        </MenuItem>
                      ))}
                    </Select>
                    {errors.room_id && (
                      <Typography color="error" variant="caption">
                        {errors.room_id.message}
                      </Typography>
                    )}
                  </FormControl>
                )}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <LocalizationProvider dateAdapter={AdapterDateFns} adapterLocale={fr}>
                <Controller
                  name="start_time"
                  control={control}
                  render={({ field }) => (
                    <DateTimePicker
                      {...field}
                      label="Date et heure de début"
                      slotProps={{
                        textField: {
                          fullWidth: true,
                          error: !!errors.start_time,
                          helperText: errors.start_time?.message
                        }
                      }}
                    />
                  )}
                />
              </LocalizationProvider>
            </Grid>

            <Grid item xs={12} md={6}>
              <Controller
                name="duration"
                control={control}
                render={({ field }) => (
                  <TextField
                    {...field}
                    fullWidth
                    type="number"
                    label="Durée (minutes)"
                    error={!!errors.duration}
                    helperText={errors.duration?.message}
                  />
                )}
              />
            </Grid>

            <Grid item xs={12}>
              <Controller
                name="description"
                control={control}
                render={({ field }) => (
                  <TextField
                    {...field}
                    fullWidth
                    multiline
                    rows={3}
                    label="Description (optionnel)"
                    error={!!errors.description}
                    helperText={errors.description?.message}
                  />
                )}
              />
            </Grid>

            {/* Récurrence */}
            <Grid item xs={12}>
              <Divider sx={{ my: 2 }} />
              <Typography variant="h6" gutterBottom>
                Récurrence
              </Typography>
            </Grid>

            <Grid item xs={12}>
              <Controller
                name="is_recurring"
                control={control}
                render={({ field }) => (
                  <FormControlLabel
                    control={
                      <Switch
                        checked={field.value}
                        onChange={field.onChange}
                      />
                    }
                    label="Cours récurrent"
                  />
                )}
              />
            </Grid>

            {isRecurring && (
              <>
                <Grid item xs={12}>
                  <Typography variant="subtitle2" gutterBottom>
                    Jours de répétition
                  </Typography>
                  <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
                    {DAYS_OF_WEEK.map((day) => (
                      <Chip
                        key={day.value}
                        label={day.label}
                        onClick={() => handleDayToggle(day.value)}
                        color={selectedDays.includes(day.value) ? 'primary' : 'default'}
                        variant={selectedDays.includes(day.value) ? 'filled' : 'outlined'}
                      />
                    ))}
                  </Box>
                  {errors.recurrence_pattern && (
                    <Typography color="error" variant="caption" display="block">
                      {errors.recurrence_pattern.message}
                    </Typography>
                  )}
                </Grid>

                <Grid item xs={12} md={6}>
                  <LocalizationProvider dateAdapter={AdapterDateFns} adapterLocale={fr}>
                    <Controller
                      name="recurrence_end_date"
                      control={control}
                      render={({ field }) => (
                        <DateTimePicker
                          {...field}
                          label="Date de fin de récurrence"
                          slotProps={{
                            textField: {
                              fullWidth: true,
                              error: !!errors.recurrence_end_date,
                              helperText: errors.recurrence_end_date?.message
                            }
                          }}
                        />
                      )}
                    />
                  </LocalizationProvider>
                </Grid>

                <Grid item xs={12} md={6}>
                  <Controller
                    name="exclude_holidays"
                    control={control}
                    render={({ field }) => (
                      <FormControlLabel
                        control={
                          <Switch
                            checked={field.value}
                            onChange={field.onChange}
                          />
                        }
                        label="Exclure les jours fériés"
                      />
                    )}
                  />
                </Grid>
              </>
            )}

            {/* Conflits */}
            {showConflicts && conflicts.length > 0 && (
              <Grid item xs={12}>
                <Alert severity="warning" sx={{ mt: 2 }}>
                  <Typography variant="subtitle2" gutterBottom>
                    Conflits détectés :
                  </Typography>
                  {conflicts.map((conflict, index) => (
                    <Typography key={index} variant="body2">
                      • {conflict.course_name} - {conflict.room_name} - {new Date(conflict.start_time).toLocaleString('fr-FR')}
                    </Typography>
                  ))}
                </Alert>
              </Grid>
            )}
          </Grid>
        </DialogContent>

        <DialogActions>
          <Button onClick={handleClose}>Annuler</Button>
          <Button
            type="submit"
            variant="contained"
            disabled={isSubmitting}
          >
            {isSubmitting ? 'Enregistrement...' : (course ? 'Modifier' : 'Créer')}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
} 