import { useState, useCallback, useEffect } from 'react';
import { Helmet } from 'react-helmet-async';
// @mui
import {
  Card,
  Stack,
  Button,
  Container,
  TextField,
  Grid,
  Chip,
  IconButton,
  Tooltip,
  Typography,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  CircularProgress,
} from '@mui/material';
// hooks
import { useSnackbar } from 'notistack';
// components
import Iconify from '../../components/iconify';
import apiClient from '../../utils/api';

// ----------------------------------------------------------------------

export default function AbsenceManagementPage() {
  const { enqueueSnackbar } = useSnackbar();

  const [absences, setAbsences] = useState([]);
  const [courses, setCourses] = useState([]);
  const [loading, setLoading] = useState(false);
  const [stats, setStats] = useState({
    totalAbsences: 0,
    pendingAbsences: 0,
    approvedAbsences: 0,
    rejectedAbsences: 0,
  });

  // Pagination
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [total, setTotal] = useState(0);

  // Dialog states
  const [openForm, setOpenForm] = useState(false);
  const [openDetails, setOpenDetails] = useState(false);
  const [selectedAbsence, setSelectedAbsence] = useState(null);

  // Form states
  const [formData, setFormData] = useState({
    courseId: '',
    justification: '',
    documentPath: '',
  });

  // Load absences
  const loadAbsences = useCallback(async () => {
    setLoading(true);
    try {
      const response = await apiClient.get(`/absences/my?page=${page + 1}&limit=${rowsPerPage}`);
      setAbsences(response.data.data || []);
      setTotal(response.data.total || 0);
    } catch (error) {
      console.error('Error loading absences:', error);
      enqueueSnackbar('Erreur lors du chargement des absences', { variant: 'error' });
    } finally {
      setLoading(false);
    }
  }, [page, rowsPerPage, enqueueSnackbar]);

  // Load courses
  const loadCourses = useCallback(async () => {
    try {
      const response = await apiClient.get('/admin/courses');
      // Filtrer seulement les cours passés
      const now = new Date();
      const pastCourses = response.data.filter(course => new Date(course.end_time) < now);
      setCourses(pastCourses);
    } catch (error) {
      console.error('Error loading courses:', error);
      enqueueSnackbar('Erreur lors du chargement des cours', { variant: 'error' });
    }
  }, [enqueueSnackbar]);

  // Load stats
  const loadStats = useCallback(async () => {
    try {
      const response = await apiClient.get('/absences/stats');
      setStats(response.data);
    } catch (error) {
      console.error('Error loading stats:', error);
    }
  }, []);

  // Load data on component mount
  useEffect(() => {
    loadAbsences();
    loadCourses();
    loadStats();
  }, [loadAbsences, loadCourses, loadStats]);

  // Handle form submission
  const handleSubmit = async () => {
    if (!formData.courseId || !formData.justification) {
      enqueueSnackbar('Veuillez remplir tous les champs obligatoires', { variant: 'warning' });
      return;
    }

    try {
      await apiClient.post('/absences', formData);
      enqueueSnackbar('Absence créée avec succès', { variant: 'success' });
      setOpenForm(false);
      setFormData({ courseId: '', justification: '', documentPath: '' });
      loadAbsences();
      loadStats();
    } catch (error) {
      console.error('Error creating absence:', error);
      enqueueSnackbar(error.response?.data?.error || 'Erreur lors de la création de l\'absence', { variant: 'error' });
    }
  };

  // Handle pagination
  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  // Get status color
  const getStatusColor = (status) => {
    switch (status) {
      case 'pending':
        return 'warning';
      case 'approved':
        return 'success';
      case 'rejected':
        return 'error';
      default:
        return 'default';
    }
  };

  // Get status label
  const getStatusLabel = (status) => {
    switch (status) {
      case 'pending':
        return 'En attente';
      case 'approved':
        return 'Approuvée';
      case 'rejected':
        return 'Rejetée';
      default:
        return status;
    }
  };

  return (
    <>
      <Helmet>
        <title>Gestion des Absences | EduQR</title>
      </Helmet>

      <Container maxWidth="xl">
        <Stack spacing={3}>
          <Stack direction="row" alignItems="center" justifyContent="space-between" spacing={4}>
            <Typography variant="h4">Mes Absences</Typography>
            <Button
              variant="contained"
              startIcon={<Iconify icon="eva:plus-fill" />}
              onClick={() => setOpenForm(true)}
            >
              Soumettre un justificatif
            </Button>
          </Stack>

          {/* Stats Cards */}
          <Grid container spacing={3}>
            <Grid item xs={12} sm={6} md={3}>
              <Card sx={{ p: 3, textAlign: 'center' }}>
                <Typography variant="h4" color="primary">
                  {stats.totalAbsences}
                </Typography>
                <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                  Total des absences
                </Typography>
              </Card>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <Card sx={{ p: 3, textAlign: 'center' }}>
                <Typography variant="h4" color="warning.main">
                  {stats.pendingAbsences}
                </Typography>
                <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                  En attente
                </Typography>
              </Card>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <Card sx={{ p: 3, textAlign: 'center' }}>
                <Typography variant="h4" color="success.main">
                  {stats.approvedAbsences}
                </Typography>
                <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                  Approuvées
                </Typography>
              </Card>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <Card sx={{ p: 3, textAlign: 'center' }}>
                <Typography variant="h4" color="error.main">
                  {stats.rejectedAbsences}
                </Typography>
                <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                  Rejetées
                </Typography>
              </Card>
            </Grid>
          </Grid>

          {/* Absences Table */}
          <Card>
            <TableContainer>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Cours</TableCell>
                    <TableCell>Matière</TableCell>
                    <TableCell>Professeur</TableCell>
                    <TableCell>Date du cours</TableCell>
                    <TableCell>Justification</TableCell>
                    <TableCell>Statut</TableCell>
                    <TableCell>Date de soumission</TableCell>
                    <TableCell align="right">Actions</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {loading ? (
                    <TableRow>
                      <TableCell colSpan={8} align="center">
                        <CircularProgress />
                      </TableCell>
                    </TableRow>
                  ) : absences.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={8} align="center">
                        <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                          Aucune absence trouvée
                        </Typography>
                      </TableCell>
                    </TableRow>
                  ) : (
                    absences.map((absence) => (
                      <TableRow key={absence.id}>
                        <TableCell>{absence.course.name}</TableCell>
                        <TableCell>{absence.course.subject.name}</TableCell>
                        <TableCell>
                          {absence.course.teacher.first_name} {absence.course.teacher.last_name}
                        </TableCell>
                        <TableCell>
                          {new Date(absence.course.start_time).toLocaleDateString('fr-FR')}
                        </TableCell>
                        <TableCell>
                          <Typography variant="body2" noWrap sx={{ maxWidth: 200 }}>
                            {absence.justification}
                          </Typography>
                        </TableCell>
                        <TableCell>
                          <Chip
                            label={getStatusLabel(absence.status)}
                            color={getStatusColor(absence.status)}
                            size="small"
                          />
                        </TableCell>
                        <TableCell>
                          {new Date(absence.created_at).toLocaleDateString('fr-FR')}
                        </TableCell>
                        <TableCell align="right">
                          <Tooltip title="Voir les détails">
                            <IconButton
                              onClick={() => {
                                setSelectedAbsence(absence);
                                setOpenDetails(true);
                              }}
                            >
                              <Iconify icon="eva:eye-fill" />
                            </IconButton>
                          </Tooltip>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </TableContainer>
            <TablePagination
              rowsPerPageOptions={[5, 10, 25]}
              component="div"
              count={total}
              rowsPerPage={rowsPerPage}
              page={page}
              onPageChange={handleChangePage}
              onRowsPerPageChange={handleChangeRowsPerPage}
              labelRowsPerPage="Lignes par page"
            />
          </Card>
        </Stack>

        {/* Create Absence Dialog */}
        <Dialog open={openForm} onClose={() => setOpenForm(false)} maxWidth="sm" fullWidth>
          <DialogTitle>Soumettre un justificatif d'absence</DialogTitle>
          <DialogContent>
            <Stack spacing={3} sx={{ mt: 2 }}>
              <FormControl fullWidth>
                <InputLabel>Cours concerné *</InputLabel>
                <Select
                  value={formData.courseId}
                  onChange={(e) => setFormData({ ...formData, courseId: e.target.value })}
                  label="Cours concerné *"
                >
                  {courses.map((course) => (
                    <MenuItem key={course.id} value={course.id}>
                      {course.name} - {course.subject.name} ({new Date(course.start_time).toLocaleDateString('fr-FR')})
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>

              <TextField
                fullWidth
                multiline
                rows={4}
                label="Justification *"
                value={formData.justification}
                onChange={(e) => setFormData({ ...formData, justification: e.target.value })}
                placeholder="Expliquez la raison de votre absence..."
              />

              <TextField
                fullWidth
                label="Chemin du document (optionnel)"
                value={formData.documentPath}
                onChange={(e) => setFormData({ ...formData, documentPath: e.target.value })}
                placeholder="/uploads/justificatifs/document.pdf"
                helperText="Chemin vers le fichier justificatif (PDF, JPG, PNG)"
              />
            </Stack>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setOpenForm(false)}>Annuler</Button>
            <Button onClick={handleSubmit} variant="contained">
              Soumettre
            </Button>
          </DialogActions>
        </Dialog>

        {/* Absence Details Dialog */}
        <Dialog open={openDetails} onClose={() => setOpenDetails(false)} maxWidth="md" fullWidth>
          <DialogTitle>Détails de l'absence</DialogTitle>
          <DialogContent>
            {selectedAbsence && (
              <Stack spacing={3} sx={{ mt: 2 }}>
                <Grid container spacing={2}>
                  <Grid item xs={12} md={6}>
                    <Typography variant="subtitle2" color="text.secondary">
                      Cours
                    </Typography>
                    <Typography variant="body1">{selectedAbsence.course.name}</Typography>
                  </Grid>
                  <Grid item xs={12} md={6}>
                    <Typography variant="subtitle2" color="text.secondary">
                      Matière
                    </Typography>
                    <Typography variant="body1">{selectedAbsence.course.subject.name}</Typography>
                  </Grid>
                  <Grid item xs={12} md={6}>
                    <Typography variant="subtitle2" color="text.secondary">
                      Professeur
                    </Typography>
                    <Typography variant="body1">
                      {selectedAbsence.course.teacher.first_name} {selectedAbsence.course.teacher.last_name}
                    </Typography>
                  </Grid>
                  <Grid item xs={12} md={6}>
                    <Typography variant="subtitle2" color="text.secondary">
                      Date du cours
                    </Typography>
                    <Typography variant="body1">
                      {new Date(selectedAbsence.course.start_time).toLocaleDateString('fr-FR')}
                    </Typography>
                  </Grid>
                  <Grid item xs={12}>
                    <Typography variant="subtitle2" color="text.secondary">
                      Justification
                    </Typography>
                    <Typography variant="body1">{selectedAbsence.justification}</Typography>
                  </Grid>
                  {selectedAbsence.document_path && (
                    <Grid item xs={12}>
                      <Typography variant="subtitle2" color="text.secondary">
                        Document justificatif
                      </Typography>
                      <Typography variant="body1">{selectedAbsence.document_path}</Typography>
                    </Grid>
                  )}
                  <Grid item xs={12}>
                    <Typography variant="subtitle2" color="text.secondary">
                      Statut
                    </Typography>
                    <Chip
                      label={getStatusLabel(selectedAbsence.status)}
                      color={getStatusColor(selectedAbsence.status)}
                    />
                  </Grid>
                  {selectedAbsence.reviewer && (
                    <>
                      <Grid item xs={12} md={6}>
                        <Typography variant="subtitle2" color="text.secondary">
                          Validé par
                        </Typography>
                        <Typography variant="body1">
                          {selectedAbsence.reviewer.first_name} {selectedAbsence.reviewer.last_name}
                        </Typography>
                      </Grid>
                      <Grid item xs={12} md={6}>
                        <Typography variant="subtitle2" color="text.secondary">
                          Date de validation
                        </Typography>
                        <Typography variant="body1">
                          {selectedAbsence.reviewed_at
                            ? new Date(selectedAbsence.reviewed_at).toLocaleDateString('fr-FR')
                            : '-'}
                        </Typography>
                      </Grid>
                      <Grid item xs={12}>
                        <Typography variant="subtitle2" color="text.secondary">
                          Commentaire
                        </Typography>
                        <Typography variant="body1">{selectedAbsence.review_comment}</Typography>
                      </Grid>
                    </>
                  )}
                </Grid>
              </Stack>
            )}
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setOpenDetails(false)}>Fermer</Button>
          </DialogActions>
        </Dialog>
      </Container>
    </>
  );
} 