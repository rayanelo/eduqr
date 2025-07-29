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
  Box,
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
  Divider,
} from '@mui/material';
// hooks
import { useSnackbar } from 'notistack';
// components
import Iconify from '../../components/iconify';
import apiClient from '../../utils/api';

// ----------------------------------------------------------------------

export default function AdminAbsencePage() {
  const { enqueueSnackbar } = useSnackbar();

  const [absences, setAbsences] = useState([]);
  const [loading, setLoading] = useState(false);
  const [stats, setStats] = useState({
    total_absences: 0,
    pending_absences: 0,
    approved_absences: 0,
    rejected_absences: 0,
  });

  // Pagination
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [total, setTotal] = useState(0);

  // Dialog states
  const [openReview, setOpenReview] = useState(false);
  const [openDetails, setOpenDetails] = useState(false);
  const [openFilters, setOpenFilters] = useState(false);
  const [selectedAbsence, setSelectedAbsence] = useState(null);

  // Review form states
  const [reviewData, setReviewData] = useState({
    status: '',
    reviewComment: '',
  });

  // Filter states
  const [filters, setFilters] = useState({
    studentId: '',
    courseId: '',
    status: '',
    startDate: '',
    endDate: '',
    page: 1,
    limit: 10,
  });

  // Load absences
  const loadAbsences = useCallback(async () => {
    setLoading(true);
    try {
              const response = await apiClient.get(`/api/v1/admin/absences?page=${page + 1}&limit=${rowsPerPage}`);
      setAbsences(response.data.data || []);
      setTotal(response.data.total || 0);
    } catch (error) {
      console.error('Error loading absences:', error);
      enqueueSnackbar('Erreur lors du chargement des absences', { variant: 'error' });
    } finally {
      setLoading(false);
    }
  }, [page, rowsPerPage, enqueueSnackbar]);

  // Load absences with filters
  const loadAbsencesWithFilters = useCallback(async () => {
    setLoading(true);
    try {
      const response = await apiClient.post('/absences/filter', {
        ...filters,
        page: page + 1,
        limit: rowsPerPage,
      });
      setAbsences(response.data.data || []);
      setTotal(response.data.total || 0);
    } catch (error) {
      console.error('Error loading absences with filters:', error);
      enqueueSnackbar('Erreur lors du chargement des absences', { variant: 'error' });
    } finally {
      setLoading(false);
    }
  }, [filters, page, rowsPerPage, enqueueSnackbar]);

  // Load stats
  const loadStats = useCallback(async () => {
    try {
      const response = await apiClient.get('/api/v1/absences/stats');
      setStats(response.data);
    } catch (error) {
      console.error('Error loading stats:', error);
    }
  }, []);

  // Load data on component mount
  useEffect(() => {
    loadAbsences();
    loadStats();
  }, [loadAbsences, loadStats]);

  // Handle review submission
  const handleReviewSubmit = async () => {
    if (!reviewData.status || !reviewData.reviewComment) {
      enqueueSnackbar('Veuillez remplir tous les champs obligatoires', { variant: 'warning' });
      return;
    }

    try {
      await apiClient.post(`/api/v1/absences/${selectedAbsence.id}/review`, reviewData);
      enqueueSnackbar('Absence traitée avec succès', { variant: 'success' });
      setOpenReview(false);
      setReviewData({ status: '', reviewComment: '' });
      setSelectedAbsence(null);
      loadAbsences();
      loadStats();
    } catch (error) {
      console.error('Error reviewing absence:', error);
      enqueueSnackbar(error.response?.data?.error || 'Erreur lors du traitement de l\'absence', { variant: 'error' });
    }
  };

  // Handle filter submission
  const handleFilterSubmit = () => {
    setPage(0);
    loadAbsencesWithFilters();
    setOpenFilters(false);
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

  // Open review dialog
  const handleOpenReview = (absence) => {
    setSelectedAbsence(absence);
    setOpenReview(true);
  };

  return (
    <>
      <Helmet>
        <title>Gestion Globale des Absences | EduQR</title>
      </Helmet>

      <Container maxWidth="xl">
        <Stack spacing={3}>
          <Stack direction="row" alignItems="center" justifyContent="space-between" spacing={4}>
            <Typography variant="h4">Gestion Globale des Absences</Typography>
            <Stack direction="row" spacing={2}>
              <Button
                variant="outlined"
                startIcon={<Iconify icon="eva:filter-fill" />}
                onClick={() => setOpenFilters(true)}
              >
                Filtres
              </Button>
            </Stack>
          </Stack>

          {/* Stats Cards */}
          <Grid container spacing={3}>
            <Grid item xs={12} sm={6} md={3}>
              <Card sx={{ p: 3, textAlign: 'center' }}>
                <Typography variant="h4" color="primary">
                  {stats.total_absences}
                </Typography>
                <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                  Total des absences
                </Typography>
              </Card>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <Card sx={{ p: 3, textAlign: 'center' }}>
                <Typography variant="h4" color="warning.main">
                  {stats.pending_absences}
                </Typography>
                <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                  En attente
                </Typography>
              </Card>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <Card sx={{ p: 3, textAlign: 'center' }}>
                <Typography variant="h4" color="success.main">
                  {stats.approved_absences}
                </Typography>
                <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                  Approuvées
                </Typography>
              </Card>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <Card sx={{ p: 3, textAlign: 'center' }}>
                <Typography variant="h4" color="error.main">
                  {stats.rejected_absences}
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
                    <TableCell>Étudiant</TableCell>
                    <TableCell>Cours</TableCell>
                    <TableCell>Date et heure du cours</TableCell>
                    <TableCell>Justification</TableCell>
                    <TableCell>Statut</TableCell>
                    <TableCell>Validé par</TableCell>
                    <TableCell align="right">Actions</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {loading ? (
                    <TableRow>
                      <TableCell colSpan={7} align="center">
                        <CircularProgress />
                      </TableCell>
                    </TableRow>
                  ) : absences.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={7} align="center">
                        <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                          Aucune absence trouvée
                        </Typography>
                      </TableCell>
                    </TableRow>
                  ) : (
                    absences.map((absence) => (
                      <TableRow key={absence.id}>
                        <TableCell>
                          {absence.student.first_name} {absence.student.last_name}
                        </TableCell>
                        <TableCell>
                          <Tooltip title={absence.course.name} arrow>
                            <Typography variant="body2" noWrap>
                              {absence.course.name.length > 15 
                                ? `${absence.course.name.substring(0, 15)}...` 
                                : absence.course.name}
                            </Typography>
                          </Tooltip>
                        </TableCell>
                        <TableCell>
                          {new Date(absence.course.start_time).toLocaleDateString('fr-FR')} à {new Date(absence.course.start_time).toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' })}
                        </TableCell>
                        <TableCell>
                          <Tooltip title={absence.justification} arrow>
                            <Typography variant="body2" noWrap sx={{ maxWidth: 200 }}>
                              {absence.justification.length > 15 
                                ? `${absence.justification.substring(0, 15)}...` 
                                : absence.justification}
                            </Typography>
                          </Tooltip>
                        </TableCell>
                        <TableCell>
                          <Chip
                            label={getStatusLabel(absence.status)}
                            color={getStatusColor(absence.status)}
                            size="small"
                          />
                        </TableCell>
                        <TableCell>
                          {absence.reviewer ? (
                            `${absence.reviewer.first_name} ${absence.reviewer.last_name}`
                          ) : (
                            <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                              -
                            </Typography>
                          )}
                        </TableCell>
                        <TableCell align="right">
                          <Stack direction="row" spacing={1} justifyContent="flex-end">
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
                            {absence.status === 'pending' && (
                              <Tooltip title="Traiter l'absence">
                                <IconButton
                                  color="primary"
                                  onClick={() => handleOpenReview(absence)}
                                >
                                  <Iconify icon="eva:edit-fill" />
                                </IconButton>
                              </Tooltip>
                            )}
                          </Stack>
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

        {/* Filters Dialog */}
        <Dialog open={openFilters} onClose={() => setOpenFilters(false)} maxWidth="sm" fullWidth>
          <DialogTitle>Filtres avancés</DialogTitle>
          <DialogContent>
            <Stack spacing={3} sx={{ mt: 2 }}>
              <TextField
                fullWidth
                label="ID Étudiant"
                value={filters.studentId}
                onChange={(e) => setFilters({ ...filters, studentId: e.target.value })}
                placeholder="Ex: 1"
              />

              <TextField
                fullWidth
                label="ID Cours"
                value={filters.courseId}
                onChange={(e) => setFilters({ ...filters, courseId: e.target.value })}
                placeholder="Ex: 1"
              />

              <FormControl fullWidth>
                <InputLabel>Statut</InputLabel>
                <Select
                  value={filters.status}
                  onChange={(e) => setFilters({ ...filters, status: e.target.value })}
                  label="Statut"
                >
                  <MenuItem value="">Tous</MenuItem>
                  <MenuItem value="pending">En attente</MenuItem>
                  <MenuItem value="approved">Approuvée</MenuItem>
                  <MenuItem value="rejected">Rejetée</MenuItem>
                </Select>
              </FormControl>

              <TextField
                fullWidth
                label="Date de début"
                type="date"
                value={filters.startDate}
                onChange={(e) => setFilters({ ...filters, startDate: e.target.value })}
                InputLabelProps={{ shrink: true }}
              />

              <TextField
                fullWidth
                label="Date de fin"
                type="date"
                value={filters.endDate}
                onChange={(e) => setFilters({ ...filters, endDate: e.target.value })}
                InputLabelProps={{ shrink: true }}
              />
            </Stack>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setOpenFilters(false)}>Annuler</Button>
            <Button onClick={handleFilterSubmit} variant="contained">
              Appliquer
            </Button>
          </DialogActions>
        </Dialog>

        {/* Review Absence Dialog */}
        <Dialog open={openReview} onClose={() => setOpenReview(false)} maxWidth="sm" fullWidth>
          <DialogTitle>Traiter l'absence</DialogTitle>
          <DialogContent>
            {selectedAbsence && (
              <Stack spacing={3} sx={{ mt: 2 }}>
                <Box>
                  <Typography variant="subtitle2" color="text.secondary">
                    Étudiant
                  </Typography>
                  <Typography variant="body1">
                    {selectedAbsence.student.first_name} {selectedAbsence.student.last_name}
                  </Typography>
                </Box>
                <Box>
                  <Typography variant="subtitle2" color="text.secondary">
                    Cours
                  </Typography>
                  <Typography variant="body1">
                    {selectedAbsence.course.name} - {selectedAbsence.course.subject.name}
                  </Typography>
                </Box>
                <Box>
                  <Typography variant="subtitle2" color="text.secondary">
                    Professeur
                  </Typography>
                  <Typography variant="body1">
                    {selectedAbsence.course.teacher.first_name} {selectedAbsence.course.teacher.last_name}
                  </Typography>
                </Box>
                <Box>
                  <Typography variant="subtitle2" color="text.secondary">
                    Justification
                  </Typography>
                  <Typography variant="body1">{selectedAbsence.justification}</Typography>
                </Box>
                {selectedAbsence.document_path && (
                  <Box>
                    <Typography variant="subtitle2" color="text.secondary">
                      Document justificatif
                    </Typography>
                    <Typography variant="body1">{selectedAbsence.document_path}</Typography>
                  </Box>
                )}
                <Divider />
                <FormControl fullWidth>
                  <InputLabel>Décision *</InputLabel>
                  <Select
                    value={reviewData.status}
                    onChange={(e) => setReviewData({ ...reviewData, status: e.target.value })}
                    label="Décision *"
                  >
                    <MenuItem value="approved">Approuver</MenuItem>
                    <MenuItem value="rejected">Rejeter</MenuItem>
                  </Select>
                </FormControl>

                <TextField
                  fullWidth
                  multiline
                  rows={4}
                  label="Commentaire (optionnel)"
                  value={reviewData.reviewComment}
                  onChange={(e) => setReviewData({ ...reviewData, reviewComment: e.target.value })}
                  placeholder="Expliquez votre décision (optionnel)..."
                />
              </Stack>
            )}
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setOpenReview(false)}>Annuler</Button>
            <Button onClick={handleReviewSubmit} variant="contained">
              Valider
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
                      Étudiant
                    </Typography>
                    <Typography variant="body1">
                      {selectedAbsence.student.first_name} {selectedAbsence.student.last_name}
                    </Typography>
                  </Grid>
                  <Grid item xs={12} md={6}>
                    <Typography variant="subtitle2" color="text.secondary">
                      Email
                    </Typography>
                    <Typography variant="body1">{selectedAbsence.student.email}</Typography>
                  </Grid>
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
                          Traité par
                        </Typography>
                        <Typography variant="body1">
                          {selectedAbsence.reviewer.first_name} {selectedAbsence.reviewer.last_name}
                        </Typography>
                      </Grid>
                      <Grid item xs={12} md={6}>
                        <Typography variant="subtitle2" color="text.secondary">
                          Date de traitement
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
            {selectedAbsence && selectedAbsence.status === 'pending' && (
              <Button
                onClick={() => {
                  setOpenDetails(false);
                  handleOpenReview(selectedAbsence);
                }}
                variant="contained"
              >
                Traiter
              </Button>
            )}
          </DialogActions>
        </Dialog>
      </Container>
    </>
  );
} 