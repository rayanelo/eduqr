import PropTypes from 'prop-types';
import { Alert } from '@mui/material';
import { Fade } from '@mui/material';
import Iconify from '../iconify';

export default function PasswordMatchIndicator({ password, confirmPassword }) {
  if (!confirmPassword) {
    return null;
  }

  const isMatch = password && confirmPassword && password === confirmPassword;

  return (
    <Fade in={true} timeout={200}>
      <Alert 
        severity={isMatch ? 'success' : 'error'} 
        icon={<Iconify icon={isMatch ? 'eva:checkmark-circle-2-fill' : 'eva:close-circle-fill'} />}
        sx={{ py: 0.5 }}
      >
        {isMatch 
          ? 'Les mots de passe correspondent' 
          : 'Les mots de passe ne correspondent pas'
        }
      </Alert>
    </Fade>
  );
}

PasswordMatchIndicator.propTypes = {
  password: PropTypes.string,
  confirmPassword: PropTypes.string,
}; 