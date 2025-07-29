// @mui
import { Stack, Box, Drawer } from '@mui/material';
import { useTheme } from '@mui/material/styles';
// config
import { NAV } from '../../../config-global';
// utils
import { hideScrollbarX } from '../../../utils/cssStyles';
// components
import Logo from '../../../components/logo';
import { NavSectionMini } from '../../../components/nav-section';
//
import { useNavData } from './config-navigation';
import NavToggleButton from './NavToggleButton';

// ----------------------------------------------------------------------

export default function NavMini() {
  const theme = useTheme();
  const navData = useNavData();
  const isRTL = theme.direction === 'rtl';
  
  return (
    <Box
      component="nav"
      sx={{
        flexShrink: { lg: 0 },
        width: { lg: NAV.W_DASHBOARD_MINI },
      }}
    >
      <NavToggleButton
        sx={{
          top: 22,
          ...(isRTL 
            ? { right: NAV.W_DASHBOARD_MINI - 12 } 
            : { left: NAV.W_DASHBOARD_MINI - 12 }
          ),
        }}
      />

      <Drawer
        open
        variant="permanent"
        anchor={isRTL ? 'right' : 'left'}
        PaperProps={{
          sx: {
            width: NAV.W_DASHBOARD_MINI,
            bgcolor: 'transparent',
            ...(isRTL 
              ? { borderLeft: (theme) => `dashed 1px ${theme.palette.divider}` }
              : { borderRight: (theme) => `dashed 1px ${theme.palette.divider}` }
            ),
          },
        }}
      >
        <Stack
          sx={{
            pb: 2,
            height: 1,
            width: NAV.W_DASHBOARD_MINI,
            ...hideScrollbarX,
          }}
        >
          <Logo sx={{ mx: 'auto', my: 2 }} />

          <NavSectionMini data={navData} />
        </Stack>
      </Drawer>
    </Box>
  );
}
