import { useState, useCallback } from 'react';
import { Helmet } from 'react-helmet-async';
// @mui
import {
  Card,
  Stack,
  Button,
  Container,
  TextField,
  InputAdornment,
  Grid,
  Box,
} from '@mui/material';
// hooks
import { useUsers } from '../../hooks/useUsers';
import { usePermissions } from '../../hooks/usePermissions';
// routes
import { PATH_DASHBOARD } from '../../routes/paths';
// components
import Iconify from '../../components/iconify';
import { useSnackbar } from '../../components/snackbar';
import CustomBreadcrumbs from '../../components/custom-breadcrumbs';
import { useSettingsContext } from '../../components/settings';
import DataTable from '../../components/data-table/DataTable';
import UserForm from '../../components/user-form/UserForm';
import ConfirmDialog from '../../components/confirm-dialog/ConfirmDialog';
import UserInfo from '../../components/user-info/UserInfo';

// ----------------------------------------------------------------------

// Colonnes de base du tableau
const BASE_TABLE_HEAD = [
  { id: 'avatar', label: 'Avatar', align: 'center', width: 80, type: 'avatar' },
  { id: 'name', label: 'Nom complet', align: 'left', minWidth: 200 },
  { id: 'email', label: 'Email', align: 'left', minWidth: 200 },
  { id: 'phone', label: 'Téléphone', align: 'left', minWidth: 150 },
  { id: 'address', label: 'Adresse', align: 'left', minWidth: 200 },
  { id: 'role', label: 'Rôle', align: 'center', width: 120, type: 'role' },
  { id: 'created_at', label: 'Créé le', align: 'center', width: 120, type: 'date' },
  { id: 'actions', label: 'Actions', align: 'center', width: 80 },
];

// Fonction pour obtenir les colonnes selon les permissions
const getTableColumns = (permissions) => {
  const columns = [...BASE_TABLE_HEAD];
  
  // Si l'utilisateur ne peut pas voir les actions, on retire la colonne actions
  if (!permissions.canSeeActions) {
    return columns.filter(col => col.id !== 'actions');
  }
  
  return columns;
};

// ----------------------------------------------------------------------

