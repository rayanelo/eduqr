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
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Paper,
  Fade,
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

// Critères de validation du mot de passe
const PASSWORD_CRITERIA = [
  {
    key: 'length',
    label: 'Au moins 8 caractères',
    test: (password) => password.length >= 8,
    icon: 'eva:checkmark-circle-2-fill',
  },
  {
    key: 'uppercase',
    label: 'Au moins 1 lettre majuscule',
    test: (password) => /[A-Z]/.test(password),
    icon: 'eva:checkmark-circle-2-fill',
  },
  {
    key: 'lowercase',
    label: 'Au moins 1 lettre minuscule',
    test: (password) => /[a-z]/.test(password),
    icon: 'eva:checkmark-circle-2-fill',
  },
  {
    key: 'number',
    label: 'Au moins 1 chiffre',
    test: (password) => /[0-9]/.test(password),
    icon: 'eva:checkmark-circle-2-fill',
  },
  {
    key: 'special',
    label: 'Au moins 1 caractère spécial (!@#$%^&*)',
    test: (password) => /[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]/.test(password),
    icon: 'eva:checkmark-circle-2-fill',
  },
];

export default function PasswordTab({ onChangePassword, onValidatePassword, isLoading }) {
  const [showCurrentPassword, setShowCurrentPassword] = useState(false);
  const [showNewPassword, setShowNewPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [passwordStrength, setPasswordStrength] = useState(null);
  const [isValidating, setIsValidating] = useState(false);
  const [criteriaStatus, setCriteriaStatus] = useState({});

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
  const confirmPassword = watch('confirm_password');

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
    let strengthLevel = '';
    
    switch (score) {
      case 0:
        feedback = 'Mot de passe très faible - Ajoutez plus de caractères et de variété';
        strengthLevel = 'Très faible';
        break;
      case 1:
        feedback = 'Mot de passe faible - Ajoutez des lettres majuscules, minuscules, chiffres ou caractères spéciaux';
        strengthLevel = 'Faible';
        break;
      case 2:
        feedback = 'Mot de passe moyen - Ajoutez plus de variété de caractères';
        strengthLevel = 'Moyen';
        break;
      case 3:
        feedback = 'Mot de passe bon - Presque parfait !';
        strengthLevel = 'Bon';
        break;
      case 4:
        feedback = 'Mot de passe fort - Excellent niveau de sécurité';
        strengthLevel = 'Fort';
        break;
      case 5:
        feedback = 'Mot de passe très fort - Niveau de sécurité optimal';
        strengthLevel = 'Très fort';
        break;
      default:
        feedback = 'Mot de passe invalide';
        strengthLevel = 'Invalide';
        break;
    }
    
    return {
      score,
      feedback,
      strengthLevel,
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
        
        // Mettre à jour le statut des critères
        setCriteriaStatus(clientStrength.criteria);
        
        // Validation côté serveur après un délai (pour les mots de passe plus longs)
        if (newPassword.length >= 3) {
          setIsValidating(true);
          try {
            const serverStrength = await onValidatePassword(newPassword);
            // Ajouter le strengthLevel si il n'existe pas
            if (!serverStrength.strengthLevel) {
              serverStrength.strengthLevel = getStrengthLevelFromScore(serverStrength.score);
            }
            setPasswordStrength(serverStrength);
            setCriteriaStatus(serverStrength.criteria);
          } catch (error) {
            console.error('Error validating password:', error);
            // Garder la validation côté client en cas d'erreur serveur
          } finally {
            setIsValidating(false);
          }
        }
      } else {
        setPasswordStrength(null);
        setCriteriaStatus({});
      }
    };

    const timeoutId = setTimeout(validatePassword, 200);
    return () => clearTimeout(timeoutId);
  }, [newPassword, onValidatePassword]);

  // Fonction pour obtenir le niveau de force à partir du score
  const getStrengthLevelFromScore = (score) => {
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
        return 'Invalide';
    }
  };

  const onSubmit = async (data) => {
    try {
      await onChangePassword(data);
      reset();
      setPasswordStrength(null);
      setCriteriaStatus({});
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

  const getStrengthIcon = (score) => {
    switch (score) {
      case 0:
      case 1:
        return 'eva:alert-triangle-fill';
      case 2:
        return 'eva:alert-circle-fill';
      case 3:
        return 'eva:checkmark-circle-fill';
      case 4:
      case 5:
        return 'eva:shield-fill';
      default:
        return 'eva:help-circle-fill';
    }
  };

  const isPasswordMatch = newPassword && confirmPassword && newPassword === confirmPassword;

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
            <Fade in={true} timeout={300}>
              <Paper elevation={1} sx={{ p: 2, bgcolor: 'background.neutral' }}>
                <Stack spacing={2}>
                  {/* Header avec score et feedback */}
                  <Stack direction="row" alignItems="center" spacing={2}>
                    <Iconify 
                      icon={getStrengthIcon(passwordStrength.score)} 
                      color={getStrengthColor(passwordStrength.score)}
                      width={24}
                    />
                    <Box sx={{ flexGrow: 1 }}>
                      <Stack direction="row" alignItems="center" spacing={1}>
                        <Typography variant="body2" color="text.secondary">
                          Force du mot de passe:
                        </Typography>
                        <Chip
                          label={passwordStrength.strengthLevel || getStrengthLevelFromScore(passwordStrength.score)}
                          color={getStrengthColor(passwordStrength.score)}
                          size="small"
                          variant="filled"
                        />
                        {isValidating && (
                          <Typography variant="caption" color="text.secondary">
                            Validation serveur...
                          </Typography>
                        )}
                      </Stack>
                    </Box>
                  </Stack>

                  {/* Barre de progression */}
                  <Box>
                    <LinearProgress
                      variant="determinate"
                      value={(passwordStrength.score / 5) * 100}
                      color={getStrengthColor(passwordStrength.score)}
                      sx={{ 
                        height: 8, 
                        borderRadius: 4,
                        bgcolor: 'background.paper',
                        '& .MuiLinearProgress-bar': {
                          borderRadius: 4,
                        }
                      }}
                    />
                  </Box>

                  {/* Message de feedback */}
                  <Typography variant="body2" color="text.secondary">
                    {passwordStrength.feedback}
                  </Typography>

                  {/* Critères de validation */}
                  <Box>
                    <Typography variant="subtitle2" color="text.primary" sx={{ mb: 1 }}>
                      Critères de sécurité:
                    </Typography>
                    <List dense sx={{ py: 0 }}>
                      {PASSWORD_CRITERIA.map((criterion) => {
                        const isMet = criteriaStatus[criterion.key];
                        return (
                          <ListItem key={criterion.key} sx={{ py: 0.5, px: 0 }}>
                            <ListItemIcon sx={{ minWidth: 32 }}>
                              <Iconify
                                icon={isMet ? 'eva:checkmark-circle-2-fill' : 'eva:radio-button-off-fill'}
                                color={isMet ? 'success.main' : 'text.disabled'}
                                width={16}
                              />
                            </ListItemIcon>
                            <ListItemText
                              primary={criterion.label}
                              primaryTypographyProps={{
                                variant: 'body2',
                                color: isMet ? 'text.primary' : 'text.disabled',
                                sx: {
                                  textDecoration: isMet ? 'none' : 'none',
                                  fontWeight: isMet ? 500 : 400,
                                }
                              }}
                            />
                          </ListItem>
                        );
                      })}
                    </List>
                  </Box>

                  {/* Alerte si le mot de passe n'est pas assez fort */}
                  {!passwordStrength.is_valid && (
                    <Alert severity="warning" sx={{ mt: 1 }}>
                      Le mot de passe doit respecter au moins 4 critères de sécurité sur 5.
                    </Alert>
                  )}
                </Stack>
              </Paper>
            </Fade>
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

          {/* Indicateur de correspondance des mots de passe */}
          {confirmPassword && (
            <Fade in={true} timeout={200}>
              <Alert 
                severity={isPasswordMatch ? 'success' : 'error'} 
                icon={<Iconify icon={isPasswordMatch ? 'eva:checkmark-circle-2-fill' : 'eva:close-circle-fill'} />}
                sx={{ py: 0.5 }}
              >
                {isPasswordMatch 
                  ? 'Les mots de passe correspondent' 
                  : 'Les mots de passe ne correspondent pas'
                }
              </Alert>
            </Fade>
          )}

          {/* Actions */}
          <Stack direction="row" spacing={2} justifyContent="flex-end">
            <LoadingButton
              type="submit"
              variant="contained"
              loading={isSubmitting || isLoading}
              disabled={
                (passwordStrength && !passwordStrength.is_valid) || 
                !isPasswordMatch ||
                !newPassword ||
                !confirmPassword
              }
              startIcon={<Iconify icon="eva:lock-fill" />}
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