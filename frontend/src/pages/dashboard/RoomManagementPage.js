import { useState, useCallback, useEffect } from 'react';
import { useSnackbar } from 'notistack';
import {
  Container,
  Card,
  Stack,
  Button,
  TextField,
  Typography,
  Chip,
  IconButton,
  Tooltip,
  Alert,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Box,
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Search as SearchIcon,
  Clear as ClearIcon,
  ExpandMore as ExpandMoreIcon,
} from '@mui/icons-material';

import { useSettingsContext } from '../../components/settings';
import { useRooms } from '../../hooks/useRooms';
import { usePermissions } from '../../hooks/usePermissions';
import { PATH_DASHBOARD } from '../../routes/paths';
import CustomBreadcrumbs from '../../components/custom-breadcrumbs/CustomBreadcrumbs';
import RoomFormDialog from '../../sections/rooms/RoomFormDialog';
import ConfirmDialog from '../../components/confirm-dialog/ConfirmDialog';

export default function RoomManagementPage() {
  const { enqueueSnackbar } = useSnackbar();
  const settings = useSettingsContext();
  const { canManageRooms } = usePermissions();

  // Room management
  const {
    rooms,
    isLoading,
    error,
    getAllRooms,
    createRoom,
    updateRoom,
    deleteRoom,
  } = useRooms();

  // Local state
  const [filterName, setFilterName] = useState('');
  const [filterBuilding, setFilterBuilding] = useState('');
  const [filterFloor, setFilterFloor] = useState('');
  const [openForm, setOpenForm] = useState(false);
  const [selectedRoom, setSelectedRoom] = useState(null);
  const [openDeleteDialog, setOpenDeleteDialog] = useState(false);
  const [roomToDelete, setRoomToDelete] = useState(null);

  // Load rooms on component mount
  const loadRooms = useCallback(() => {
    getAllRooms();
  }, [getAllRooms]);

  useEffect(() => {
    loadRooms();
  }, [loadRooms]);



  // Filter handlers
  const handleFilterName = (event) => {
    setFilterName(event.target.value);
  };

  const handleFilterBuilding = (event) => {
    setFilterBuilding(event.target.value);
  };

  const handleFilterFloor = (event) => {
    setFilterFloor(event.target.value);
  };

  const handleSearch = () => {
    getAllRooms({
      name: filterName,
      building: filterBuilding,
      floor: filterFloor,
    });
  };

  const handleResetFilters = () => {
    setFilterName('');
    setFilterBuilding('');
    setFilterFloor('');
    getAllRooms();
  };

  // Form handlers
  const handleOpenForm = (room = null) => {
    setSelectedRoom(room);
    setOpenForm(true);
  };

  const handleCloseForm = () => {
    setOpenForm(false);
    setSelectedRoom(null);
  };

  const handleFormSuccess = () => {
    handleCloseForm();
    loadRooms();
    enqueueSnackbar(
      selectedRoom ? 'Salle modifiée avec succès' : 'Salle créée avec succès',
      { variant: 'success' }
    );
  };

  // Delete handlers
  const handleDeleteClick = (room) => {
    setRoomToDelete(room);
    setOpenDeleteDialog(true);
  };

  const handleConfirmDelete = async () => {
    if (!roomToDelete) return;

    try {
      await deleteRoom(roomToDelete.id);
      setOpenDeleteDialog(false);
      setRoomToDelete(null);
      loadRooms();
      enqueueSnackbar('Salle supprimée avec succès', { variant: 'success' });
    } catch (error) {
      enqueueSnackbar('Erreur lors de la suppression', { variant: 'error' });
    }
  };

  // Filter rooms based on search criteria
  const filteredRooms = rooms.filter((room) => {
    const matchesName = !filterName || room.name.toLowerCase().includes(filterName.toLowerCase());
    const matchesBuilding = !filterBuilding || (room.building && room.building.toLowerCase().includes(filterBuilding.toLowerCase()));
    const matchesFloor = !filterFloor || (room.floor && room.floor.toLowerCase().includes(filterFloor.toLowerCase()));
    return matchesName && matchesBuilding && matchesFloor;
  });

  // Get only main rooms (no parent)
  const mainRooms = filteredRooms.filter(room => !room.parent_id);



  const getSubRooms = (parentId) => {
    return filteredRooms.filter(room => room.parent_id === parentId);
  };

  if (isLoading) {
    return (
      <Container maxWidth={settings.themeStretch ? false : 'xl'}>
        <Typography>Chargement...</Typography>
      </Container>
    );
  }

  return (
    <Container maxWidth={settings.themeStretch ? false : 'xl'}>
      <CustomBreadcrumbs
        heading="Gestion des salles"
        links={[
          { name: 'Dashboard', href: PATH_DASHBOARD.root },
          { name: 'Salles' },
        ]}
        action={
          canManageRooms && (
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              onClick={() => handleOpenForm()}
            >
              Nouvelle salle
            </Button>
          )
        }
      />

      <Card sx={{ mt: 3 }}>
        <Stack spacing={3} sx={{ p: 3 }}>
          {/* Search and filters */}
          <Card variant="outlined" sx={{ p: 2, mb: 2 }}>
            <Stack spacing={2}>
              <Typography variant="subtitle2" color="text.secondary">
                Filtres de recherche
              </Typography>
              <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
                <TextField
                  size="small"
                  label="Nom de la salle"
                  placeholder="Rechercher par nom..."
                  value={filterName}
                  onChange={handleFilterName}
                  onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                  sx={{ minWidth: 200 }}
                />
                <TextField
                  size="small"
                  label="Bâtiment"
                  placeholder="Filtrer par bâtiment..."
                  value={filterBuilding}
                  onChange={handleFilterBuilding}
                  onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                  sx={{ minWidth: 150 }}
                />
                <TextField
                  size="small"
                  label="Étage"
                  placeholder="Filtrer par étage..."
                  value={filterFloor}
                  onChange={handleFilterFloor}
                  onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                  sx={{ minWidth: 120 }}
                />
                <Stack direction="row" spacing={1}>
                  <Button
                    variant="contained"
                    size="small"
                    startIcon={<SearchIcon />}
                    onClick={handleSearch}
                  >
                    Rechercher
                  </Button>
                  <Button
                    variant="outlined"
                    size="small"
                    startIcon={<ClearIcon />}
                    onClick={handleResetFilters}
                  >
                    Effacer
                  </Button>
                </Stack>
              </Stack>
            </Stack>
          </Card>

          {/* Error message */}
          {error && (
            <Alert severity="error">
              {error}
            </Alert>
          )}

          {/* Rooms list with cards */}
          <Stack spacing={2}>
            {mainRooms.length === 0 ? (
              <Card sx={{ p: 4, textAlign: 'center' }}>
                <Typography variant="h6" color="text.secondary" gutterBottom>
                  Aucune salle trouvée
                </Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                  Il n'y a pas encore de salles dans cette liste. Commencez par ajouter une nouvelle salle.
                </Typography>
                {canManageRooms && (
                  <Button
                    variant="contained"
                    startIcon={<AddIcon />}
                    onClick={() => handleOpenForm()}
                  >
                    Ajouter la première salle
                  </Button>
                )}
              </Card>
            ) : (
              <Stack spacing={1}>
                {mainRooms.map((room) => (
                  <Accordion key={room.id} variant="outlined" sx={{ '&:before': { display: 'none' } }}>
                    <AccordionSummary 
                      expandIcon={null}
                      sx={{ 
                        '& .MuiAccordionSummary-content': { 
                          margin: 0,
                          alignItems: 'center'
                        }
                      }}
                    >
                      <Stack direction="row" alignItems="center" spacing={2} sx={{ flexGrow: 1 }}>
                        {room.is_modular && room.children?.length > 0 && (
                          <ExpandMoreIcon sx={{ color: 'text.secondary', fontSize: 20 }} />
                        )}
                        <Tooltip title={room.name} placement="top">
                          <Typography 
                            variant="subtitle1" 
                            sx={{ 
                              fontWeight: 'bold', 
                              width: 180,
                              overflow: 'hidden',
                              textOverflow: 'ellipsis',
                              whiteSpace: 'nowrap'
                            }}
                          >
                            {room.name}
                          </Typography>
                        </Tooltip>
                        
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, flexGrow: 1 }}>
                          <Chip 
                            label={room.is_modular ? "Principale modulable" : "Principale"} 
                            color={room.is_modular ? "success" : "default"}
                            size="small"
                            sx={{ width: 140 }}
                          />
                          
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, width: 160 }}>
                            <Typography variant="body2" color="text.secondary" sx={{ fontSize: '0.75rem' }}>
                              Bâtiment:
                            </Typography>
                            <Typography variant="body2" sx={{ fontWeight: 'medium', fontSize: '0.75rem' }}>
                              {room.building || 'Non spécifié'}
                            </Typography>
                          </Box>
                          
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, width: 120 }}>
                            <Typography variant="body2" color="text.secondary" sx={{ fontSize: '0.75rem' }}>
                              Étage:
                            </Typography>
                            <Typography variant="body2" sx={{ fontWeight: 'medium', fontSize: '0.75rem' }}>
                              {room.floor || 'Non spécifié'}
                            </Typography>
                          </Box>
                          
                          {room.is_modular && (
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, width: 100 }}>
                              <Typography variant="body2" color="text.secondary" sx={{ fontSize: '0.75rem' }}>
                                Sous-salles:
                              </Typography>
                              <Typography variant="body2" sx={{ fontWeight: 'medium', fontSize: '0.75rem' }}>
                                {room.children?.length || 0}
                              </Typography>
                            </Box>
                          )}
                        </Box>
                      </Stack>
                      
                      <Stack direction="row" spacing={1} sx={{ ml: 2 }}>
                        <Tooltip title="Modifier">
                          <IconButton
                            size="small"
                            onClick={(e) => {
                              e.stopPropagation();
                              handleOpenForm(room);
                            }}
                            disabled={!canManageRooms}
                          >
                            <EditIcon />
                          </IconButton>
                        </Tooltip>
                        <Tooltip title="Supprimer">
                          <IconButton
                            size="small"
                            color="error"
                            onClick={(e) => {
                              e.stopPropagation();
                              handleDeleteClick(room);
                            }}
                            disabled={!canManageRooms}
                          >
                            <DeleteIcon />
                          </IconButton>
                        </Tooltip>
                      </Stack>
                    </AccordionSummary>
                    
                    {/* Sub-rooms section */}
                    {room.is_modular && room.children && room.children.length > 0 && (
                      <AccordionDetails sx={{ pt: 0, pb: 2 }}>
                        <Stack spacing={1}>
                          {room.children.map((subRoom) => (
                            <Card key={subRoom.id} variant="outlined" sx={{ bgcolor: 'grey.50' }}>
                              <Stack
                                direction="row"
                                alignItems="center"
                                spacing={2}
                                sx={{ p: 1.5 }}
                              >
                                <Stack direction="row" alignItems="center" spacing={2} sx={{ flexGrow: 1 }}>
                                  <Tooltip title={subRoom.name} placement="top">
                                    <Typography 
                                      variant="body2" 
                                      sx={{ 
                                        fontWeight: 'medium', 
                                        width: 120,
                                        overflow: 'hidden',
                                        textOverflow: 'ellipsis',
                                        whiteSpace: 'nowrap'
                                      }}
                                    >
                                      {subRoom.name}
                                    </Typography>
                                  </Tooltip>
                                  
                                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, flexGrow: 1 }}>
                                    <Chip 
                                      label="Sous-salle" 
                                      color="info"
                                      size="small"
                                      sx={{ width: 100 }}
                                    />
                                    
                                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, width: 160 }}>
                                      <Typography variant="body2" color="text.secondary" sx={{ fontSize: '0.75rem' }}>
                                        Bâtiment:
                                      </Typography>
                                      <Typography variant="body2" sx={{ fontWeight: 'medium', fontSize: '0.75rem' }}>
                                        {subRoom.building || 'Non spécifié'}
                                      </Typography>
                                    </Box>
                                    
                                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, width: 120 }}>
                                      <Typography variant="body2" color="text.secondary" sx={{ fontSize: '0.75rem' }}>
                                        Étage:
                                      </Typography>
                                      <Typography variant="body2" sx={{ fontWeight: 'medium', fontSize: '0.75rem' }}>
                                        {subRoom.floor || 'Non spécifié'}
                                      </Typography>
                                    </Box>
                                  </Box>
                                </Stack>
                                
                                <Stack direction="row" spacing={0.5}>
                                  <Tooltip title="Modifier">
                                    <IconButton
                                      size="small"
                                      onClick={() => handleOpenForm(subRoom)}
                                      disabled={!canManageRooms}
                                    >
                                      <EditIcon fontSize="small" />
                                    </IconButton>
                                  </Tooltip>
                                  <Tooltip title="Supprimer">
                                    <IconButton
                                      size="small"
                                      color="error"
                                      onClick={() => handleDeleteClick(subRoom)}
                                      disabled={!canManageRooms}
                                    >
                                      <DeleteIcon fontSize="small" />
                                    </IconButton>
                                  </Tooltip>
                                </Stack>
                              </Stack>
                            </Card>
                          ))}
                        </Stack>
                      </AccordionDetails>
                    )}
                  </Accordion>
                ))}
              </Stack>
            )}
          </Stack>
        </Stack>
      </Card>

      {/* Room form dialog */}
      <RoomFormDialog
        open={openForm}
        onClose={handleCloseForm}
        onSuccess={handleFormSuccess}
        room={selectedRoom}
      />

      {/* Delete confirmation dialog */}
      <ConfirmDialog
        open={openDeleteDialog}
        onClose={() => setOpenDeleteDialog(false)}
        onConfirm={handleConfirmDelete}
        title="Confirmer la suppression"
        content={`Êtes-vous sûr de vouloir supprimer la salle "${roomToDelete?.name}" ? Cette action est irréversible.`}
        confirmText="Supprimer"
        cancelText="Annuler"
      />
    </Container>
  );
} 