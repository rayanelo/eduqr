import PropTypes from 'prop-types';
// @mui
import {
  Box,
  Card,
  Stack,
  Avatar,
  Typography,
  Chip,
} from '@mui/material';
// components
import { useAuthContext } from '../../auth/JwtContext';
import { usePermissions, ROLES } from '../../hooks/usePermissions';

// ----------------------------------------------------------------------

export default function UserInfo() {
  const { user } = useAuthContext();
  const { currentUserRole, getCreatableRoles, hasRoleOrHigher } = usePermissions();

  if (!user) {
    return null;
  }

  const getRoleColor = (role) => {
    switch (role) {
      case ROLES.SUPER_ADMIN:
        return 'error';
      case ROLES.ADMIN:
        return 'warning';
      case ROLES.PROFESSEUR:
        return 'info';
      case ROLES.ETUDIANT:
        return 'default';
      default:
        return 'default';
    }
  };

  const getRoleLabel = (role) => {
    switch (role) {
      case ROLES.SUPER_ADMIN:
        return 'Super Admin';
      case ROLES.ADMIN:
        return 'Admin';
      case ROLES.PROFESSEUR:
        return 'Professeur';
      case ROLES.ETUDIANT:
        return 'Étudiant';
      default:
        return role;
    }
  };

  const creatableRoles = getCreatableRoles();

  return (
    <Card sx={{ p: 2, mb: 2 }}>
      <Stack direction="row" alignItems="center" spacing={1.5}>
        <Avatar
          src={user.avatar || '/assets/images/avatars/default-avatar.png'}
          alt={user.first_name}
          sx={{ width: 48, height: 48 }}
        />
        
        <Box sx={{ flexGrow: 1 }}>
          <Typography variant="subtitle1" sx={{ mb: 0.5 }}>
            {user.first_name} {user.last_name}
          </Typography>
          
          <Typography variant="body2" color="text.secondary" sx={{ mb: 0.5 }}>
            {user.email}
          </Typography>
          
          <Stack direction="row" spacing={1} alignItems="center">
            <Chip
              label={getRoleLabel(currentUserRole)}
              color={getRoleColor(currentUserRole)}
              size="small"
            />
            
            {hasRoleOrHigher(ROLES.ADMIN) && (
              <Typography variant="caption" color="text.secondary">
                • Peut gérer les utilisateurs
              </Typography>
            )}
          </Stack>
        </Box>
      </Stack>

      {creatableRoles.length > 0 && (
        <Box sx={{ mt: 1.5 }}>
          <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mb: 0.5 }}>
            Rôles que vous pouvez créer :
          </Typography>
          <Stack direction="row" spacing={0.5} flexWrap="wrap">
            {creatableRoles.map((role) => (
              <Chip
                key={role}
                label={getRoleLabel(role)}
                color={getRoleColor(role)}
                size="small"
                variant="outlined"
                sx={{ fontSize: '0.75rem' }}
              />
            ))}
          </Stack>
        </Box>
      )}
    </Card>
  );
} 