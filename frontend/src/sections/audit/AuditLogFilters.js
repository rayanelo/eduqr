import PropTypes from 'prop-types';
import { useState } from 'react';
import {
  Box,
  Stack,
  TextField,
  MenuItem,
  Button,
  Grid,
  Typography,
  Collapse,
} from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { fr } from 'date-fns/locale';

import Iconify from '../../components/iconify';

const ACTION_OPTIONS = [
  { value: '', label: 'Toutes les actions' },
  { value: 'create', label: 'Création' },
  { value: 'update', label: 'Modification' },
  { value: 'delete', label: 'Suppression' },
  { value: 'login', label: 'Connexion' },
  { value: 'logout', label: 'Déconnexion' },
];

const RESOURCE_TYPE_OPTIONS = [
  { value: '', label: 'Toutes les ressources' },
  { value: 'user', label: 'Utilisateur' },
  { value: 'room', label: 'Salle' },
  { value: 'subject', label: 'Matière' },
  { value: 'course', label: 'Cours' },
  { value: 'event', label: 'Événement' },
];

export default function AuditLogFilters({ filters, onFilterChange, loading }) {
  const [showAdvanced, setShowAdvanced] = useState(false);

  const handleFilterChange = (field, value) => {
    onFilterChange({
      ...filters,
      [field]: value,
    });
  };

  const handleReset = () => {
    onFilterChange({
      page: 1,
      limit: 20,
      action: '',
      resource_type: '',
      resource_id: '',
      user_id: '',
      start_date: '',
      end_date: '',
      search: '',
    });
  };

  return (
    <LocalizationProvider dateAdapter={AdapterDateFns} adapterLocale={fr}>
      <Box sx={{ p: 3 }}>
        <Stack direction="row" alignItems="center" justifyContent="space-between" sx={{ mb: 3 }}>
          <Typography variant="h6">Filtres</Typography>
          <Stack direction="row" spacing={1}>
            <Button
              variant="outlined"
              size="small"
              onClick={() => setShowAdvanced(!showAdvanced)}
              startIcon={<Iconify icon={showAdvanced ? 'eva:arrow-up-fill' : 'eva:arrow-down-fill'} />}
            >
              {showAdvanced ? 'Masquer' : 'Avancés'}
            </Button>
            <Button
              variant="outlined"
              size="small"
              onClick={handleReset}
              startIcon={<Iconify icon="eva:refresh-fill" />}
              disabled={loading}
            >
              Réinitialiser
            </Button>
          </Stack>
        </Stack>

        <Grid container spacing={2}>
          {/* Recherche textuelle */}
          <Grid item xs={12} md={6}>
            <TextField
              fullWidth
              label="Recherche"
              placeholder="Rechercher dans les descriptions, emails, rôles..."
              value={filters.search}
              onChange={(e) => handleFilterChange('search', e.target.value)}
              disabled={loading}
            />
          </Grid>

          {/* Action */}
          <Grid item xs={12} sm={6} md={3}>
            <TextField
              fullWidth
              select
              label="Action"
              value={filters.action}
              onChange={(e) => handleFilterChange('action', e.target.value)}
              disabled={loading}
            >
              {ACTION_OPTIONS.map((option) => (
                <MenuItem key={option.value} value={option.value}>
                  {option.label}
                </MenuItem>
              ))}
            </TextField>
          </Grid>

          {/* Type de ressource */}
          <Grid item xs={12} sm={6} md={3}>
            <TextField
              fullWidth
              select
              label="Type de ressource"
              value={filters.resource_type}
              onChange={(e) => handleFilterChange('resource_type', e.target.value)}
              disabled={loading}
            >
              {RESOURCE_TYPE_OPTIONS.map((option) => (
                <MenuItem key={option.value} value={option.value}>
                  {option.label}
                </MenuItem>
              ))}
            </TextField>
          </Grid>

          {/* Filtres avancés */}
          <Collapse in={showAdvanced} timeout="auto" unmountOnExit>
            <Grid item xs={12}>
              <Grid container spacing={2}>
                {/* ID de ressource */}
                <Grid item xs={12} sm={6} md={3}>
                  <TextField
                    fullWidth
                    label="ID de ressource"
                    placeholder="Ex: 123"
                    value={filters.resource_id}
                    onChange={(e) => handleFilterChange('resource_id', e.target.value)}
                    disabled={loading}
                  />
                </Grid>

                {/* ID d'utilisateur */}
                <Grid item xs={12} sm={6} md={3}>
                  <TextField
                    fullWidth
                    label="ID d'utilisateur"
                    placeholder="Ex: 456"
                    value={filters.user_id}
                    onChange={(e) => handleFilterChange('user_id', e.target.value)}
                    disabled={loading}
                  />
                </Grid>

                {/* Date de début */}
                <Grid item xs={12} sm={6} md={3}>
                  <DatePicker
                    label="Date de début"
                    value={filters.start_date ? new Date(filters.start_date) : null}
                    onChange={(date) => {
                      const dateStr = date ? date.toISOString().split('T')[0] : '';
                      handleFilterChange('start_date', dateStr);
                    }}
                    renderInput={(params) => <TextField {...params} fullWidth disabled={loading} />}
                    format="dd/MM/yyyy"
                  />
                </Grid>

                {/* Date de fin */}
                <Grid item xs={12} sm={6} md={3}>
                  <DatePicker
                    label="Date de fin"
                    value={filters.end_date ? new Date(filters.end_date) : null}
                    onChange={(date) => {
                      const dateStr = date ? date.toISOString().split('T')[0] : '';
                      handleFilterChange('end_date', dateStr);
                    }}
                    renderInput={(params) => <TextField {...params} fullWidth disabled={loading} />}
                    format="dd/MM/yyyy"
                  />
                </Grid>
              </Grid>
            </Grid>
          </Collapse>
        </Grid>
      </Box>
    </LocalizationProvider>
  );
}

AuditLogFilters.propTypes = {
  filters: PropTypes.object.isRequired,
  onFilterChange: PropTypes.func.isRequired,
  loading: PropTypes.bool.isRequired,
}; 