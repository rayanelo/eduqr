// @mui
import { Alert, Tooltip, Stack, Typography, Link, Box } from '@mui/material';
// auth
import { useAuthContext } from '../../auth/useAuthContext';
// layouts
import LoginLayout from '../../layouts/login';
//
import AuthLoginForm from './AuthLoginForm';
import AuthWithSocial from './AuthWithSocial';

// ----------------------------------------------------------------------

export default function Login() {
  const { method } = useAuthContext();

  return (
    <LoginLayout>
      <Stack spacing={2} sx={{ mb: 5, position: 'relative' }}>
        <Typography variant="h4">Sign in to Minimal</Typography>

        <Stack direction="row" spacing={0.5}>
          <Typography variant="body2">New user?</Typography>

          <Link variant="subtitle2">Create an account</Link>
        </Stack>

        <Tooltip title={method} placement="left">
          <Box
            component="img"
            alt={method}
            src={`/assets/icons/auth/ic_${method}.png`}
            sx={{ width: 32, height: 32, position: 'absolute', right: 0 }}
          />
        </Tooltip>
      </Stack>

      <Alert severity="info" sx={{ mb: 3 }}>
        <Typography variant="body2" sx={{ mb: 1 }}>
          <strong>Comptes de test disponibles :</strong>
        </Typography>
        <Typography variant="body2" sx={{ mb: 0.5 }}>
          • Admin: <strong>demo@minimals.cc</strong> / <strong>demo1234</strong>
        </Typography>
        <Typography variant="body2" sx={{ mb: 0.5 }}>
          • User: <strong>john@example.com</strong> / <strong>password123</strong>
        </Typography>
        <Typography variant="body2">
          • User: <strong>jane@example.com</strong> / <strong>password456</strong>
        </Typography>
      </Alert>

      <AuthLoginForm />

      <AuthWithSocial />
    </LoginLayout>
  );
}
