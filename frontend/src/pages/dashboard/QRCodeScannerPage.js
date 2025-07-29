import { useState, useCallback } from 'react';
import { Helmet } from 'react-helmet-async';
// @mui
import {
  Card,
  Stack,
  Button,
  Container,
  Typography,
  Box,
  Alert,
  CircularProgress,
  Chip,
  Grid,
  Paper,
} from '@mui/material';
// hooks
import { useSnackbar } from 'notistack';
// components
import Iconify from '../../components/iconify';
import apiClient from '../../utils/api';

// ----------------------------------------------------------------------

export default function QRCodeScannerPage() {
  const { enqueueSnackbar } = useSnackbar();

  const [scanning, setScanning] = useState(false);
  const [scanResult, setScanResult] = useState(null);
  const [loading, setLoading] = useState(false);

  // Handle QR code scan
  const handleScan = useCallback(async (qrCodeData) => {
    if (!qrCodeData) return;
    
    setLoading(true);
    try {
      const response = await apiClient.post('/presences/scan', {
        qr_code_data: qrCodeData,
      });
      
      setScanResult(response.data);
      enqueueSnackbar('Présence enregistrée avec succès !', { variant: 'success' });
    } catch (error) {
      console.error('Error scanning QR code:', error);
      enqueueSnackbar(error.response?.data?.error || 'Erreur lors du scan du QR code', { variant: 'error' });
    } finally {
      setLoading(false);
    }
  }, [enqueueSnackbar]);

  // Simulate QR code scanning (in a real app, this would use a camera)
  const simulateScan = () => {
    setScanning(true);
    // Simulate scanning process
    setTimeout(() => {
      const mockQRData = 'eyJjb3Vyc2VfaWQiOjEsInRva2VuIjoiZXhhbXBsZS10b2tlbiIsInRpbWVzdGFtcCI6MTYzNTQzMjEwMH0='; // Mock QR data
      handleScan(mockQRData);
      setScanning(false);
    }, 2000);
  };

  // Get status color
  const getStatusColor = (status) => {
    switch (status) {
      case 'present':
        return 'success';
      case 'late':
        return 'warning';
      case 'absent':
        return 'error';
      default:
        return 'default';
    }
  };

  // Get status label
  const getStatusLabel = (status) => {
    switch (status) {
      case 'present':
        return 'Présent';
      case 'late':
        return 'En retard';
      case 'absent':
        return 'Absent';
      default:
        return status;
    }
  };

  // Format time
  const formatTime = (dateString) => {
    if (!dateString) return '-';
    return new Date(dateString).toLocaleString('fr-FR');
  };

  return (
    <>
      <Helmet>
        <title>Scanner QR Code | EduQR</title>
      </Helmet>

      <Container maxWidth="md">
        <Stack spacing={3}>
          <Typography variant="h4" textAlign="center">
            Scanner QR Code de Présence
          </Typography>

          {/* Scanner Interface */}
          <Card sx={{ p: 4, textAlign: 'center' }}>
            <Typography variant="h6" gutterBottom>
              Scanner le QR Code
            </Typography>
            
            <Box sx={{ my: 4 }}>
              {scanning ? (
                <Box sx={{ position: 'relative' }}>
                  <CircularProgress size={120} />
                  <Typography variant="body1" sx={{ mt: 2 }}>
                    Scan en cours...
                  </Typography>
                </Box>
              ) : (
                <Box
                  sx={{
                    width: 300,
                    height: 300,
                    border: '2px dashed',
                    borderColor: 'grey.300',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    mx: 'auto',
                    mb: 3,
                  }}
                >
                  <Typography variant="body2" color="text.secondary">
                    Zone de scan
                  </Typography>
                </Box>
              )}
            </Box>

            <Stack direction="row" spacing={2} justifyContent="center">
              <Button
                variant="contained"
                size="large"
                startIcon={<Iconify icon="eva:camera-fill" />}
                onClick={simulateScan}
                disabled={scanning || loading}
              >
                {scanning ? 'Scan en cours...' : 'Scanner QR Code'}
              </Button>
            </Stack>

            <Typography variant="body2" color="text.secondary" sx={{ mt: 3 }}>
              Pointez votre caméra vers le QR code affiché dans la salle
            </Typography>
          </Card>

          {/* Scan Result */}
          {scanResult && (
            <Card sx={{ p: 4 }}>
              <Typography variant="h6" gutterBottom>
                Résultat du Scan
              </Typography>
              
              <Grid container spacing={3}>
                <Grid item xs={12} md={6}>
                  <Box>
                    <Typography variant="subtitle2" color="text.secondary">
                      Cours
                    </Typography>
                    <Typography variant="body1">
                      {scanResult.presence.course.name}
                    </Typography>
                  </Box>
                </Grid>
                <Grid item xs={12} md={6}>
                  <Box>
                    <Typography variant="subtitle2" color="text.secondary">
                      Matière
                    </Typography>
                    <Typography variant="body1">
                      {scanResult.presence.course.subject.name}
                    </Typography>
                  </Box>
                </Grid>
                <Grid item xs={12} md={6}>
                  <Box>
                    <Typography variant="subtitle2" color="text.secondary">
                      Professeur
                    </Typography>
                    <Typography variant="body1">
                      {scanResult.presence.course.teacher.first_name} {scanResult.presence.course.teacher.last_name}
                    </Typography>
                  </Box>
                </Grid>
                <Grid item xs={12} md={6}>
                  <Box>
                    <Typography variant="subtitle2" color="text.secondary">
                      Salle
                    </Typography>
                    <Typography variant="body1">
                      {scanResult.presence.course.room.name}
                    </Typography>
                  </Box>
                </Grid>
                <Grid item xs={12} md={6}>
                  <Box>
                    <Typography variant="subtitle2" color="text.secondary">
                      Statut
                    </Typography>
                    <Chip
                      label={getStatusLabel(scanResult.presence.status)}
                      color={getStatusColor(scanResult.presence.status)}
                      sx={{ mt: 1 }}
                    />
                  </Box>
                </Grid>
                <Grid item xs={12} md={6}>
                  <Box>
                    <Typography variant="subtitle2" color="text.secondary">
                      Heure de scan
                    </Typography>
                    <Typography variant="body1">
                      {formatTime(scanResult.presence.scanned_at)}
                    </Typography>
                  </Box>
                </Grid>
              </Grid>

              <Alert severity="success" sx={{ mt: 3 }}>
                Votre présence a été enregistrée avec succès !
              </Alert>
            </Card>
          )}

          {/* Instructions */}
          <Paper sx={{ p: 3, bgcolor: 'background.neutral' }}>
            <Typography variant="h6" gutterBottom>
              Instructions
            </Typography>
            
            <Stack spacing={2}>
              <Typography variant="body2">
                • Assurez-vous que le QR code est bien visible sur l'écran de la salle
              </Typography>
              <Typography variant="body2">
                • Pointez votre caméra vers le QR code
              </Typography>
              <Typography variant="body2">
                • Votre présence sera automatiquement enregistrée
              </Typography>
              <Typography variant="body2" color="warning.main">
                ⚠️ Vous ne pouvez scanner qu'une seule fois par cours
              </Typography>
              <Typography variant="body2" color="error.main">
                ⚠️ Le QR code n'est valide que pendant la durée du cours
              </Typography>
            </Stack>
          </Paper>

          {/* Troubleshooting */}
          <Paper sx={{ p: 3, bgcolor: 'background.neutral' }}>
            <Typography variant="h6" gutterBottom>
              Problèmes courants
            </Typography>
            
            <Stack spacing={2}>
              <Typography variant="body2">
                <strong>Le scan ne fonctionne pas :</strong>
              </Typography>
              <Typography variant="body2" sx={{ ml: 2 }}>
                • Vérifiez que le QR code est bien affiché dans la salle
              </Typography>
              <Typography variant="body2" sx={{ ml: 2 }}>
                • Assurez-vous que le cours est en cours
              </Typography>
              <Typography variant="body2" sx={{ ml: 2 }}>
                • Contactez votre professeur si le problème persiste
              </Typography>
            </Stack>
          </Paper>
        </Stack>
      </Container>
    </>
  );
} 