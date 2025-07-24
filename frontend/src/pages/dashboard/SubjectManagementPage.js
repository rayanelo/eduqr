import React, { useState, useEffect } from 'react';
import {
  Container,
  Card,
  Stack,
  Button,
  Typography,
  IconButton,
  Tooltip,
  Alert,
  Box,
  TextField,
  InputAdornment,
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import { useSubjects } from '../../hooks/useSubjects';
import { usePermissions } from '../../hooks/usePermissions';
import { useContext } from 'react';
import { AuthContext } from '../../auth/JwtContext';
import DataTable from '../../components/data-table/DataTable';
import ConfirmDialog from '../../components/confirm-dialog/ConfirmDialog';
import SubjectFormDialog from '../../sections/subjects/SubjectFormDialog';
import Iconify from '../../components/iconify';

export default function SubjectManagementPage() {
  const TABLE_HEAD = [
    { id: 'name', label: 'Nom de la matière', align: 'left', minWidth: 200 },
    { id: 'code', label: 'Code', align: 'left', minWidth: 120 },
    { id: 'description', label: 'Description', align: 'left', minWidth: 300 },
    { id: 'actions', label: 'Actions', align: 'center', width: 120 },
  ];
  
  const navigate = useNavigate();
  const { user } = useContext(AuthContext);
  const { canManageSubjects } = usePermissions();
  const {
    subjects,
    loading,
    error,
    fetchSubjects,
    deleteSubject,
    setError,
  } = useSubjects();

  const [openFormDialog, setOpenFormDialog] = useState(false);
  const [openDeleteDialog, setOpenDeleteDialog] = useState(false);
  const [selectedSubject, setSelectedSubject] = useState(null);
  const [filterName, setFilterName] = useState('');

  useEffect(() => {
    console.log('SubjectManagementPage - User:', user);
    console.log('SubjectManagementPage - canManageSubjects:', canManageSubjects);
    fetchSubjects();
  }, [fetchSubjects, user, canManageSubjects]);

  // Vérification de sécurité
  if (!user) {
    return (
      <Container maxWidth="xl">
        <Alert severity="warning" sx={{ mt: 3 }}>
          Chargement de l'utilisateur...
        </Alert>
      </Container>
    );
  }

  if (canManageSubjects === undefined) {
    return (
      <Container maxWidth="xl">
        <Alert severity="warning" sx={{ mt: 3 }}>
          Chargement des permissions...
        </Alert>
      </Container>
    );
  }

  if (!canManageSubjects) {
    return (
      <Container maxWidth="xl">
        <Alert severity="error" sx={{ mt: 3 }}>
          Vous n'avez pas les permissions nécessaires pour accéder à cette page.
        </Alert>
      </Container>
    );
  }

  const handleOpenFormDialog = (subject = null) => {
    setSelectedSubject(subject);
    setOpenFormDialog(true);
  };

  const handleCloseFormDialog = () => {
    setSelectedSubject(null);
    setOpenFormDialog(false);
  };

  const handleOpenDeleteDialog = (subject) => {
    setSelectedSubject(subject);
    setOpenDeleteDialog(true);
  };

  const handleCloseDeleteDialog = () => {
    setSelectedSubject(null);
    setOpenDeleteDialog(false);
  };

  const handleDeleteSubject = async () => {
    if (selectedSubject) {
      try {
        await deleteSubject(selectedSubject.id);
        handleCloseDeleteDialog();
      } catch (error) {
        // L'erreur est gérée dans le hook
      }
    }
  };

  const filteredSubjects = subjects.filter((subject) =>
    subject.name.toLowerCase().includes(filterName.toLowerCase()) ||
    (subject.code && subject.code.toLowerCase().includes(filterName.toLowerCase()))
  );

  const dataFiltered = filteredSubjects.map((subject) => ({
    id: subject.id,
    name: subject.name,
    code: subject.code || '-',
    description: subject.description || '-',
    actions: (
      <Stack direction="row" spacing={1}>
        <Tooltip title="Modifier">
          <IconButton
            onClick={() => handleOpenFormDialog(subject)}
            color="primary"
            size="small"
          >
            <EditIcon />
          </IconButton>
        </Tooltip>
        <Tooltip title="Supprimer">
          <IconButton
            onClick={() => handleOpenDeleteDialog(subject)}
            color="error"
            size="small"
          >
            <DeleteIcon />
          </IconButton>
        </Tooltip>
      </Stack>
    ),
  }));


  return (
    <Container maxWidth="xl">
      <Stack spacing={3}>
        <Stack direction="row" alignItems="center" justifyContent="space-between">
          <Typography variant="h4">Gestion des Matières</Typography>
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => handleOpenFormDialog()}
          >
            Nouvelle Matière
          </Button>
        </Stack>

        {error && (
          <Alert severity="error" onClose={() => setError(null)}>
            {error}
          </Alert>
        )}

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
                placeholder="Rechercher par nom ou code..."
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
            data={dataFiltered}
            columns={TABLE_HEAD}
            tableData={dataFiltered}
            onAddNew={() => handleOpenFormDialog()}
            isFiltered={!!filterName}
          />
        </Card>
      </Stack>

      {/* Dialog de formulaire */}
      <SubjectFormDialog
        open={openFormDialog}
        onClose={handleCloseFormDialog}
        subject={selectedSubject}
        onSuccess={() => {
          handleCloseFormDialog();
          fetchSubjects();
        }}
      />

      {/* Dialog de confirmation de suppression */}
      <ConfirmDialog
        open={openDeleteDialog}
        onClose={handleCloseDeleteDialog}
        onConfirm={handleDeleteSubject}
        title="Supprimer la matière"
        content={`Êtes-vous sûr de vouloir supprimer la matière "${selectedSubject?.name}" ? Cette action est irréversible.`}
        confirmText="Supprimer"
        cancelText="Annuler"
        confirmColor="error"
      />
    </Container>
  );
} 