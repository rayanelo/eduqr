import { Suspense, lazy } from 'react';
// components
import LoadingScreen from '../components/loading-screen';

// ----------------------------------------------------------------------

const Loadable = (Component) => (props) =>
  (
    <Suspense fallback={<LoadingScreen />}>
      <Component {...props} />
    </Suspense>
  );

// ----------------------------------------------------------------------

export const LoginPage = Loadable(lazy(() => import('../pages/LoginPage')));
export const RegisterPage = Loadable(lazy(() => import('../pages/RegisterPage')));

export const PageOne = Loadable(lazy(() => import('../pages/PageOne')));
export const PageTwo = Loadable(lazy(() => import('../pages/PageTwo')));
export const PageThree = Loadable(lazy(() => import('../pages/PageThree')));
export const PageFour = Loadable(lazy(() => import('../pages/PageFour')));
export { default as PageFive } from '../pages/PageFive';
export { default as PageSix } from '../pages/PageSix';
export { default as UserManagementPage } from '../pages/dashboard/UserManagementPage';
export { default as RoomManagementPage } from '../pages/dashboard/RoomManagementPage';
export { default as SubjectManagementPage } from '../pages/dashboard/SubjectManagementPage';
export { default as CourseManagementPage } from '../pages/dashboard/CourseManagementPage';
export { default as ProfilePage } from '../pages/dashboard/ProfilePage';
export { default as CalendarPage } from '../pages/dashboard/CalendarPage';
export { default as AuditLogPage } from '../pages/dashboard/AuditLogPage';

export const Page404 = Loadable(lazy(() => import('../pages/Page404')));
