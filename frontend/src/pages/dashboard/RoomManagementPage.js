import { useState, useCallback, useEffect } from 'react';
import { Helmet } from 'react-helmet-async';
// @mui
import {
  Card,
  Stack,
  Button,
  Container,
  TextField,
  InputAdornment,
  Chip,
  IconButton,
  Tooltip,
  Typography,
} from '@mui/material';
// hooks
import { useRooms } from '../../hooks/useRooms';
import { usePermissions } from '../../hooks/usePermissions';
import { useDeletion } from '../../hooks/useDeletion';
// routes
import { PATH_DASHBOARD } from '../../routes/paths';
// components
import Iconify from '../../components/iconify';
import { useSnackbar } from '../../components/snackbar';
import CustomBreadcrumbs from '../../components/custom-breadcrumbs';
import { useSettingsContext } from '../../components/settings';
import { DataTable } from '../../components/data-table';
import DeleteConfirmDialog from '../../components/confirm-dialog/DeleteConfirmDialog';

// ----------------------------------------------------------------------

// Fonction pour obtenir la couleur du statut modulaire
const getModularColor = (isModular) => {
  return isModular ? 'info' : 'default';
};

// Fonction pour traduire le statut modulaire
const translateModular = (isModular) => {
  return isModular ? 'Modulable' : 'Simple';
};

// Colonnes du tableau avec rendu personnalisé
const getTableColumns = (permissions) => [
  {
    id: 'name',
    label: 'Nom de la salle',
    align: 'left',
    minWidth: 200,
    render: (value) => (
      <Typography variant="body2" sx={{ fontWeight: 'medium' }}>
        {value}
      </Typography>
    ),
  },
  {
    id: 'building',
    label: 'Bâtiment',
    align: 'left',
    minWidth: 150,
  },
  {
    id: 'floor',
    label: 'Étage',
    align: 'left',
    minWidth: 100,
  },
  {
    id: 'is_modular',
    label: 'Type',
    align: 'center',
    minWidth: 120,
    render: (value) => (
      <Chip
        label={translateModular(value)}
        color={getModularColor(value)}
        size="small"
      />
    ),
  },
  {
    id: 'parent_name',
    label: 'Salle parente',
    align: 'left',
    minWidth: 150,
    render: (value) => value || '-',
  },
  {
    id: 'children_count',
    label: 'Sous-salles',
    align: 'center',
    minWidth: 100,
    render: (value) => value > 0 ? value : '-',
  },
  {
    id: 'actions',
    label: 'Actions',
    align: 'center',
    minWidth: 120,
    render: (value, row) => {
      return (
        <Stack direction="row" spacing={1} justifyContent="center">
          {permissions.canView && (
            <Tooltip title="Voir">
              <IconButton
                size="small"
                color="info"
                onClick={() => row.onView(row.id)}
              >
                <Iconify icon="eva:eye-fill" />
              </IconButton>
            </Tooltip>
          )}
          
          {permissions.canEdit && (
            <Tooltip title="Modifier">
              <IconButton
                size="small"
                color="primary"
                onClick={() => row.onEdit(row.id)}
              >
                <Iconify icon="eva:edit-fill" />
              </IconButton>
            </Tooltip>
          )}
          
          {permissions.canDelete && (
            <Tooltip title="Supprimer">
              <IconButton
                size="small"
                color="error"
                onClick={() => row.onDelete(row.id)}
              >
                <Iconify icon="eva:trash-2-fill" />
              </IconButton>
            </Tooltip>
          )}
        </Stack>
      );
    },
  },
];

// ----------------------------------------------------------------------

