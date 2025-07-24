import { Navigate, useRoutes } from 'react-router-dom';

// layouts
import DashboardLayout from '../layouts/dashboard';
import CompactLayout from '../layouts/compact';

// guards
import GuestGuard from '../auth/GuestGuard';
import AuthGuard from '../auth/AuthGuard';

// import pages
import LoginPage from '../pages/LoginPage';
import Page404 from '../pages/Page404';
import PageFive from '../pages/PageFive';
import PageSix from '../pages/PageSix';
import UserManagementPage from '../pages/dashboard/UserManagementPage';
import RoomManagementPage from '../pages/dashboard/RoomManagementPage';
import SubjectManagementPage from '../pages/dashboard/SubjectManagementPage';
import ProfilePage from '../pages/dashboard/ProfilePage';
import CalendarPage from '../pages/dashboard/CalendarPage';

// ----------------------------------------------------------------------

export default function Router() {
  return useRoutes([
    {
      path: 'dashboard',
      element: (
        <AuthGuard>
          <DashboardLayout />
        </AuthGuard>
      ),
      children: [
        { element: <Navigate to="/dashboard/user-management" replace />, index: true },
        { path: 'user-management', element: <UserManagementPage /> },
        { path: 'room-management', element: <RoomManagementPage /> },
        { path: 'subject-management', element: <SubjectManagementPage /> },
        { path: 'profile', element: <ProfilePage /> },
        { path: 'calendar', element: <CalendarPage /> },
        { path: 'five', element: <PageFive /> },
        { path: 'six', element: <PageSix /> },
      ],
    },
    {
      path: '/',
      element: <Navigate to="/login" replace />,
    },
    {
      path: 'login',
      element: (
        <GuestGuard>
          <LoginPage />
        </GuestGuard>
      ),
    },
    {
      element: <CompactLayout />,
      children: [{ path: '404', element: <Page404 /> }],
    },
    { path: '*', element: <Navigate to="/404" replace /> },
  ]);
}
