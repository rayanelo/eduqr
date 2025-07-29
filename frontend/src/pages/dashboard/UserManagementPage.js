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
  Chip,
  IconButton,
  Tooltip,
  Typography,
} from '@mui/material';
// hooks
import { useUsers } from '../../hooks/useUsers';
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
import UserForm from '../../components/user-form/UserForm';
import DeleteConfirmDialog from '../../components/confirm-dialog/DeleteConfirmDialog';
import UserInfo from '../../components/user-info/UserInfo';
import { CustomAvatar } from '../../components/custom-avatar';

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
    render: (value, row) => (
      <CustomAvatar
        firstName={row.first_name}
        lastName={row.last_name}
        role={row.role}
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
  },
  {
    id: 'phone',
    label: 'Téléphone',
    align: 'left',
    minWidth: 150,
  },
  {
    id: 'role',
    label: 'Rôle',
    align: 'center',
    minWidth: 120,
    render: (value) => (
      <Chip
        label={translateRole(value)}
        color={getRoleColor(value)}
        size="small"
      />
    ),
  },
  {
    id: 'created_at',
    label: 'Date de création',
    align: 'center',
    minWidth: 150,
    render: (value) => formatDate(value),
  },
  {
    id: 'actions',
    label: 'Actions',
    align: 'center',
    minWidth: 120,
    render: (value, row) => {
      const user = { id: row.id, role: row.role };
      
      return (
        <Stack direction="row" spacing={1} justifyContent="center">
          {permissions.canView && permissions.canView(user) && (
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
          
          {permissions.canEdit && permissions.canEdit(user) && (
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
          
          {permissions.canDelete && permissions.canDelete(user) && (
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

export default function UserManagementPage() {
  const themeStretch = useSettingsContext();
  const { enqueueSnackbar } = useSnackbar();
  const { users, loading, error, fetchUsers, createUser, updateUser } = useUsers();
  const { canManageUsers, canDeleteUser, canEditUser, canViewUser } = usePermissions();
  const { deleteUser, isDeleting } = useDeletion();

  const [filterName, setFilterName] = useState('');
  const [openForm, setOpenForm] = useState(false);
  const [openConfirm, setOpenConfirm] = useState(false);
  const [selectedUser, setSelectedUser] = useState(null);
  const [deleteConflicts, setDeleteConflicts] = useState([]);
  const [futureCourses, setFutureCourses] = useState([]);
  const [pastCourses, setPastCourses] = useState([]);

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  const handleOpenForm = (user = null) => {
    setSelectedUser(user);
    setOpenForm(true);
  };

  const handleCloseForm = () => {
    setOpenForm(false);
    setSelectedUser(null);
  };

  const handleSubmitForm = async (data) => {
    try {
      if (selectedUser) {
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

  // Gestion de la suppression sécurisée
  const handleOpenConfirm = useCallback((user) => {
    // Vérifier les permissions de suppression
    if (!canDeleteUser(user)) {
      enqueueSnackbar('Vous n\'avez pas les permissions pour supprimer cet utilisateur', { variant: 'error' });
      return;
    }
    
    setSelectedUser(user);
    setOpenConfirm(true);
  }, [canDeleteUser, enqueueSnackbar]);

  const handleCloseConfirm = () => {
    setOpenConfirm(false);
    setSelectedUser(null);
    setDeleteConflicts([]);
    setFutureCourses([]);
    setPastCourses([]);
  };

  const handleConfirmDelete = async () => {
    if (!selectedUser) return;

    try {
      const result = await deleteUser(selectedUser.id);
      
      if (result.success) {
        // Rafraîchir la liste des utilisateurs
        await fetchUsers();
        handleCloseConfirm();
      } else if (result.hasLinkedCourses) {
        // Afficher les cours liés dans le dialog
        setDeleteConflicts([]); // Réinitialiser les conflits
        setFutureCourses(result.data.future_courses || []);
        setPastCourses(result.data.past_courses || []);
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

  const handleConfirmDeleteWithCourses = async () => {
    if (!selectedUser) return;

    try {
      const result = await deleteUser(selectedUser.id, true); // true = confirmer avec les cours
      
      if (result.success) {
        // Rafraîchir la liste des utilisateurs
        await fetchUsers();
        handleCloseConfirm();
      }
    } catch (error) {
      enqueueSnackbar(error.message || 'Erreur lors de la suppression', { variant: 'error' });
    }
  };

  // Gestion des actions du tableau
  const handleEditRow = useCallback((id) => {
    const user = users.find((u) => u.id === id);
    if (canEditUser(user)) {
      handleOpenForm(user);
    } else {
      enqueueSnackbar('Vous n\'avez pas les permissions pour modifier cet utilisateur', { variant: 'error' });
    }
  }, [users, canEditUser, enqueueSnackbar]);

  const handleViewRow = useCallback((id) => {
    const user = users.find((u) => u.id === id);
    if (canViewUser(user)) {
      setSelectedUser(user);
      // TODO: Implémenter l'affichage des détails utilisateur
      enqueueSnackbar('Fonctionnalité en cours de développement', { variant: 'info' });
    } else {
      enqueueSnackbar('Vous n\'avez pas les permissions pour voir cet utilisateur', { variant: 'error' });
    }
  }, [users, canViewUser, enqueueSnackbar]);

  const handleDeleteRow = useCallback((id) => {
    const user = users.find((u) => u.id === id);
    handleOpenConfirm(user);
  }, [users, handleOpenConfirm]);

  // Filtrage des données selon les permissions
  const filteredUsers = users.filter((user) => {
    // Vérifier si l'utilisateur peut voir cet utilisateur
    if (!canViewUser(user)) return false;
    
    // Appliquer le filtre de recherche
    return `${user.first_name} ${user.last_name}`.toLowerCase().includes(filterName.toLowerCase()) ||
           user.email.toLowerCase().includes(filterName.toLowerCase());
  });

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
  const tableColumns = getTableColumns({ 
    canSeeActions: canManageUsers,
    canEdit: canEditUser,
    canDelete: canDeleteUser,
    canView: canViewUser,
  });

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
                canManageUsers && (
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
          <Alert severity="error" sx={{ mb: 2 }}>
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
            isLoading={loading}
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
        isEdit={!!selectedUser}
        isLoading={loading}
      />

      {/* Dialog de confirmation de suppression */}
      <DeleteConfirmDialog
        open={openConfirm}
        onClose={handleCloseConfirm}
        onConfirm={handleConfirmDelete}
        onConfirmWithCourses={handleConfirmDeleteWithCourses}
        title="Supprimer l'utilisateur"
        message={`Êtes-vous sûr de vouloir supprimer ${selectedUser?.first_name} ${selectedUser?.last_name} ?`}
        resourceName={`${selectedUser?.first_name} ${selectedUser?.last_name}`}
        resourceType="user"
        isDeleting={isDeleting}
        conflicts={deleteConflicts}
        futureCourses={futureCourses}
        pastCourses={pastCourses}
      />


    </>
  );
} 