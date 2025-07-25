import PropTypes from 'prop-types';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  Box,
  Alert,
  List,
  ListItem,
  ListItemText,
  Chip,
} from '@mui/material';
import { LoadingButton } from '@mui/lab';

const DeleteConfirmDialog = ({
  open,
  onClose,
  onConfirm,
  title,
  message,
  resourceName,
  resourceType,
  isDeleting,
  warnings = [],
  conflicts = [],
  isRecurring = false,
  onDeleteRecurring,
  futureCourses = [],
  pastCourses = [],
  onConfirmWithCourses,
}) => {
  const getResourceTypeText = () => {
    switch (resourceType) {
      case 'user':
        return 'utilisateur';
      case 'room':
        return 'salle';
      case 'subject':
        return 'matière';
      case 'course':
        return 'cours';
      default:
        return 'élément';
    }
  };

  const getDeleteButtonText = () => {
    if (isRecurring) {
      return 'Supprimer la série complète';
    }
    return `Supprimer ${getResourceTypeText()}`;
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle sx={{ color: 'error.main' }}>
        {title || `Supprimer ${getResourceTypeText()}`}
      </DialogTitle>
      
      <DialogContent>
        <Typography variant="body1" sx={{ mb: 2 }}>
          {message || `Êtes-vous sûr de vouloir supprimer ${resourceName || `cet ${getResourceTypeText()}`} ?`}
        </Typography>

        {/* Avertissements */}
        {warnings.length > 0 && (
          <Alert severity="warning" sx={{ mb: 2 }}>
            <Typography variant="subtitle2" sx={{ mb: 1 }}>
              Avertissements :
            </Typography>
            <List dense>
              {warnings.map((warning, index) => (
                <ListItem key={index} sx={{ py: 0 }}>
                  <ListItemText primary={warning} />
                </ListItem>
              ))}
            </List>
          </Alert>
        )}

        {/* Conflits empêchant la suppression */}
        {conflicts.length > 0 && (
          <Alert severity="error" sx={{ mb: 2 }}>
            <Typography variant="subtitle2" sx={{ mb: 1 }}>
              Conflits détectés :
            </Typography>
            <List dense>
              {conflicts.map((conflict, index) => (
                <ListItem key={index} sx={{ py: 0 }}>
                  <ListItemText 
                    primary={conflict.courseName || conflict.roomName || conflict.subjectName}
                    secondary={conflict.date || conflict.description}
                  />
                </ListItem>
              ))}
            </List>
          </Alert>
        )}

        {/* Cours liés à supprimer */}
        {(futureCourses.length > 0 || pastCourses.length > 0) && (
          <Alert severity="warning" sx={{ mb: 2 }}>
            <Typography variant="subtitle2" sx={{ mb: 1 }}>
              Cours liés qui seront supprimés :
            </Typography>
            
            {futureCourses.length > 0 && (
              <Box sx={{ mb: 2 }}>
                <Typography variant="body2" color="error.main" sx={{ mb: 1 }}>
                  Cours futurs ({futureCourses.length}) :
                </Typography>
                <List dense>
                  {futureCourses.slice(0, 5).map((course, index) => (
                    <ListItem key={index} sx={{ py: 0 }}>
                      <ListItemText 
                        primary={course.name}
                        secondary={`${course.subject?.name} - ${course.room?.name} - ${new Date(course.start_time).toLocaleDateString()}`}
                      />
                    </ListItem>
                  ))}
                  {futureCourses.length > 5 && (
                    <ListItem sx={{ py: 0 }}>
                      <ListItemText 
                        primary={`... et ${futureCourses.length - 5} autres cours futurs`}
                        sx={{ fontStyle: 'italic' }}
                      />
                    </ListItem>
                  )}
                </List>
              </Box>
            )}
            
            {pastCourses.length > 0 && (
              <Box>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                  Cours passés ({pastCourses.length}) :
                </Typography>
                <List dense>
                  {pastCourses.slice(0, 3).map((course, index) => (
                    <ListItem key={index} sx={{ py: 0 }}>
                      <ListItemText 
                        primary={course.name}
                        secondary={`${course.subject?.name} - ${course.room?.name} - ${new Date(course.start_time).toLocaleDateString()}`}
                      />
                    </ListItem>
                  ))}
                  {pastCourses.length > 3 && (
                    <ListItem sx={{ py: 0 }}>
                      <ListItemText 
                        primary={`... et ${pastCourses.length - 3} autres cours passés`}
                        sx={{ fontStyle: 'italic' }}
                      />
                    </ListItem>
                  )}
                </List>
              </Box>
            )}
          </Alert>
        )}

        {/* Options pour les cours récurrents */}
        {isRecurring && conflicts.length === 0 && (
          <Box sx={{ mt: 2 }}>
            <Typography variant="subtitle2" sx={{ mb: 1 }}>
              Ce cours fait partie d'une série récurrente. Que souhaitez-vous faire ?
            </Typography>
            <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
              <Chip 
                label="Supprimer toute la série" 
                color="error" 
                variant="outlined"
                onClick={onConfirm}
                disabled={isDeleting}
              />
              <Chip 
                label="Supprimer cette occurrence seulement" 
                color="warning" 
                variant="outlined"
                onClick={onDeleteRecurring}
                disabled={isDeleting}
              />
            </Box>
          </Box>
        )}
      </DialogContent>

      <DialogActions sx={{ p: 3, pt: 1 }}>
        <Button onClick={onClose} disabled={isDeleting}>
          Annuler
        </Button>
        
        {/* Boutons pour les cours liés */}
        {(futureCourses.length > 0 || pastCourses.length > 0) && conflicts.length === 0 && (
          <Box sx={{ display: 'flex', gap: 1 }}>
            <LoadingButton
              onClick={onConfirm}
              loading={isDeleting}
              variant="outlined"
              color="warning"
            >
              Supprimer sans les cours
            </LoadingButton>
            <LoadingButton
              onClick={onConfirmWithCourses}
              loading={isDeleting}
              variant="contained"
              color="error"
            >
              Supprimer avec les cours
            </LoadingButton>
          </Box>
        )}
        
        {/* Bouton normal pour les autres cas */}
        {!isRecurring && conflicts.length === 0 && (futureCourses.length === 0 && pastCourses.length === 0) && (
          <LoadingButton
            onClick={onConfirm}
            loading={isDeleting}
            variant="contained"
            color="error"
          >
            {getDeleteButtonText()}
          </LoadingButton>
        )}
      </DialogActions>
    </Dialog>
  );
};

