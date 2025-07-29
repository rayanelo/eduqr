import { useMemo } from 'react';

import { useLocales } from '../../../locales';

import { usePermissions } from '../../../hooks/usePermissions';

// components
import SvgColor from '../../../components/svg-color';

// ----------------------------------------------------------------------

export function useNavData() {
  const { translate } = useLocales();
  const { canManageUsers, canManageRooms, canManageSubjects, canAccessAuditLogs } = usePermissions();

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
            icon: <SvgColor src="/assets/icons/navbar/ic_user.svg" />,
          },
          {
            title: translate('calendar'),
            path: '/dashboard/calendar',
            icon: <SvgColor src="/assets/icons/navbar/ic_calendar.svg" />,
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
                  icon: <SvgColor src="/assets/icons/navbar/ic_user.svg" />,
                },
              ]
            : []),
          ...(canManageRooms
            ? [
                {
                  title: 'Salles',
                  path: '/dashboard/room-management',
                  icon: <SvgColor src="/assets/icons/navbar/ic_booking.svg" />,
                },
              ]
            : []),
          ...(canManageSubjects
            ? [
                {
                  title: 'Matières',
                  path: '/dashboard/subject-management',
                  icon: <SvgColor src="/assets/icons/navbar/ic_file.svg" />,
                },
                {
                  title: 'Cours',
                  path: '/dashboard/course-management',
                  icon: <SvgColor src="/assets/icons/navbar/ic_kanban.svg" />,
                },
              ]
            : []),
          ...(canAccessAuditLogs
            ? [
                {
                  title: 'Journal d\'Activité',
                  path: '/dashboard/audit-logs',
                  icon: <SvgColor src="/assets/icons/navbar/ic_analytics.svg" />,
                },
              ]
            : []),
        ],
      },
    ],
    [translate, canManageUsers, canManageRooms, canManageSubjects, canAccessAuditLogs]
  );

  return data;
}

// Configuration par défaut pour la barre de recherche (utilise useNavData)
const navConfig = useNavData;

export default navConfig;
