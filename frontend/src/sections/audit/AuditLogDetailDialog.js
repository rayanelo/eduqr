import PropTypes from 'prop-types';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  Box,
  Stack,
  Chip,
  Divider,
  Grid,
  Paper,
} from '@mui/material';
import { format } from 'date-fns';
import { fr } from 'date-fns/locale';

import Iconify from '../../components/iconify';

export default function AuditLogDetailDialog({ open, log, onClose }) {
  if (!log) return null;

  const getActionColor = (action) => {
    switch (action) {
      case 'create':
        return 'success';
      case 'update':
        return 'info';
      case 'delete':
        return 'error';
      case 'login':
        return 'primary';
      case 'logout':
        return 'warning';
      default:
        return 'default';
    }
  };

  const getActionLabel = (action) => {
    switch (action) {
      case 'create':
        return 'Création';
      case 'update':
        return 'Modification';
      case 'delete':
        return 'Suppression';
      case 'login':
        return 'Connexion';
      case 'logout':
        return 'Déconnexion';
      default:
        return action;
    }
  };

  const getResourceTypeLabel = (resourceType) => {
    switch (resourceType) {
      case 'user':
        return 'Utilisateur';
      case 'room':
        return 'Salle';
      case 'subject':
        return 'Matière';
      case 'course':
        return 'Cours';
      case 'event':
        return 'Événement';
      default:
        return resourceType;
    }
  };

  const getRoleColor = (role) => {
    switch (role) {
      case 'super_admin':
        return 'error';
      case 'admin':
        return 'warning';
      case 'professeur':
        return 'info';
      case 'etudiant':
        return 'success';
      default:
        return 'default';
    }
  };

  const getRoleLabel = (role) => {
    switch (role) {
      case 'super_admin':
        return 'Super Administrateur';
      case 'admin':
        return 'Administrateur';
      case 'professeur':
        return 'Professeur';
      case 'etudiant':
        return 'Étudiant';
      default:
        return role;
    }
  };

  const formatJSON = (jsonString) => {
    if (!jsonString) return null;
    try {
      const parsed = JSON.parse(jsonString);
      return JSON.stringify(parsed, null, 2);
    } catch (error) {
      return jsonString;
    }
  };

  const renderJSONData = (title, data, color = 'primary') => {
    if (!data) return null;

    const formattedData = formatJSON(data);
    if (!formattedData) return null;

    return (
      <Box sx={{ mt: 2 }}>
        <Typography variant="subtitle2" color={`${color}.main`} sx={{ mb: 1 }}>
          {title}
        </Typography>
        <Paper
          variant="outlined"
          sx={{
            p: 2,
            bgcolor: 'background.neutral',
            fontFamily: 'monospace',
            fontSize: '0.875rem',
            maxHeight: 200,
            overflow: 'auto',
          }}
        >
          <pre style={{ margin: 0, whiteSpace: 'pre-wrap' }}>{formattedData}</pre>
        </Paper>
      </Box>
    );
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>
        <Stack direction="row" alignItems="center" spacing={2}>
          <Iconify icon="eva:eye-fill" color="primary.main" />
          <Typography variant="h6">Détails du Log d'Audit</Typography>
        </Stack>
      </DialogTitle>

      <DialogContent>
        <Stack spacing={3}>
          {/* Informations générales */}
          <Box>
            <Typography variant="h6" sx={{ mb: 2 }}>
              Informations générales
            </Typography>
            <Grid container spacing={2}>
              <Grid item xs={12} sm={6}>
                <Stack spacing={1}>
                  <Typography variant="body2" color="text.secondary">
                    ID du log
                  </Typography>
                  <Typography variant="body1" fontWeight="medium">
                    #{log.id}
                  </Typography>
                </Stack>
              </Grid>
              <Grid item xs={12} sm={6}>
                <Stack spacing={1}>
                  <Typography variant="body2" color="text.secondary">
                    Date et heure
                  </Typography>
                  <Typography variant="body1" fontWeight="medium">
                    {format(new Date(log.created_at), 'dd/MM/yyyy HH:mm:ss', { locale: fr })}
                  </Typography>
                </Stack>
              </Grid>
            </Grid>
          </Box>

          <Divider />

          {/* Utilisateur */}
          <Box>
            <Typography variant="h6" sx={{ mb: 2 }}>
              Utilisateur
            </Typography>
            <Grid container spacing={2}>
              <Grid item xs={12} sm={6}>
                <Stack spacing={1}>
                  <Typography variant="body2" color="text.secondary">
                    Email
                  </Typography>
                  <Typography variant="body1" fontWeight="medium">
                    {log.user_email}
                  </Typography>
                </Stack>
              </Grid>
              <Grid item xs={12} sm={6}>
                <Stack spacing={1}>
                  <Typography variant="body2" color="text.secondary">
                    Rôle
                  </Typography>
                  <Chip
                    label={getRoleLabel(log.user_role)}
                    color={getRoleColor(log.user_role)}
                    size="small"
                    variant="filled"
                  />
                </Stack>
              </Grid>
            </Grid>
          </Box>

          <Divider />

          {/* Action */}
          <Box>
            <Typography variant="h6" sx={{ mb: 2 }}>
              Action
            </Typography>
            <Grid container spacing={2}>
              <Grid item xs={12} sm={6}>
                <Stack spacing={1}>
                  <Typography variant="body2" color="text.secondary">
                    Type d'action
                  </Typography>
                  <Chip
                    label={getActionLabel(log.action)}
                    color={getActionColor(log.action)}
                    size="small"
                    variant="filled"
                  />
                </Stack>
              </Grid>
              <Grid item xs={12} sm={6}>
                <Stack spacing={1}>
                  <Typography variant="body2" color="text.secondary">
                    Ressource
                  </Typography>
                  <Typography variant="body1" fontWeight="medium">
                    {getResourceTypeLabel(log.resource_type)}
                    {log.resource_id && ` (ID: ${log.resource_id})`}
                  </Typography>
                </Stack>
              </Grid>
            </Grid>
          </Box>

          <Divider />

          {/* Description */}
          <Box>
            <Typography variant="h6" sx={{ mb: 2 }}>
              Description
            </Typography>
            <Paper variant="outlined" sx={{ p: 2, bgcolor: 'background.neutral' }}>
              <Typography variant="body1">{log.description}</Typography>
            </Paper>
          </Box>

          {/* Anciennes et nouvelles valeurs */}
          {log.old_values && (
            <>
              <Divider />
              {renderJSONData('Anciennes valeurs', log.old_values, 'error')}
            </>
          )}

          {log.new_values && (
            <>
              <Divider />
              {renderJSONData('Nouvelles valeurs', log.new_values, 'success')}
            </>
          )}

          <Divider />

          {/* Informations techniques */}
          <Box>
            <Typography variant="h6" sx={{ mb: 2 }}>
              Informations techniques
            </Typography>
            <Grid container spacing={2}>
              <Grid item xs={12} sm={6}>
                <Stack spacing={1}>
                  <Typography variant="body2" color="text.secondary">
                    Adresse IP
                  </Typography>
                  <Typography variant="body1" fontWeight="medium">
                    {log.ip_address}
                  </Typography>
                </Stack>
              </Grid>
              <Grid item xs={12} sm={6}>
                <Stack spacing={1}>
                  <Typography variant="body2" color="text.secondary">
                    User Agent
                  </Typography>
                  <Typography variant="body1" sx={{ wordBreak: 'break-all' }}>
                    {log.user_agent || 'Non spécifié'}
                  </Typography>
                </Stack>
              </Grid>
            </Grid>
          </Box>
        </Stack>
      </DialogContent>

      <DialogActions>
        <Button onClick={onClose} variant="outlined">
          Fermer
        </Button>
      </DialogActions>
    </Dialog>
  );
}

AuditLogDetailDialog.propTypes = {
  open: PropTypes.bool.isRequired,
  log: PropTypes.object,
  onClose: PropTypes.func.isRequired,
}; 