DeleteConfirmDialog.propTypes = {
  open: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  onConfirm: PropTypes.func.isRequired,
  title: PropTypes.string,
  message: PropTypes.string,
  resourceName: PropTypes.string,
  resourceType: PropTypes.oneOf(['user', 'room', 'subject', 'course']).isRequired,
  isDeleting: PropTypes.bool,
  warnings: PropTypes.arrayOf(PropTypes.string),
  conflicts: PropTypes.arrayOf(PropTypes.shape({
    courseName: PropTypes.string,
    roomName: PropTypes.string,
    subjectName: PropTypes.string,
    date: PropTypes.string,
    description: PropTypes.string,
  })),
  isRecurring: PropTypes.bool,
  onDeleteRecurring: PropTypes.func,
  futureCourses: PropTypes.arrayOf(PropTypes.shape({
    id: PropTypes.number,
    name: PropTypes.string,
    subject: PropTypes.shape({
      name: PropTypes.string,
    }),
    room: PropTypes.shape({
      name: PropTypes.string,
    }),
    start_time: PropTypes.string,
  })),
  pastCourses: PropTypes.arrayOf(PropTypes.shape({
    id: PropTypes.number,
    name: PropTypes.string,
    subject: PropTypes.shape({
      name: PropTypes.string,
    }),
    room: PropTypes.shape({
      name: PropTypes.string,
    }),
    start_time: PropTypes.string,
  })),
  onConfirmWithCourses: PropTypes.func,
};

export default DeleteConfirmDialog; 