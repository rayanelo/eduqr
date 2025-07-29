import { useMemo } from 'react';

import { useLocales } from '../../../locales';

import { usePermissions } from '../../../hooks/usePermissions';

// components
import Iconify from '../../../components/iconify';

// ----------------------------------------------------------------------

export function useNavData() {
  const { translate } = useLocales();
  const { canManageUsers, canManageRooms, canManageSubjects, canAccessAuditLogs, canManageAbsences, canViewAbsences, canSubmitAbsences, canScanQRCode, canViewQRByRoom } = usePermissions();

  const data = useMemo(
    () => [
      // GENERAL
      // ----------------------------------------------------------------------
      {
        subheader: 'general',
        items: [
          {
            title: translate('profile'),
            path: '/dashboard/profile',
            icon: <Iconify icon="mdi:account" />,
          },
          {
            title: translate('calendar'),
            path: '/dashboard/calendar',
            icon: <Iconify icon="mdi:calendar" />,
          },
        ],
      },

      // MANAGEMENT
      // ----------------------------------------------------------------------
      {
        subheader: 'gestions',
        items: [
          ...(canManageUsers
            ? [
                {
                  title: translate('user'),
                  path: '/dashboard/user-management',
                  icon: <Iconify icon="mdi:account-group" />,
                },
              ]
            : []),
          ...(canManageRooms
            ? [
                {
                  title: 'Salles',
                  path: '/dashboard/room-management',
                  icon: <Iconify icon="mdi:office-building" />,
                },
              ]
            : []),
          ...(canManageSubjects
            ? [
                {
                  title: 'Matières',
                  path: '/dashboard/subject-management',
                  icon: <Iconify icon="mdi:book-open-variant" />,
                },
                {
                  title: 'Cours',
                  path: '/dashboard/course-management',
                  icon: <Iconify icon="mdi:calendar" />,
                },
              ]
            : []),
          ...(canAccessAuditLogs
            ? [
                {
                  title: 'Journal d\'Activité',
                  path: '/dashboard/audit-logs',
                  icon: <Iconify icon="mdi:file-document" />,
                },
              ]
            : []),
          ...(canSubmitAbsences
            ? [
                {
                  title: 'Mes Absences',
                  path: '/dashboard/my-absences',
                  icon: <Iconify icon="mdi:account-clock" />,
                },
              ]
            : []),
          ...(canViewAbsences && !canManageAbsences
            ? [
                {
                  title: 'Absences à Traiter',
                  path: '/dashboard/teacher-absences',
                  icon: <Iconify icon="mdi:check-circle" />,
                },
              ]
            : []),
          ...(canManageAbsences
            ? [
                {
                  title: 'Gestion des Absences',
                  path: '/dashboard/admin-absences',
                  icon: <Iconify icon="mdi:calendar-remove" />,
                },
              ]
            : []),
          ...(canViewQRByRoom
            ? [
                {
                  title: 'QR Codes par Salle',
                  path: '/dashboard/qr-by-room',
                  icon: <Iconify icon="mdi:qrcode" />,
                },
              ]
            : []),
          ...(canScanQRCode
            ? [
                {
                  title: 'Scanner QR Code',
                  path: '/dashboard/qr-scanner',
                  icon: <Iconify icon="mdi:qrcode-scan" />,
                },
              ]
            : []),
        ],
      },
    ],
    [translate, canManageUsers, canManageRooms, canManageSubjects, canAccessAuditLogs, canManageAbsences, canViewAbsences, canSubmitAbsences, canScanQRCode, canViewQRByRoom]
  );

  return data;
}

// Configuration par défaut pour la barre de recherche (utilise useNavData)
const navConfig = useNavData;

export default navConfig;
