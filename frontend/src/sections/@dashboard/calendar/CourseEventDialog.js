import PropTypes from 'prop-types';
import { useState } from 'react';
// @mui
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Box,
  Typography,
  Button,
  IconButton,
  Tooltip,
  Alert,
  Chip,
  Stack,
  TextField,
  Avatar,
  Grid,
  Card,
  CardContent,
} from '@mui/material';
// hooks
import { useSnackbar } from '../../../components/snackbar';
import { usePermissions } from '../../../hooks/usePermissions';
// components
import Iconify from '../../../components/iconify';
import apiClient from '../../../utils/api';

// ----------------------------------------------------------------------

CourseEventDialog.propTypes = {
  open: PropTypes.bool,
  onClose: PropTypes.func,
  event: PropTypes.object,
  onUpdate: PropTypes.func,
  onDelete: PropTypes.func,
};

export default function CourseEventDialog({ open, onClose, event, onUpdate, onDelete }) {
  const { enqueueSnackbar } = useSnackbar();
  const { user } = usePermissions();
  const [comment, setComment] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  if (!event) return null;

  const course = event.extendedProps?.course;
  const isAdmin = user?.role === 'admin' || user?.role === 'super_admin';
  const isTeacher = user?.role === 'professeur';
  const isStudent = user?.role === 'etudiant';
  const isMyCourse = isTeacher && course?.teacher?.id === user?.id;
  const isRecurring = event.extendedProps?.isRecurring || false;

  const handleAddComment = async () => {
    if (!comment.trim()) return;

    setIsSubmitting(true);
    try {
      await apiClient.post(`/api/v1/admin/courses/${event.extendedProps.course.id}/comments`, {
        comment: comment.trim(),
        teacher_id: user.id,
      });
      
      enqueueSnackbar('Commentaire ajouté avec succès', { variant: 'success' });
      setComment('');
      onClose();
    } catch (error) {
      enqueueSnackbar('Erreur lors de l\'ajout du commentaire', { variant: 'error' });
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDeleteCourse = async () => {
    setIsSubmitting(true);
    try {
      await onDelete();
      setShowDeleteConfirm(false);
      onClose();
    } catch (error) {
      enqueueSnackbar('Erreur lors de la suppression du cours', { variant: 'error' });
    } finally {
      setIsSubmitting(false);
    }
  };

  const formatTime = (timeString) => {
    return new Date(timeString).toLocaleTimeString('fr-FR', {
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('fr-FR', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  const formatDuration = (startTime, endTime) => {
    const start = new Date(startTime);
    const end = new Date(endTime);
    const diffMs = end - start;
    const diffMins = Math.round(diffMs / 60000);
    const hours = Math.floor(diffMins / 60);
    const mins = diffMins % 60;
    
    if (hours > 0) {
      return `${hours}h${mins > 0 ? ` ${mins}min` : ''}`;
    }
    return `${mins}min`;
  };

  const getRoleColor = (role) => {
    switch (role) {
      case 'admin':
      case 'super_admin':
        return 'error';
      case 'professeur':
        return 'warning';
      case 'etudiant':
        return 'info';
      default:
        return 'default';
    }
  };

  const getRoleLabel = (role) => {
    switch (role) {
      case 'admin':
        return 'Admin';
      case 'super_admin':
        return 'Super Admin';
      case 'professeur':
        return 'Professeur';
      case 'etudiant':
        return 'Étudiant';
      default:
        return role;
    }
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle sx={{ pb: 1, px: 3, pt: 2 }}>
        <Box display="flex" alignItems="center" justifyContent="space-between">
          <Box display="flex" alignItems="center" gap={1.5}>
            <Box
              sx={{
                width: 8,
                height: 8,
                borderRadius: '50%',
                backgroundColor: '#2196F3',
                flexShrink: 0,
              }}
            />
            <Typography variant="h6" component="h2" sx={{ fontWeight: 600 }}>
              {event.title}
            </Typography>
            {isRecurring && (
              <Chip
                label="Récurrent"
                color="secondary"
                size="small"
                icon={<Iconify icon="eva:refresh-fill" />}
                sx={{ height: 24, fontSize: '0.75rem' }}
              />
            )}
          </Box>
          <Box>
            {isAdmin && (
              <>
                {!isRecurring && (
                  <Tooltip title="Modifier le cours">
                    <IconButton 
                      color="primary" 
                      size="small"
                      onClick={() => {
                        window.location.href = `/dashboard/courses/edit/${event.extendedProps.course.id}`;
                      }}
                      sx={{ mr: 0.5 }}
                    >
                      <Iconify icon="eva:edit-fill" />
                    </IconButton>
                  </Tooltip>
                )}
                <Tooltip title="Supprimer le cours">
                  <IconButton 
                    color="error" 
                    size="small"
                    onClick={() => setShowDeleteConfirm(true)}
                  >
                    <Iconify icon="eva:trash-2-fill" />
                  </IconButton>
                </Tooltip>
              </>
            )}
          </Box>
        </Box>
      </DialogTitle>

      <DialogContent sx={{ pt: 1, px: 3, pb: 2 }}>
        {showDeleteConfirm ? (
          <Alert severity="warning" sx={{ mb: 2 }}>
            <Typography variant="body1" sx={{ fontWeight: 600, mb: 1 }}>
              Confirmer la suppression
            </Typography>
            <Typography variant="body2" sx={{ mb: 2 }}>
              Êtes-vous sûr de vouloir supprimer ce cours ? Cette action est irréversible.
            </Typography>
            <Box>
              <Button
                variant="contained"
                color="error"
                size="small"
                onClick={handleDeleteCourse}
                disabled={isSubmitting}
                sx={{ mr: 1 }}
              >
                Confirmer la suppression
              </Button>
              <Button
                variant="outlined"
                size="small"
                onClick={() => setShowDeleteConfirm(false)}
                disabled={isSubmitting}
              >
                Annuler
              </Button>
            </Box>
          </Alert>
        ) : (
          <Box>
            {/* Informations principales - Layout compact */}
            <Grid container spacing={2} sx={{ mb: 2 }}>
              <Grid item xs={12} sm={6}>
                <Card variant="outlined" sx={{ height: '100%', boxShadow: 'none', borderColor: 'divider' }}>
                  <CardContent sx={{ p: 2, '&:last-child': { pb: 2 } }}>
                    <Typography variant="subtitle2" sx={{ mb: 1.5, display: 'flex', alignItems: 'center', gap: 0.5, color: 'primary.main', fontWeight: 600 }}>
                      <Iconify icon="eva:clock-fill" sx={{ fontSize: 16 }} />
                      Horaires
                    </Typography>
                    <Stack spacing={0.5}>
                      <Box display="flex" justifyContent="space-between">
                        <Typography variant="caption" color="text.secondary">Date:</Typography>
                        <Typography variant="body2" sx={{ fontWeight: 500 }}>
                          {formatDate(event.start)}
                        </Typography>
                      </Box>
                      <Box display="flex" justifyContent="space-between">
                        <Typography variant="caption" color="text.secondary">Heure:</Typography>
                        <Typography variant="body2" sx={{ fontWeight: 500 }}>
                          {formatTime(event.start)} - {formatTime(event.end)}
                        </Typography>
                      </Box>
                      <Box display="flex" justifyContent="space-between">
                        <Typography variant="caption" color="text.secondary">Durée:</Typography>
                        <Typography variant="body2" sx={{ fontWeight: 500 }}>
                          {formatDuration(event.start, event.end)}
                        </Typography>
                      </Box>
                    </Stack>
                  </CardContent>
                </Card>
              </Grid>

              <Grid item xs={12} sm={6}>
                <Card variant="outlined" sx={{ height: '100%', boxShadow: 'none', borderColor: 'divider' }}>
                  <CardContent sx={{ p: 2, '&:last-child': { pb: 2 } }}>
                    <Typography variant="subtitle2" sx={{ mb: 1.5, display: 'flex', alignItems: 'center', gap: 0.5, color: 'primary.main', fontWeight: 600 }}>
                      <Iconify icon="eva:map-pin-fill" sx={{ fontSize: 16 }} />
                      Localisation
                    </Typography>
                    <Stack spacing={0.5}>
                      <Box display="flex" justifyContent="space-between">
                        <Typography variant="caption" color="text.secondary">Salle:</Typography>
                        <Typography variant="body2" sx={{ fontWeight: 500 }}>
                          {event.extendedProps.room}
                        </Typography>
                      </Box>
                      <Box display="flex" justifyContent="space-between">
                        <Typography variant="caption" color="text.secondary">Bâtiment:</Typography>
                        <Typography variant="body2" sx={{ fontWeight: 500 }}>
                          {course?.room?.building || 'Non spécifié'}
                        </Typography>
                      </Box>
                      <Box display="flex" justifyContent="space-between">
                        <Typography variant="caption" color="text.secondary">Étage:</Typography>
                        <Typography variant="body2" sx={{ fontWeight: 500 }}>
                          {course?.room?.floor || 'Non spécifié'}
                        </Typography>
                      </Box>
                    </Stack>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>

            {/* Informations sur le cours - Layout compact */}
            <Grid container spacing={2} sx={{ mb: 2 }}>
              <Grid item xs={12} sm={6}>
                <Card variant="outlined" sx={{ height: '100%', boxShadow: 'none', borderColor: 'divider' }}>
                  <CardContent sx={{ p: 2, '&:last-child': { pb: 2 } }}>
                    <Typography variant="subtitle2" sx={{ mb: 1.5, display: 'flex', alignItems: 'center', gap: 0.5, color: 'primary.main', fontWeight: 600 }}>
                      <Iconify icon="eva:book-fill" sx={{ fontSize: 16 }} />
                      Matière
                    </Typography>
                    <Typography variant="body2" sx={{ fontWeight: 500, mb: 0.5 }}>
                      {event.extendedProps.subject}
                    </Typography>
                    {course?.subject?.description && (
                      <Typography variant="caption" color="text.secondary" sx={{ lineHeight: 1.2 }}>
                        {course.subject.description}
                      </Typography>
                    )}
                  </CardContent>
                </Card>
              </Grid>

              <Grid item xs={12} sm={6}>
                <Card variant="outlined" sx={{ height: '100%', boxShadow: 'none', borderColor: 'divider' }}>
                  <CardContent sx={{ p: 2, '&:last-child': { pb: 2 } }}>
                    <Typography variant="subtitle2" sx={{ mb: 1.5, display: 'flex', alignItems: 'center', gap: 0.5, color: 'primary.main', fontWeight: 600 }}>
                      <Iconify icon="eva:person-fill" sx={{ fontSize: 16 }} />
                      Professeur
                    </Typography>
                    <Box display="flex" alignItems="center" gap={1.5}>
                      <Avatar
                        src={course?.teacher?.avatar_url || '/assets/images/avatars/default-avatar.png'}
                        alt={event.extendedProps.teacher}
                        sx={{ width: 36, height: 36 }}
                      />
                      <Box>
                        <Typography variant="body2" sx={{ fontWeight: 500, mb: 0.5 }}>
                          {event.extendedProps.teacher}
                        </Typography>
                        <Chip
                          label={getRoleLabel(course?.teacher?.role)}
                          color={getRoleColor(course?.teacher?.role)}
                          size="small"
                          sx={{ height: 20, fontSize: '0.7rem' }}
                        />
                      </Box>
                    </Box>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>

            {/* Description - Conditionnelle et compacte */}
            {event.extendedProps.description && (
              <Card variant="outlined" sx={{ mb: 2, boxShadow: 'none', borderColor: 'divider' }}>
                <CardContent sx={{ p: 2, '&:last-child': { pb: 2 } }}>
                  <Typography variant="subtitle2" sx={{ mb: 1, display: 'flex', alignItems: 'center', gap: 0.5, color: 'primary.main', fontWeight: 600 }}>
                    <Iconify icon="eva:file-text-fill" sx={{ fontSize: 16 }} />
                    Description
                  </Typography>
                  <Typography variant="body2" sx={{ lineHeight: 1.4 }}>
                    {event.extendedProps.description}
                  </Typography>
                </CardContent>
              </Card>
            )}

            {/* Commentaires pour les professeurs - Compact */}
            {isTeacher && isMyCourse && (
              <Card variant="outlined" sx={{ mb: 2, boxShadow: 'none', borderColor: 'divider' }}>
                <CardContent sx={{ p: 2, '&:last-child': { pb: 2 } }}>
                  <Typography variant="subtitle2" sx={{ mb: 1.5, display: 'flex', alignItems: 'center', gap: 0.5, color: 'primary.main', fontWeight: 600 }}>
                    <Iconify icon="eva:message-circle-fill" sx={{ fontSize: 16 }} />
                    Ajouter un commentaire
                  </Typography>
                  <TextField
                    fullWidth
                    multiline
                    rows={2}
                    size="small"
                    placeholder="Ajoutez un commentaire pour ce cours..."
                    value={comment}
                    onChange={(e) => setComment(e.target.value)}
                    sx={{ mb: 1.5 }}
                  />
                  <Button
                    variant="contained"
                    size="small"
                    onClick={handleAddComment}
                    disabled={!comment.trim() || isSubmitting}
                    startIcon={<Iconify icon="eva:send-fill" />}
                  >
                    Ajouter le commentaire
                  </Button>
                </CardContent>
              </Card>
            )}

            {/* Informations contextuelles - Compactes */}
            {isAdmin && (
              <Alert severity="info" sx={{ mt: 1.5, py: 1 }}>
                <Typography variant="caption">
                  En tant qu'administrateur, vous pouvez {isRecurring ? 'supprimer' : 'modifier ou supprimer'} ce cours.
                  {isRecurring && ' Les cours récurrents ne peuvent pas être modifiés directement.'}
                </Typography>
              </Alert>
            )}

            {isStudent && (
              <Alert severity="info" sx={{ mt: 1.5, py: 1 }}>
                <Typography variant="caption">
                  Vous pouvez consulter les informations de ce cours. Pour toute question, contactez votre professeur.
                </Typography>
              </Alert>
            )}
          </Box>
        )}
      </DialogContent>

      <DialogActions sx={{ px: 3, pb: 2, pt: 1 }}>
        <Button variant="outlined" color="inherit" size="small" onClick={onClose}>
          Fermer
        </Button>
      </DialogActions>
    </Dialog>
  );
} 