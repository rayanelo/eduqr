// @mui
import { styled, alpha } from '@mui/material/styles';
import { Box, Typography } from '@mui/material';
// auth
import { useAuthContext } from '../../../auth/useAuthContext';
// components
import { CustomAvatar } from '../../../components/custom-avatar';

// ----------------------------------------------------------------------

const StyledRoot = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  padding: theme.spacing(2, 2.5),
  borderRadius: Number(theme.shape.borderRadius) * 1.5,
  backgroundColor: alpha(theme.palette.grey[500], 0.12),
}));

// ----------------------------------------------------------------------

export default function NavAccount() {
  const { user } = useAuthContext();

  const fullName = `${user?.first_name || ''} ${user?.last_name || ''}`.trim();

  return (
    <StyledRoot>
      <CustomAvatar 
        src={user?.photoURL} 
        alt={fullName} 
        firstName={user?.first_name}
        lastName={user?.last_name}
        role={user?.role}
      />

      <Box sx={{ ml: 2, minWidth: 0 }}>
        <Typography variant="subtitle2" noWrap>
          {fullName}
        </Typography>

        <Typography variant="body2" noWrap sx={{ color: 'text.secondary' }}>
          {user?.role}
        </Typography>
      </Box>
    </StyledRoot>
  );
}
