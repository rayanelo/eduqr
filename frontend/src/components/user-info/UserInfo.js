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
    <Card sx={{ p: 1.5, mb: 1.5 }}>
      <Stack direction="row" alignItems="center" spacing={1}>
        <Avatar
          src={user.avatar || '/assets/images/avatars/default-avatar.png'}
          alt={user.first_name}
          sx={{ width: 40, height: 40 }}
        />
        
        <Box sx={{ flexGrow: 1 }}>
          <Typography variant="subtitle2" sx={{ mb: 0.25, lineHeight: 1.2 }}>
            {user.first_name} {user.last_name}
          </Typography>
          
          <Typography variant="caption" color="text.secondary" sx={{ mb: 0.25, display: 'block' }}>
            {user.email}
          </Typography>
          
          <Stack direction="row" spacing={0.5} alignItems="center">
            <Chip
              label={getRoleLabel(currentUserRole)}
              color={getRoleColor(currentUserRole)}
              size="small"
              sx={{ height: 20, fontSize: '0.7rem' }}
            />
            
            {hasRoleOrHigher(ROLES.ADMIN) && (
              <Typography variant="caption" color="text.secondary" sx={{ fontSize: '0.7rem' }}>
                • Peut gérer les utilisateurs
              </Typography>
            )}
          </Stack>
        </Box>
      </Stack>

      {creatableRoles.length > 0 && (
        <Box sx={{ mt: 1 }}>
          <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mb: 0.25, fontSize: '0.7rem' }}>
            Rôles que vous pouvez créer :
          </Typography>
          <Stack direction="row" spacing={0.25} flexWrap="wrap">
            {creatableRoles.map((role) => (
              <Chip
                key={role}
                label={getRoleLabel(role)}
                color={getRoleColor(role)}
                size="small"
                variant="outlined"
                sx={{ fontSize: '0.65rem', height: 18 }}
              />
            ))}
          </Stack>
        </Box>
      )}
    </Card>
  );
} 