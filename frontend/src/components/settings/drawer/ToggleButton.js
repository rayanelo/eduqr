import PropTypes from 'prop-types';
// @mui
import { alpha, useTheme } from '@mui/material/styles';
import { Tooltip, Box } from '@mui/material';
// utils
import { bgBlur } from '../../../utils/cssStyles';
//
import { IconButtonAnimate } from '../../animate';
import Iconify from '../../iconify';
//
import BadgeDot from './BadgeDot';

// ----------------------------------------------------------------------

ToggleButton.propTypes = {
  open: PropTypes.bool,
  onToggle: PropTypes.func,
  notDefault: PropTypes.bool,
};

export default function ToggleButton({ notDefault, open, onToggle }) {
  const theme = useTheme();
  const isRTL = theme.direction === 'rtl';

  return (
    <Box
      sx={{
        p: 0.5,
        ...(isRTL ? { left: 24 } : { right: 24 }),
        bottom: 24,
        zIndex: 999,
        position: 'fixed',
        borderRadius: '50%',
        boxShadow: `-12px 12px 32px -4px ${alpha(
          theme.palette.mode === 'light' ? theme.palette.grey[600] : theme.palette.common.black,
          0.36
        )}`,
        ...bgBlur({ color: theme.palette.background.default }),
      }}
    >
      {notDefault && !open && (
        <BadgeDot
          sx={{
            top: 8,
            ...(isRTL ? { left: 10 } : { right: 10 }),
          }}
        />
      )}

      <Tooltip title="Settings">
        <IconButtonAnimate color="primary" onClick={onToggle} sx={{ p: 1.25 }}>
          <Iconify icon="solar:settings-linear" />
        </IconButtonAnimate>
      </Tooltip>
    </Box>
  );
}
