import PropTypes from 'prop-types';
// @mui
import { Box, Typography } from '@mui/material';

// ----------------------------------------------------------------------

export default function CustomBreadcrumbs({ links, action, sx, ...other }) {
  return (
    <Box sx={{ mb: 3, ...sx }}>
      <Typography variant="h4" sx={{ mb: 1 }}>
        {links[links.length - 1]?.name || 'Page'}
      </Typography>
      {action && (
        <Box sx={{ flexShrink: 0 }}>{action}</Box>
      )}
    </Box>
  );
}

CustomBreadcrumbs.propTypes = {
  action: PropTypes.node,
  links: PropTypes.array,
  sx: PropTypes.object,
}; 