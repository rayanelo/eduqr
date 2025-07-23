// @mui
import { Stack, Typography, Link } from '@mui/material';
// layouts
import LoginLayout from '../layouts/login';
//
import AuthRegisterForm from '../sections/auth/AuthRegisterForm';

// ----------------------------------------------------------------------

export default function RegisterPage() {
  return (
    <LoginLayout>
      <Stack spacing={2} sx={{ mb: 5, position: 'relative' }}>
        <Typography variant="h4">Create an account</Typography>

        <Stack direction="row" spacing={0.5}>
          <Typography variant="body2">Already have an account?</Typography>

          <Link href="/login" variant="subtitle2">
            Sign in
          </Link>
        </Stack>
      </Stack>

      <AuthRegisterForm />
    </LoginLayout>
  );
} 