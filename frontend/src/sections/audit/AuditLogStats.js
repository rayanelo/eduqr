import PropTypes from 'prop-types';
import {
  Box,
  Card,
  Grid,
  Typography,
  Stack,
  Chip,
  LinearProgress,
} from '@mui/material';

import Iconify from '../../components/iconify';

export default function AuditLogStats({ stats }) {
  if (!stats) return null;

  const getActionIcon = (action) => {
    switch (action) {
      case 'create':
        return 'eva:plus-circle-fill';
      case 'update':
        return 'eva:edit-fill';
      case 'delete':
        return 'eva:trash-2-fill';
      case 'login':
        return 'eva:log-in-fill';
      case 'logout':
        return 'eva:log-out-fill';
      default:
        return 'eva:activity-fill';
    }
  };

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

  const getResourceIcon = (resourceType) => {
    switch (resourceType) {
      case 'user':
        return 'eva:person-fill';
      case 'room':
        return 'eva:home-fill';
      case 'subject':
        return 'eva:book-fill';
      case 'course':
        return 'eva:calendar-fill';
      case 'event':
        return 'eva:clock-fill';
      default:
        return 'eva:cube-fill';
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
        return 'Super Admin';
      case 'admin':
        return 'Admin';
      case 'professeur':
        return 'Professeur';
      case 'etudiant':
        return 'Étudiant';
      default:
        return role;
    }
  };

  const totalLogs = stats.total_logs || 0;
  const byAction = stats.by_action || [];
  const byResource = stats.by_resource || [];
  const byRole = stats.by_role || [];

  return (
    <Grid container spacing={3}>
      {/* Total des logs */}
      <Grid item xs={12} sm={6} md={3}>
        <Card sx={{ p: 3 }}>
          <Stack direction="row" alignItems="center" justifyContent="space-between">
            <Box>
              <Typography variant="h4" sx={{ mb: 1 }}>
                {totalLogs.toLocaleString()}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Total des logs
              </Typography>
            </Box>
            <Iconify icon="eva:activity-fill" width={48} color="primary.main" />
          </Stack>
        </Card>
      </Grid>

      {/* Répartition par action */}
      <Grid item xs={12} md={9}>
        <Card sx={{ p: 3 }}>
          <Typography variant="h6" sx={{ mb: 2 }}>
            Répartition par action
          </Typography>
          <Grid container spacing={2}>
            {byAction.map((item) => {
              const percentage = totalLogs > 0 ? (item.count / totalLogs) * 100 : 0;
              return (
                <Grid item xs={12} sm={6} md={4} key={item.action}>
                  <Stack spacing={1}>
                    <Stack direction="row" alignItems="center" spacing={1}>
                      <Iconify 
                        icon={getActionIcon(item.action)} 
                        color={getActionColor(item.action)}
                        width={20}
                      />
                      <Typography variant="body2" fontWeight="medium">
                        {getActionLabel(item.action)}
                      </Typography>
                      <Chip
                        label={item.count}
                        size="small"
                        color={getActionColor(item.action)}
                        variant="filled"
                      />
                    </Stack>
                    <LinearProgress
                      variant="determinate"
                      value={percentage}
                      color={getActionColor(item.action)}
                      sx={{ height: 6, borderRadius: 3 }}
                    />
                    <Typography variant="caption" color="text.secondary">
                      {percentage.toFixed(1)}%
                    </Typography>
                  </Stack>
                </Grid>
              );
            })}
          </Grid>
        </Card>
      </Grid>

      {/* Répartition par ressource */}
      <Grid item xs={12} md={6}>
        <Card sx={{ p: 3 }}>
          <Typography variant="h6" sx={{ mb: 2 }}>
            Répartition par ressource
          </Typography>
          <Stack spacing={2}>
            {byResource.map((item) => {
              const percentage = totalLogs > 0 ? (item.count / totalLogs) * 100 : 0;
              return (
                <Box key={item.resource_type}>
                  <Stack direction="row" alignItems="center" justifyContent="space-between" sx={{ mb: 1 }}>
                    <Stack direction="row" alignItems="center" spacing={1}>
                      <Iconify icon={getResourceIcon(item.resource_type)} width={16} />
                      <Typography variant="body2" fontWeight="medium">
                        {getResourceTypeLabel(item.resource_type)}
                      </Typography>
                    </Stack>
                    <Chip label={item.count} size="small" />
                  </Stack>
                  <LinearProgress
                    variant="determinate"
                    value={percentage}
                    sx={{ height: 4, borderRadius: 2 }}
                  />
                </Box>
              );
            })}
          </Stack>
        </Card>
      </Grid>

      {/* Répartition par rôle */}
      <Grid item xs={12} md={6}>
        <Card sx={{ p: 3 }}>
          <Typography variant="h6" sx={{ mb: 2 }}>
            Répartition par rôle
          </Typography>
          <Stack spacing={2}>
            {byRole.map((item) => {
              const percentage = totalLogs > 0 ? (item.count / totalLogs) * 100 : 0;
              return (
                <Box key={item.user_role}>
                  <Stack direction="row" alignItems="center" justifyContent="space-between" sx={{ mb: 1 }}>
                    <Typography variant="body2" fontWeight="medium">
                      {getRoleLabel(item.user_role)}
                    </Typography>
                    <Chip 
                      label={item.count} 
                      size="small" 
                      color={getRoleColor(item.user_role)}
                      variant="filled"
                    />
                  </Stack>
                  <LinearProgress
                    variant="determinate"
                    value={percentage}
                    color={getRoleColor(item.user_role)}
                    sx={{ height: 4, borderRadius: 2 }}
                  />
                </Box>
              );
            })}
          </Stack>
        </Card>
      </Grid>
    </Grid>
  );
}

AuditLogStats.propTypes = {
  stats: PropTypes.object.isRequired,
};

// Helper functions
function getActionLabel(action) {
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
}

function getResourceTypeLabel(resourceType) {
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
} 