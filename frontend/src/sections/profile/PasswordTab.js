import PropTypes from 'prop-types';
import { useState, useEffect } from 'react';
import * as Yup from 'yup';
// form
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
// @mui
import {
  Stack,
  Typography,
  Box,
  Divider,
  LinearProgress,
  Alert,
  Chip,
  InputAdornment,
  IconButton,
} from '@mui/material';
import { LoadingButton } from '@mui/lab';
// components
import Iconify from '../../components/iconify';
import FormProvider, { RHFTextField } from '../../components/hook-form';

// ----------------------------------------------------------------------

const PasswordSchema = Yup.object().shape({
  current_password: Yup.string().required('Le mot de passe actuel est requis'),
  new_password: Yup.string()
    .min(8, 'Le mot de passe doit contenir au moins 8 caractères')
    .required('Le nouveau mot de passe est requis'),
  confirm_password: Yup.string()
    .oneOf([Yup.ref('new_password'), null], 'Les mots de passe doivent être identiques')
    .required('La confirmation du mot de passe est requise'),
});

export default function PasswordTab({ onChangePassword, onValidatePassword, isLoading }) {
  const [showCurrentPassword, setShowCurrentPassword] = useState(false);
  const [showNewPassword, setShowNewPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [passwordStrength, setPasswordStrength] = useState(null);
  const [isValidating, setIsValidating] = useState(false);

  const defaultValues = {
    current_password: '',
    new_password: '',
    confirm_password: '',
  };

  const methods = useForm({
    resolver: yupResolver(PasswordSchema),
    defaultValues,
  });

  const {
    reset,
    handleSubmit,
    watch,
    formState: { isSubmitting },
  } = methods;

  const newPassword = watch('new_password');

  // Validation côté client en temps réel
  const validatePasswordClient = (password) => {
    if (!password) return null;
    
    const criteria = {
      length: password.length >= 8,
      uppercase: /[A-Z]/.test(password),
      lowercase: /[a-z]/.test(password),
      number: /[0-9]/.test(password),
      special: /[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]/.test(password),
    };
    
    const score = Object.values(criteria).filter(Boolean).length;
    
    let feedback = '';
    switch (score) {
      case 0:
        feedback = 'Mot de passe très faible';
        break;
      case 1:
        feedback = 'Mot de passe faible';
        break;
      case 2:
        feedback = 'Mot de passe moyen';
        break;
      case 3:
        feedback = 'Mot de passe bon';
        break;
      case 4:
        feedback = 'Mot de passe fort';
        break;
      case 5:
        feedback = 'Mot de passe très fort';
        break;
      default:
        feedback = 'Mot de passe invalide';
        break;
    }
    
    return {
      score,
      feedback,
      is_valid: score >= 4,
      criteria,
    };
  };

  // Validation en temps réel du mot de passe
  useEffect(() => {
    const validatePassword = async () => {
      if (newPassword && newPassword.length >= 1) {
        // Validation côté client immédiate
        const clientStrength = validatePasswordClient(newPassword);
        setPasswordStrength(clientStrength);
        
        // Validation côté serveur après un délai (pour les mots de passe plus longs)
        if (newPassword.length >= 3) {
          setIsValidating(true);
          try {
            const serverStrength = await onValidatePassword(newPassword);
            setPasswordStrength(serverStrength);
          } catch (error) {
            console.error('Error validating password:', error);
            // Garder la validation côté client en cas d'erreur serveur
          } finally {
            setIsValidating(false);
          }
        }
      } else {
        setPasswordStrength(null);
      }
    };

    const timeoutId = setTimeout(validatePassword, 200); // Réduit de 500ms à 200ms
    return () => clearTimeout(timeoutId);
  }, [newPassword, onValidatePassword]);

  const onSubmit = async (data) => {
    try {
      await onChangePassword(data);
      reset();
      setPasswordStrength(null);
    } catch (error) {
      console.error('Error changing password:', error);
    }
  };

  const getStrengthColor = (score) => {
    switch (score) {
      case 0:
      case 1:
        return 'error';
      case 2:
        return 'warning';
      case 3:
        return 'info';
      case 4:
      case 5:
        return 'success';
      default:
        return 'default';
    }
  };

  const getStrengthText = (score) => {
    switch (score) {
      case 0:
        return 'Très faible';
      case 1:
        return 'Faible';
      case 2:
        return 'Moyen';
      case 3:
        return 'Bon';
      case 4:
        return 'Fort';
      case 5:
        return 'Très fort';
      default:
        return '';
    }
  };

  return (
      <Stack sx={{ p: 3 }} spacing={3}>
        {/* Header */}
        <Typography variant="h6">Sécurité</Typography>

        <Divider />

        {/* Form */}
        <FormProvider methods={methods} onSubmit={handleSubmit(onSubmit)}>
          <Stack spacing={3}>
            {/* Mot de passe actuel */}
            <RHFTextField
              name="current_password"
              label="Mot de passe actuel"
              type={showCurrentPassword ? 'text' : 'password'}
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton
                      onClick={() => setShowCurrentPassword(!showCurrentPassword)}
                      edge="end"
                    >
                      <Iconify icon={showCurrentPassword ? 'eva:eye-fill' : 'eva:eye-off-fill'} />
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />

            {/* Nouveau mot de passe */}
            <RHFTextField
              name="new_password"
              label="Nouveau mot de passe"
              type={showNewPassword ? 'text' : 'password'}
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton
                      onClick={() => setShowNewPassword(!showNewPassword)}
                      edge="end"
                    >
                      <Iconify icon={showNewPassword ? 'eva:eye-fill' : 'eva:eye-off-fill'} />
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />

            {/* Indicateur de validation en cours */}
            {newPassword && newPassword.length < 3 && (
              <Typography variant="caption" color="text.secondary">
                Continuez à taper pour voir la validation en temps réel...
              </Typography>
            )}

            {/* Indicateur de force du mot de passe */}
            {newPassword && passwordStrength && (
              <Box>
                <Stack direction="row" alignItems="center" spacing={1} sx={{ mb: 1 }}>
                  <Typography variant="body2" color="text.secondary">
                    Force du mot de passe:
                  </Typography>
                  <Chip
                    label={getStrengthText(passwordStrength.score)}
                    color={getStrengthColor(passwordStrength.score)}
                    size="small"
                  />
                  {isValidating && (
                    <Typography variant="caption" color="text.secondary">
                      Validation serveur...
                    </Typography>
                  )}
                </Stack>

                <LinearProgress
                  variant="determinate"
                  value={(passwordStrength.score / 5) * 100}
                  color={getStrengthColor(passwordStrength.score)}
                  sx={{ mb: 1 }}
                />

                <Typography variant="caption" color="text.secondary">
                  {passwordStrength.feedback}
                </Typography>

                {/* Critères */}
                <Stack spacing={0.5} sx={{ mt: 2 }}>
                  <Typography variant="caption" color="text.secondary">
                    Critères de sécurité:
                  </Typography>
                  <Stack direction="row" flexWrap="wrap" gap={1}>
                    <Chip
                      label="8+ caractères"
                      size="small"
                      variant={passwordStrength.criteria.length ? 'filled' : 'outlined'}
                      color={passwordStrength.criteria.length ? 'success' : 'default'}
                    />
                    <Chip
                      label="Majuscule"
                      size="small"
                      variant={passwordStrength.criteria.uppercase ? 'filled' : 'outlined'}
                      color={passwordStrength.criteria.uppercase ? 'success' : 'default'}
                    />
                    <Chip
                      label="Minuscule"
                      size="small"
                      variant={passwordStrength.criteria.lowercase ? 'filled' : 'outlined'}
                      color={passwordStrength.criteria.lowercase ? 'success' : 'default'}
                    />
                    <Chip
                      label="Chiffre"
                      size="small"
                      variant={passwordStrength.criteria.number ? 'filled' : 'outlined'}
                      color={passwordStrength.criteria.number ? 'success' : 'default'}
                    />
                    <Chip
                      label="Caractère spécial"
                      size="small"
                      variant={passwordStrength.criteria.special ? 'filled' : 'outlined'}
                      color={passwordStrength.criteria.special ? 'success' : 'default'}
                    />
                  </Stack>
                </Stack>

                {!passwordStrength.is_valid && (
                  <Alert severity="warning" sx={{ mt: 2 }}>
                    Le mot de passe doit respecter au moins 4 critères de sécurité sur 5.
                  </Alert>
                )}
              </Box>
            )}

            {/* Confirmation du mot de passe */}
            <RHFTextField
              name="confirm_password"
              label="Confirmer le nouveau mot de passe"
              type={showConfirmPassword ? 'text' : 'password'}
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton
                      onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                      edge="end"
                    >
                      <Iconify icon={showConfirmPassword ? 'eva:eye-fill' : 'eva:eye-off-fill'} />
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />

            {/* Actions */}
            <Stack direction="row" spacing={2} justifyContent="flex-end">
              <LoadingButton
                type="submit"
                variant="contained"
                loading={isSubmitting || isLoading}
                disabled={passwordStrength && !passwordStrength.is_valid}
              >
                Changer le mot de passe
              </LoadingButton>
            </Stack>
          </Stack>
        </FormProvider>
      </Stack>
  
  );
}

PasswordTab.propTypes = {
  onChangePassword: PropTypes.func,
  onValidatePassword: PropTypes.func,
  isLoading: PropTypes.bool,
}; 