import PropTypes from 'prop-types';
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
  IconButton,
  Tooltip,
  Typography,
  Box,
  TablePagination,
  Skeleton,
  Alert,
} from '@mui/material';
import { format } from 'date-fns';
import { fr } from 'date-fns/locale';

import Iconify from '../../components/iconify';

export default function AuditLogTable({
  logs,
  loading,
  pagination,
  onPageChange,
  onViewDetail,
  getActionColor,
  getActionLabel,
  getResourceTypeLabel,
}) {
  const handleChangePage = (event, newPage) => {
    onPageChange(newPage + 1); // API uses 1-based pagination
  };

  const handleChangeRowsPerPage = (event) => {
    // This would need to be implemented in the parent component
    console.log('Change rows per page:', event.target.value);
  };

  if (loading && (!logs || logs.length === 0)) {
    return (
      <Box sx={{ p: 3 }}>
        <Skeleton variant="rectangular" height={400} />
      </Box>
    );
  }

  if (!loading && (!logs || logs.length === 0)) {
    return (
      <Box sx={{ p: 3, textAlign: 'center' }}>
        <Alert severity="info">
          Aucun log d'audit trouvé pour les critères sélectionnés.
        </Alert>
      </Box>
    );
  }

  return (
    <Box>
      <TableContainer component={Paper} sx={{ minHeight: 400 }}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Date/Heure</TableCell>
              <TableCell>Utilisateur</TableCell>
              <TableCell>Action</TableCell>
              <TableCell>Ressource</TableCell>
              <TableCell>Description</TableCell>
              <TableCell>Adresse IP</TableCell>
              <TableCell align="center">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {logs.map((log) => (
              <TableRow key={log.id} hover>
                <TableCell>
                  <Typography variant="body2" color="text.secondary">
                    {format(new Date(log.created_at), 'dd/MM/yyyy HH:mm:ss', { locale: fr })}
                  </Typography>
                </TableCell>
                <TableCell>
                  <Box>
                    <Typography variant="body2" fontWeight="medium">
                      {log.user_email}
                    </Typography>
                    <Chip
                      label={log.user_role}
                      size="small"
                      color="default"
                      variant="outlined"
                    />
                  </Box>
                </TableCell>
                <TableCell>
                  <Chip
                    label={getActionLabel(log.action)}
                    color={getActionColor(log.action)}
                    size="small"
                    variant="filled"
                  />
                </TableCell>
                <TableCell>
                  <Box>
                    <Typography variant="body2" fontWeight="medium">
                      {getResourceTypeLabel(log.resource_type)}
                    </Typography>
                    {log.resource_id && (
                      <Typography variant="caption" color="text.secondary">
                        ID: {log.resource_id}
                      </Typography>
                    )}
                  </Box>
                </TableCell>
                <TableCell>
                  <Typography variant="body2" sx={{ maxWidth: 300 }}>
                    {log.description}
                  </Typography>
                </TableCell>
                <TableCell>
                  <Typography variant="body2" color="text.secondary">
                    {log.ip_address}
                  </Typography>
                </TableCell>
                <TableCell align="center">
                  <Tooltip title="Voir les détails">
                    <IconButton
                      size="small"
                      onClick={() => onViewDetail(log.id)}
                      color="primary"
                    >
                      <Iconify icon="eva:eye-fill" />
                    </IconButton>
                  </Tooltip>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      <TablePagination
        rowsPerPageOptions={[10, 20, 50, 100]}
        component="div"
        count={pagination.total}
        rowsPerPage={pagination.limit}
        page={pagination.page - 1} // MUI uses 0-based pagination
        onPageChange={handleChangePage}
        onRowsPerPageChange={handleChangeRowsPerPage}
        labelRowsPerPage="Lignes par page:"
        labelDisplayedRows={({ from, to, count }) =>
          `${from}-${to} sur ${count !== -1 ? count : `plus de ${to}`}`
        }
      />
    </Box>
  );
}

AuditLogTable.propTypes = {
  logs: PropTypes.array.isRequired,
  loading: PropTypes.bool.isRequired,
  pagination: PropTypes.object.isRequired,
  onPageChange: PropTypes.func.isRequired,
  onViewDetail: PropTypes.func.isRequired,
  getActionColor: PropTypes.func.isRequired,
  getActionLabel: PropTypes.func.isRequired,
  getResourceTypeLabel: PropTypes.func.isRequired,
}; 