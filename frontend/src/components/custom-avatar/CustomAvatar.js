import { forwardRef } from 'react';
// @mui
import { Avatar, Badge } from '@mui/material';

// ----------------------------------------------------------------------

// Fonction pour obtenir les initiales
const getInitials = (firstName = '', lastName = '') => {
  const first = firstName?.trim().charAt(0)?.toUpperCase() || '';
  const last = lastName?.trim().charAt(0)?.toUpperCase() || '';
  
  if (!first && !last) return 'U';
  return `${first}${last}`;
};

// Couleurs douces et épurées basées sur le rôle
const getRoleColor = (role) => {
  switch (role) {
    case 'super_admin':
      return {
        background: '#fce4ec',
        color: '#c2185b',
        border: '#ffd700',
        borderWidth: '2px'
      };
    case 'admin':
      return {
        background: '#e3f2fd',
        color: '#1976d2',
        border: '#bbdefb',
        borderWidth: '1px'
      };
    case 'professeur':
      return {
        background: '#e8f5e8',
        color: '#388e3c',
        border: '#c8e6c9',
        borderWidth: '1px'
      };
    case 'etudiant':
      return {
        background: '#f3e5f5',
        color: '#7b1fa2',
        border: '#e1bee7',
        borderWidth: '1px'
      };
    default:
      return {
        background: '#fafafa',
        color: '#757575',
        border: '#eeeeee',
        borderWidth: '1px'
      };
  }
};

// ----------------------------------------------------------------------

const CustomAvatar = forwardRef(({ 
  color, 
  firstName = '', 
  lastName = '', 
  role = '',
  BadgeProps, 
  children, 
  sx, 
  ...other 
}, ref) => {
  const initials = getInitials(firstName, lastName);
  const roleColors = getRoleColor(role);

  // Styles ultra-épurés
  const minimalStyles = {
    fontWeight: 500,
    fontSize: '0.75rem',
    letterSpacing: '0.25px',
    border: `${roleColors.borderWidth} solid ${roleColors.border}`,
    transition: 'all 0.15s ease',
  };

  const renderContent = (
    <Avatar 
      ref={ref} 
      src={!initials || initials === 'U' ? '/assets/images/avatars/default-avatar.png' : undefined}
      sx={{
        ...minimalStyles,
        background: roleColors.background,
        color: roleColors.color,
        ...sx,
      }} 
      {...other}
    >
      {initials}
      {children}
    </Avatar>
  );

  return BadgeProps ? (
    <Badge
      overlap="circular"
      anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      {...BadgeProps}
    >
      {renderContent}
    </Badge>
  ) : (
    renderContent
  );
});

export default CustomAvatar;