export default function UserManagementPage() {
  const { themeStretch } = useSettingsContext();
  const { enqueueSnackbar } = useSnackbar();
  const { canCreateUser, getCreatableRoles, canSeeActions } = usePermissions();

  const {
    users,
    isLoading,
    createUser,
    updateUser,
    deleteUser,
    updateUserRole,
  } = useUsers();

  const [openForm, setOpenForm] = useState(false);
  const [openConfirm, setOpenConfirm] = useState(false);
  const [selectedUser, setSelectedUser] = useState(null);
  const [isEdit, setIsEdit] = useState(false);
  const [filterName, setFilterName] = useState('');

  // Gestion du formulaire
  const handleOpenForm = (user = null) => {
    setSelectedUser(user);
    setIsEdit(!!user);
    setOpenForm(true);
  };

  const handleCloseForm = () => {
    setOpenForm(false);
    setSelectedUser(null);
    setIsEdit(false);
  };

  const handleSubmitForm = async (data) => {
    try {
      if (isEdit) {
        await updateUser(selectedUser.id, data);
        enqueueSnackbar('Utilisateur modifié avec succès!');
      } else {
        await createUser(data);
        enqueueSnackbar('Utilisateur créé avec succès!');
      }
    } catch (error) {
      enqueueSnackbar(error.message || 'Une erreur est survenue', { variant: 'error' });
    }
  };

  // Gestion de la suppression
  const handleOpenConfirm = (user) => {
    setSelectedUser(user);
    setOpenConfirm(true);
  };

  const handleCloseConfirm = () => {
    setOpenConfirm(false);
    setSelectedUser(null);
  };

  const handleConfirmDelete = async () => {
    try {
      await deleteUser(selectedUser.id);
      enqueueSnackbar('Utilisateur supprimé avec succès!');
      handleCloseConfirm();
    } catch (error) {
      enqueueSnackbar(error.message || 'Erreur lors de la suppression', { variant: 'error' });
    }
  };

  // Gestion des actions du tableau
  const handleEditRow = useCallback((id) => {
    const user = users.find((u) => u.id === id);
    handleOpenForm(user);
  }, [users]);

  const handleViewRow = useCallback((id) => {
    const user = users.find((u) => u.id === id);
    // Ici vous pouvez implémenter la vue détaillée
    console.log('View user:', user);
  }, [users]);

  const handleDeleteRow = useCallback((id) => {
    const user = users.find((u) => u.id === id);
    handleOpenConfirm(user);
  }, [users]);

  const handleUpdateRole = useCallback(async (id, role) => {
    try {
      await updateUserRole(id, role);
      enqueueSnackbar(`Rôle mis à jour avec succès!`);
    } catch (error) {
      enqueueSnackbar(error.message || 'Erreur lors de la mise à jour du rôle', { variant: 'error' });
    }
  }, [updateUserRole, enqueueSnackbar]);

  // Filtrage des données
  const filteredUsers = users.filter((user) =>
    `${user.first_name} ${user.last_name}`.toLowerCase().includes(filterName.toLowerCase()) ||
    user.email.toLowerCase().includes(filterName.toLowerCase())
  );

  // Préparation des données pour le tableau
  const tableData = filteredUsers.map((user) => ({
    id: user.id,
    avatar: user.avatar || '/assets/images/avatars/default-avatar.png',
    name: `${user.first_name} ${user.last_name}`,
    email: user.email,
    phone: user.phone || '-',
    address: user.address || '-',
    role: user.role,
    created_at: user.created_at,
  }));

  // Obtenir les colonnes selon les permissions
  const tableColumns = getTableColumns({ canSeeActions });

  return (
    <>
      <Helmet>
        <title> Gestion des utilisateurs | EduQR</title>
      </Helmet>

      <Container maxWidth={themeStretch ? false : 'xl'}>
        {/* Section header avec titre et UserInfo */}
        <Grid container spacing={3} sx={{ mb: 3 }}>
          <Grid item xs={12} md={9}>
            <CustomBreadcrumbs
              heading="Gestion des utilisateurs"
              links={[
                {
                  name: 'Dashboard',
                  href: PATH_DASHBOARD.root,
                },
                {
                  name: 'Utilisateurs',
                },
              ]}
              action={
                getCreatableRoles().length > 0 && (
                  <Button
                    variant="contained"
                    startIcon={<Iconify icon="eva:plus-fill" />}
                    onClick={() => handleOpenForm()}
                  >
                    Nouvel utilisateur
                  </Button>
                )
              }
            />
          </Grid>
          
          <Grid item xs={12} md={3}>
            <UserInfo />
          </Grid>
        </Grid>

        {/* Tableau en pleine largeur */}
        <Card>
          <Stack
            spacing={2.5}
            direction={{ xs: 'column', md: 'row' }}
            alignItems={{ xs: 'flex-end', md: 'center' }}
            justifyContent="space-between"
            sx={{
              p: 2.5,
              pr: { xs: 2.5, md: 1 },
            }}
          >
            <Stack direction="row" alignItems="center" spacing={2} flexGrow={1}>
              <TextField
                fullWidth
                value={filterName}
                onChange={(event) => setFilterName(event.target.value)}
                placeholder="Rechercher un utilisateur..."
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <Iconify icon="eva:search-fill" sx={{ color: 'text.disabled' }} />
                    </InputAdornment>
                  ),
                }}
              />
            </Stack>
          </Stack>

          <DataTable
            data={tableData}
            columns={tableColumns}
            tableData={tableData}
            onEditRow={handleEditRow}
            onViewRow={handleViewRow}
            onDeleteRow={handleDeleteRow}
            onUpdateRole={handleUpdateRole}
            onAddNew={() => handleOpenForm()}
            isLoading={isLoading}
            isFiltered={!!filterName}
          />
        </Card>
      </Container>

      {/* Formulaire utilisateur */}
      <UserForm
        open={openForm}
        onClose={handleCloseForm}
        onSubmit={handleSubmitForm}
        user={selectedUser}
        isEdit={isEdit}
        isLoading={isLoading}
      />

      {/* Dialog de confirmation */}
      <ConfirmDialog
        open={openConfirm}
        onClose={handleCloseConfirm}
        title="Supprimer l'utilisateur"
        content={`Êtes-vous sûr de vouloir supprimer ${selectedUser?.first_name} ${selectedUser?.last_name} ?`}
        action={
          <Button variant="contained" color="error" onClick={handleConfirmDelete}>
            Supprimer
          </Button>
        }
      />
    </>
  );
} 