export default function RoomManagementPage() {
  const themeStretch = useSettingsContext();
  const { enqueueSnackbar } = useSnackbar();
  const { rooms, loading, error, fetchRooms } = useRooms();
  const { canManageRooms, canDeleteRoom } = usePermissions();
  const { deleteRoom, isDeleting } = useDeletion();

  const [filterName, setFilterName] = useState('');
  const [openConfirm, setOpenConfirm] = useState(false);
  const [selectedRoom, setSelectedRoom] = useState(null);
  const [deleteConflicts, setDeleteConflicts] = useState([]);

  useEffect(() => {
    fetchRooms();
  }, [fetchRooms]);

  // Gestion de la suppression sécurisée
  const handleOpenConfirm = useCallback((room) => {
    // Vérifier les permissions de suppression
    if (!canDeleteRoom()) {
      enqueueSnackbar('Vous n\'avez pas les permissions pour supprimer des salles', { variant: 'error' });
      return;
    }
    
    setSelectedRoom(room);
    setOpenConfirm(true);
  }, [canDeleteRoom, enqueueSnackbar]);

  const handleCloseConfirm = () => {
    setOpenConfirm(false);
    setSelectedRoom(null);
    setDeleteConflicts([]);
  };

  const handleConfirmDelete = async () => {
    if (!selectedRoom) return;

    try {
      const result = await deleteRoom(selectedRoom.id);
      
      if (result.success) {
        // Rafraîchir la liste des salles
        await fetchRooms();
        handleCloseConfirm();
      } else {
        // Afficher les conflits empêchant la suppression
        if (result.data?.future_courses) {
          setDeleteConflicts(result.data.future_courses.map(course => ({
            courseName: course.name,
            date: new Date(course.start_time).toLocaleDateString('fr-FR'),
            description: `Cours prévu le ${new Date(course.start_time).toLocaleDateString('fr-FR')} à ${new Date(course.start_time).toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' })}`
          })));
        }
      }
    } catch (error) {
      enqueueSnackbar(error.message || 'Erreur lors de la suppression', { variant: 'error' });
    }
  };

  // Gestion des actions du tableau
  const handleEditRow = useCallback((id) => {
    // TODO: Implémenter la modification de salle
    console.log('Edit room:', id);
  }, []);

  const handleViewRow = useCallback((id) => {
    // TODO: Implémenter la vue détaillée de salle
    console.log('View room:', id);
  }, []);

  const handleDeleteRow = useCallback((id) => {
    const room = rooms.find((r) => r.id === id);
    handleOpenConfirm(room);
  }, [rooms, handleOpenConfirm]);

  // Filtrage des données
  const filteredRooms = rooms.filter((room) =>
    room.name.toLowerCase().includes(filterName.toLowerCase()) ||
    room.building.toLowerCase().includes(filterName.toLowerCase())
  );

  // Préparation des données pour le tableau
  const tableData = filteredRooms.map((room) => ({
    id: room.id,
    name: room.name,
    building: room.building,
    floor: room.floor,
    is_modular: room.is_modular,
    parent_name: room.parent?.name || null,
    children_count: room.children?.length || 0,
    // Ajout des fonctions d'action pour chaque ligne
    onEdit: handleEditRow,
    onView: handleViewRow,
    onDelete: handleDeleteRow,
  }));

  // Obtenir les colonnes selon les permissions
  const tableColumns = getTableColumns({ 
    canView: canManageRooms,
    canEdit: canManageRooms,
    canDelete: canDeleteRoom,
  });

  return (
    <>
      <Helmet>
        <title> Gestion des salles | EduQR</title>
      </Helmet>

      <Container maxWidth={themeStretch ? false : 'xl'}>
        <Stack spacing={3}>
          {/* Section header */}
          <CustomBreadcrumbs
            heading="Gestion des salles"
            links={[
              {
                name: 'Dashboard',
                href: PATH_DASHBOARD.root,
              },
              {
                name: 'Salles',
              },
            ]}
            action={
              canManageRooms && (
                <Button
                  variant="contained"
                  startIcon={<Iconify icon="eva:plus-fill" />}
                  onClick={() => console.log('Create room')}
                >
                  Nouvelle salle
                </Button>
              )
            }
          />

          {/* Filtres */}
          <Card sx={{ p: 2 }}>
            <TextField
              fullWidth
              placeholder="Rechercher une salle..."
              value={filterName}
              onChange={(e) => setFilterName(e.target.value)}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <Iconify icon="eva:search-fill" sx={{ color: 'text.disabled' }} />
                  </InputAdornment>
                ),
              }}
            />
          </Card>

          {/* Tableau des salles */}
          <DataTable
            title="Liste des salles"
            data={tableData}
            columns={tableColumns}
            loading={loading}
            error={error}
            onRefresh={fetchRooms}
          />
        </Stack>

        {/* Dialog de confirmation de suppression */}
        <DeleteConfirmDialog
          open={openConfirm}
          onClose={handleCloseConfirm}
          onConfirm={handleConfirmDelete}
          title="Supprimer la salle"
          message={`Êtes-vous sûr de vouloir supprimer la salle "${selectedRoom?.name}" ?`}
          resourceName={selectedRoom?.name}
          resourceType="room"
          isDeleting={isDeleting}
          conflicts={deleteConflicts}
        />
      </Container>
    </>
  );
} 