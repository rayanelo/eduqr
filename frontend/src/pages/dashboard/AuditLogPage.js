import { useState, useEffect, useCallback } from 'react';
import { Helmet } from 'react-helmet-async';
import {
  Box,
  Card,
  Container,
  Button,
  Tooltip,
  Alert,
} from '@mui/material';
import { useSnackbar } from 'notistack';

import { useAuditLogs } from '../../hooks/useAuditLogs';
import { usePermissions } from '../../hooks/usePermissions';
import { useContext } from 'react';
import { AuthContext } from '../../auth/JwtContext';

import Iconify from '../../components/iconify';
import { useSettingsContext } from '../../components/settings';
import CustomBreadcrumbs from '../../components/custom-breadcrumbs';
import AuditLogTable from '../../sections/audit/AuditLogTable';
import AuditLogFilters from '../../sections/audit/AuditLogFilters';
import AuditLogStats from '../../sections/audit/AuditLogStats';
import AuditLogDetailDialog from '../../sections/audit/AuditLogDetailDialog';

export default function AuditLogPage() {
  const { themeStretch } = useSettingsContext();
  const { enqueueSnackbar } = useSnackbar();
  const { canAccessAuditLogs } = usePermissions();
  const { user } = useContext(AuthContext);

  const {
    logs,
    loading,
    pagination,
    fetchAuditLogs,
    fetchAuditStats,
    cleanOldLogs,
    fetchAuditLogById,
  } = useAuditLogs();

  const [filters, setFilters] = useState({
    page: 1,
    limit: 20,
    action: '',
    resource_type: '',
    resource_id: '',
    user_id: '',
    start_date: '',
    end_date: '',
    search: '',
  });

  const [stats, setStats] = useState(null);
  const [selectedLog, setSelectedLog] = useState(null);
  const [detailDialogOpen, setDetailDialogOpen] = useState(false);

  const loadData = useCallback(async () => {
    try {
      await fetchAuditLogs(filters);
      
      // Charger les statistiques pour les 30 derniers jours
      const endDate = new Date().toISOString().split('T')[0];
      const startDate = new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];
      const statsData = await fetchAuditStats(startDate, endDate);
      setStats(statsData);
    } catch (error) {
      console.error('Erreur lors du chargement des données:', error);
    }
  }, [fetchAuditLogs, fetchAuditStats, filters]);

  // Charger les données initiales
  useEffect(() => {
    console.log('🔍 Debug AuditLogPage:', {
      canAccessAuditLogs,
      user: user,
      filters
    });
    
    if (canAccessAuditLogs) {
      loadData();
    } else {
      console.log('❌ Accès refusé aux logs d\'audit');
    }
  }, [canAccessAuditLogs, loadData]);

  const handleFilterChange = (newFilters) => {
    setFilters({ ...newFilters, page: 1 }); // Reset to first page when filters change
  };

  const handlePageChange = (newPage) => {
    setFilters({ ...filters, page: newPage });
  };

  const handleViewLogDetail = async (logId) => {
    try {
      const logDetail = await fetchAuditLogById(logId);
      setSelectedLog(logDetail);
      setDetailDialogOpen(true);
    } catch (error) {
      enqueueSnackbar('Erreur lors de la récupération des détails du log', { variant: 'error' });
    }
  };

  const handleCleanOldLogs = async () => {
    try {
      await cleanOldLogs(365); // Nettoyer les logs de plus d'un an
      enqueueSnackbar('Anciens logs nettoyés avec succès', { variant: 'success' });
      loadData(); // Recharger les données
    } catch (error) {
      enqueueSnackbar('Erreur lors du nettoyage des anciens logs', { variant: 'error' });
    }
  };

  const getActionColor = (action) => {
    switch (action) {
      case 'create':
        return 'success';
      case 'update':
        return 'info';
      case 'delete':
        return 'error';
      case 'login':
        return 'primary';
      case 'logout':
        return 'warning';
      default:
        return 'default';
    }
  };

  const getActionLabel = (action) => {
    switch (action) {
      case 'create':
        return 'Création';
      case 'update':
        return 'Modification';
      case 'delete':
        return 'Suppression';
      case 'login':
        return 'Connexion';
      case 'logout':
        return 'Déconnexion';
      default:
        return action;
    }
  };

  const getResourceTypeLabel = (resourceType) => {
    switch (resourceType) {
      case 'user':
        return 'Utilisateur';
      case 'room':
        return 'Salle';
      case 'subject':
        return 'Matière';
      case 'course':
        return 'Cours';
      case 'event':
        return 'Événement';
      default:
        return resourceType;
    }
  };

  // Vérifier les permissions
  if (!canAccessAuditLogs) {
    return (
      <Container maxWidth={themeStretch ? false : 'lg'}>
        <Alert severity="error" sx={{ mt: 3 }}>
          Vous n'avez pas les permissions nécessaires pour accéder aux logs d'audit.
        </Alert>
      </Container>
    );
  }

  return (
    <>
      <Helmet>
        <title> Journal d'Activité | EduQR</title>
      </Helmet>

      <Container maxWidth={themeStretch ? false : 'xl'}>
        <CustomBreadcrumbs
          heading="Journal d'Activité"
          links={[
            { name: 'Tableau de bord', href: '/dashboard' },
            { name: 'Journal d\'Activité' },
          ]}
          action={
            <Tooltip title="Nettoyer les anciens logs">
              <Button
                variant="outlined"
                color="warning"
                startIcon={<Iconify icon="eva:trash-2-fill" />}
                onClick={handleCleanOldLogs}
              >
                Nettoyer
              </Button>
            </Tooltip>
          }
          sx={{
            mb: { xs: 3, md: 5 },
          }}
        />

        {/* Statistiques */}
        {stats && (
          <Box sx={{ mb: 3 }}>
            <AuditLogStats stats={stats} />
          </Box>
        )}

        {/* Filtres */}
        <Card sx={{ mb: 3 }}>
          <AuditLogFilters
            filters={filters}
            onFilterChange={handleFilterChange}
            loading={loading}
          />
        </Card>

        {/* Tableau des logs */}
        <Card>
          <AuditLogTable
            logs={logs || []}
            loading={loading}
            pagination={pagination}
            onPageChange={handlePageChange}
            onViewDetail={handleViewLogDetail}
            getActionColor={getActionColor}
            getActionLabel={getActionLabel}
            getResourceTypeLabel={getResourceTypeLabel}
          />
        </Card>

        {/* Dialog de détails */}
        <AuditLogDetailDialog
          open={detailDialogOpen}
          log={selectedLog}
          onClose={() => {
            setDetailDialogOpen(false);
            setSelectedLog(null);
          }}
        />
      </Container>
    </>
  );
} 