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
  Grid,
  Alert,
  Avatar,
  Chip,
  IconButton,
  Tooltip,
  Box,
  Typography,
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
import { DataTable } from '../../components/data-table';
import UserForm from '../../components/user-form/UserForm';
import { ConfirmDialog } from '../../components/confirm-dialog';
import UserInfo from '../../components/user-info/UserInfo';

// ----------------------------------------------------------------------

// Fonction pour formater la date
const formatDate = (dateString) => {
  if (!dateString) return '-';
  try {
    const date = new Date(dateString);
    return date.toLocaleDateString('fr-FR', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
    });
  } catch (error) {
    return '-';
  }
};

// Fonction pour obtenir la couleur du rôle
const getRoleColor = (role) => {
  const roleColors = {
    'super_admin': 'error',
    'admin': 'warning',
    'professeur': 'info',
    'etudiant': 'success',
  };
  return roleColors[role] || 'default';
};

// Fonction pour traduire le rôle
const translateRole = (role) => {
  const roleTranslations = {
    'super_admin': 'Super Admin',
    'admin': 'Admin',
    'professeur': 'Professeur',
    'etudiant': 'Étudiant',
  };
  return roleTranslations[role] || role;
};

// Colonnes du tableau avec rendu personnalisé
const getTableColumns = (permissions) => [
  {
    id: 'avatar',
    label: 'Avatar',
    align: 'center',
    width: 80,
    render: (value) => (
      <Avatar
        src={value}
        alt="Avatar"
        sx={{ width: 40, height: 40 }}
      />
    ),
  },
  {
    id: 'name',
    label: 'Nom complet',
    align: 'left',
    minWidth: 200,
    render: (value) => (
      <Typography variant="body2" sx={{ fontWeight: 'medium' }}>
        {value}
      </Typography>
    ),
  },
  {
    id: 'email',
    label: 'Email',
    align: 'left',
    minWidth: 200,
    render: (value) => (
      <Typography variant="body2" sx={{ color: 'text.secondary' }}>
        {value}
      </Typography>
    ),
  },
  {
    id: 'phone',
    label: 'Téléphone',
    align: 'left',
    minWidth: 150,
    render: (value) => (
      <Typography variant="body2" sx={{ color: 'text.secondary' }}>
        {value}
      </Typography>
    ),
  },
  {
    id: 'address',
    label: 'Adresse',
    align: 'left',
    minWidth: 200,
    render: (value) => (
      <Typography variant="body2" sx={{ color: 'text.secondary' }}>
        {value}
      </Typography>
    ),
  },
  {
    id: 'role',
    label: 'Rôle',
    align: 'center',
    width: 120,
    render: (value) => (
      <Chip
        label={translateRole(value)}
        color={getRoleColor(value)}
        size="small"
        variant="soft"
      />
    ),
  },
  {
    id: 'created_at',
    label: 'Créé le',
    align: 'center',
    width: 120,
    render: (value) => (
      <Typography variant="body2" sx={{ color: 'text.secondary' }}>
        {formatDate(value)}
      </Typography>
    ),
  },
  ...(permissions.canSeeActions ? [{
    id: 'actions',
    label: 'Actions',
    align: 'center',
    width: 120,
    render: (value, row) => (
      <Box sx={{ display: 'flex', gap: 0.5, justifyContent: 'center' }}>
        <Tooltip title="Voir">
          <IconButton
            size="small"
            color="info"
            onClick={() => row.onView(row.id)}
          >
            <Iconify icon="eva:eye-fill" />
          </IconButton>
        </Tooltip>
        <Tooltip title="Modifier">
          <IconButton
            size="small"
            color="primary"
            onClick={() => row.onEdit(row.id)}
          >
            <Iconify icon="eva:edit-fill" />
          </IconButton>
        </Tooltip>
        <Tooltip title="Supprimer">
          <IconButton
            size="small"
            color="error"
            onClick={() => row.onDelete(row.id)}
          >
            <Iconify icon="eva:trash-2-fill" />
          </IconButton>
        </Tooltip>
      </Box>
    ),
  }] : []),
];

// ----------------------------------------------------------------------

export default function UserManagementPage() {
  const { themeStretch } = useSettingsContext();
  const { enqueueSnackbar } = useSnackbar();
  const { getCreatableRoles, canSeeActions } = usePermissions();

  const {
    users,
    loading: isLoading,
    error,
    setError,
    fetchUsers,
    createUser,
    updateUser,
    deleteUser,
  } = useUsers();

  const [openForm, setOpenForm] = useState(false);
  const [openConfirm, setOpenConfirm] = useState(false);
  const [selectedUser, setSelectedUser] = useState(null);
  const [isEdit, setIsEdit] = useState(false);
  const [filterName, setFilterName] = useState('');

  // Charger les utilisateurs au montage du composant
  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

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
      handleCloseForm();
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
    // Ajout des fonctions d'action pour chaque ligne
    onEdit: handleEditRow,
    onView: handleViewRow,
    onDelete: handleDeleteRow,
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
        <Grid container spacing={3} sx={{ mb: 2 }} justifyContent="space-between">
          <Grid item xs={12} md={5}>
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
          
          <Grid item xs={12} md={5}>
            <UserInfo />
          </Grid>
        </Grid>

        {/* Affichage des erreurs */}
        {error && (
          <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
            {error}
          </Alert>
        )}

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