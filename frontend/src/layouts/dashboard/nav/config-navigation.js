import { useMemo } from 'react';

import { useLocales } from '../../../locales';

import { usePermissions } from '../../../hooks/usePermissions';



// components
import SvgColor from '../../../components/svg-color';

// ----------------------------------------------------------------------

export function useNavData() {
  const { translate } = useLocales();
  const { canManageUsers, canManageRooms, canManageSubjects } = usePermissions();

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
        ],
      },
    ],
    [translate, canManageUsers, canManageRooms, canManageSubjects]
  );

  return data;
}

// Export par défaut pour compatibilité
const navConfig = [
  // GENERAL
  // ----------------------------------------------------------------------
  {
    subheader: 'general v4.2.0',
    items: [
      { title: 'One', path: '/dashboard/one', icon: <SvgColor src="/assets/icons/navbar/ic_dashboard.svg" /> },
      { title: 'Two', path: '/dashboard/two', icon: <SvgColor src="/assets/icons/navbar/ic_ecommerce.svg" /> },
      { title: 'Three', path: '/dashboard/three', icon: <SvgColor src="/assets/icons/navbar/ic_analytics.svg" /> },
      { title: 'Calendar', path: '/dashboard/calendar', icon: <SvgColor src="/assets/icons/navbar/ic_calendar.svg" /> },
      { title: 'Mon profil', path: '/dashboard/profile', icon: <SvgColor src="/assets/icons/navbar/ic_user.svg" /> },
    ],
  },

  // MANAGEMENT
  // ----------------------------------------------------------------------
  {
    subheader: 'management',
    items: [
      {
        title: 'user',
        path: '/dashboard/user',
        icon: <SvgColor src="/assets/icons/navbar/ic_user.svg" />,
        children: [
          { title: 'Four', path: '/dashboard/user/four' },
          { title: 'Five', path: '/dashboard/user/five' },
          { title: 'Six', path: '/dashboard/user/six' },
          { title: 'Gestion des utilisateurs', path: '/dashboard/user-management' },
        ],
      },
    ],
  },
];

export default navConfig;
