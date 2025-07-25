import PropTypes from 'prop-types';
import orderBy from 'lodash/orderBy';
// @mui
import { DatePicker } from '../../../components/date-picker';
import {
  Box,
  Stack,
  Drawer,
  Divider,
  Tooltip,
  Typography,
  IconButton,
  ListItemText,
  ListItemButton,
} from '@mui/material';
// utils
import { fDateTime } from '../../../utils/formatTime';
// components
import Iconify from '../../../components/iconify';
import Scrollbar from '../../../components/scrollbar';
// Suppression de l'import ColorMultiPicker

// ----------------------------------------------------------------------

CalendarFilterDrawer.propTypes = {
  events: PropTypes.array,
  picker: PropTypes.object,
  openFilter: PropTypes.bool,
  onCloseFilter: PropTypes.func,
  onResetFilter: PropTypes.func,
  onSelectEvent: PropTypes.func,
};

export default function CalendarFilterDrawer({
  events,
  picker,
  openFilter,
  onCloseFilter,
  onResetFilter,
  onSelectEvent,
}) {
  const notDefault = (picker.startDate && picker.endDate);

  return (
    <Drawer
      anchor="right"
      open={openFilter}
      onClose={onCloseFilter}
      BackdropProps={{
        invisible: true,
      }}
      PaperProps={{
        sx: { width: 320 },
      }}
    >
      <Stack
        direction="row"
        alignItems="center"
        justifyContent="space-between"
        sx={{ pl: 2, pr: 1, py: 2 }}
      >
        <Typography variant="subtitle1">Filtres</Typography>

        <Tooltip title="Réinitialiser">
          <Box sx={{ position: 'relative' }}>
            <IconButton onClick={onResetFilter}>
              <Iconify icon="ic:round-refresh" />
            </IconButton>

            {notDefault && (
              <Box
                sx={{
                  top: 6,
                  right: 4,
                  width: 8,
                  height: 8,
                  borderRadius: '50%',
                  position: 'absolute',
                  bgcolor: 'error.main',
                }}
              />
            )}
          </Box>
        </Tooltip>
      </Stack>

      <Divider />

      {/* Suppression de la section des couleurs */}

      <Typography
        variant="caption"
        sx={{
          p: 2,
          color: 'text.secondary',
          fontWeight: 'fontWeightMedium',
        }}
      >
        Période
      </Typography>

      <Stack spacing={2} sx={{ px: 2 }}>
        <DatePicker
          label="Date de début"
          value={picker.startDate}
          onChange={picker.onChangeStartDate}
          size="small"
        />

        <DatePicker
          label="Date de fin"
          value={picker.endDate}
          onChange={picker.onChangeEndDate}
          size="small"
          error={picker.isError}
          helperText={picker.isError && 'La date de fin doit être postérieure à la date de début'}
        />
      </Stack>

      <Typography
        variant="caption"
        sx={{
          p: 2,
          color: 'text.secondary',
          fontWeight: 'fontWeightMedium',
        }}
      >
        Cours ({events.length})
      </Typography>

      <Scrollbar sx={{ height: 1 }}>
        {orderBy(events, ['end'], ['desc']).map((event) => (
          <ListItemButton
            key={event.id}
            onClick={() => onSelectEvent(event.id)}
            sx={{ py: 1.5, borderBottom: (theme) => `dashed 1px ${theme.palette.divider}` }}
          >
            <Box
              sx={{
                top: 16,
                left: 0,
                width: 0,
                height: 0,
                position: 'absolute',
                borderRight: '10px solid transparent',
                borderTop: `10px solid ${event.color}`,
              }}
            />

            <ListItemText
              disableTypography
              primary={
                <Typography variant="subtitle2" sx={{ fontSize: 13, mt: 0.5 }}>
                  {event.title}
                </Typography>
              }
              secondary={
                <Typography
                  variant="caption"
                  component="div"
                  sx={{ fontSize: 11, color: 'text.disabled' }}
                >
                  {event.allDay ? (
                    fDateTime(event.start, 'dd MMM yy')
                  ) : (
                    <>
                      {`${fDateTime(event.start, 'dd MMM yy p')} - ${fDateTime(
                        event.end,
                        'dd MMM yy p'
                      )}`}
                    </>
                  )}
                </Typography>
              }
              sx={{ display: 'flex', flexDirection: 'column-reverse' }}
            />
          </ListItemButton>
        ))}
      </Scrollbar>
    </Drawer>
  );
} 