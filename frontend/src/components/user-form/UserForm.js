import PropTypes from 'prop-types';
import { useState, useEffect } from 'react';
import * as Yup from 'yup';
// form
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
// @mui
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Stack,
  Typography,
  InputAdornment,
  IconButton,
  Box,
  LinearProgress,
  Alert,
  Chip,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Paper,
  Fade,
  MenuItem,
} from '@mui/material';
import { LoadingButton } from '@mui/lab';
// components
import Iconify from '../iconify';
import FormProvider, { RHFTextField, RHFSelect } from '../hook-form';
import { usePermissions } from '../../hooks/usePermissions';

// ----------------------------------------------------------------------

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

const UserSchema = Yup.object().shape({
  first_name: Yup.string().required('Le prénom est requis'),
  last_name: Yup.string().required('Le nom est requis'),
  email: Yup.string().required('L\'email est requis').email('Email invalide'),
  phone: Yup.string().required('Le numéro de téléphone est requis'),
  address: Yup.string().required('L\'adresse est requise'),
  password: Yup.string().when('$isEdit', {
    is: false,
    then: Yup.string().required('Le mot de passe est requis').min(8, 'Minimum 8 caractères'),
    otherwise: Yup.string().optional(),
  }),
  confirm_password: Yup.string().when('$isEdit', {
    is: false,
    then: Yup.string().required('La confirmation du mot de passe est requise').oneOf([Yup.ref('password'), null], 'Les mots de passe doivent correspondre'),
    otherwise: Yup.string().optional(),
  }),
  role: Yup.string().required('Le rôle est requis'),
});

// ----------------------------------------------------------------------

UserForm.propTypes = {
  open: PropTypes.bool,
  onClose: PropTypes.func,
  onSubmit: PropTypes.func,
  user: PropTypes.object,
  isEdit: PropTypes.bool,
  isLoading: PropTypes.bool,
};

export default function UserForm({ open, onClose, onSubmit, user, isEdit = false, isLoading = false }) {
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [passwordStrength, setPasswordStrength] = useState(null);
  const [criteriaStatus, setCriteriaStatus] = useState({});
  const { getCreatableRoles } = usePermissions();
  const creatableRoles = getCreatableRoles();
  const roleOptions = getRoleOptions(creatableRoles);

  const defaultValues = {
    first_name: user?.first_name || '',
    last_name: user?.last_name || '',
    email: user?.email || '',
    phone: user?.phone || '',
    address: user?.address || '',
    password: '',
    confirm_password: '',
    role: user?.role || 'etudiant',
  };

  const methods = useForm({
    resolver: yupResolver(UserSchema),
    defaultValues,
    context: { isEdit },
  });

  const {
    reset,
    handleSubmit,
    watch,
    formState: { isSubmitting },
  } = methods;

  const password = watch('password');
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
    if (password && password.length >= 1) {
      const clientStrength = validatePasswordClient(password);
      setPasswordStrength(clientStrength);
      setCriteriaStatus(clientStrength.criteria);
    } else {
      setPasswordStrength(null);
      setCriteriaStatus({});
    }
  }, [password]);

  const handleClose = () => {
    reset();
    setPasswordStrength(null);
    setCriteriaStatus({});
    onClose();
  };

  const handleSubmitForm = async (data) => {
    try {
      await onSubmit(data);
      handleClose();
    } catch (error) {
      console.error('Error submitting form:', error);
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

  const isPasswordMatch = password && confirmPassword && password === confirmPassword;

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

  return (
    <Dialog fullWidth maxWidth="sm" open={open} onClose={handleClose}>
      <DialogTitle>
        <Typography variant="h6">
          {isEdit ? 'Modifier l\'utilisateur' : 'Créer un utilisateur'}
        </Typography>
      </DialogTitle>

      <FormProvider methods={methods} onSubmit={handleSubmit(handleSubmitForm)}>
        <DialogContent>
          <Stack spacing={3}>
            <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
              <RHFTextField name="first_name" label="Prénom" />
              <RHFTextField name="last_name" label="Nom" />
            </Stack>

            <RHFTextField name="email" label="Email" />

            <RHFTextField name="phone" label="Numéro de téléphone" />

            <RHFTextField name="address" label="Adresse" multiline rows={2} />

            {/* Mot de passe avec validation en temps réel */}
            <RHFTextField
              name="password"
              label={isEdit ? 'Nouveau mot de passe (optionnel)' : 'Mot de passe'}
              type={showPassword ? 'text' : 'password'}
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton onClick={() => setShowPassword(!showPassword)} edge="end">
                      <Iconify icon={showPassword ? 'eva:eye-fill' : 'eva:eye-off-fill'} />
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />

            {/* Indicateur de force du mot de passe (seulement pour la création) */}
            {!isEdit && password && passwordStrength && (
              <Fade in={true} timeout={300}>
                <Paper elevation={1} sx={{ p: 2, bgcolor: 'background.neutral' }}>
                  <Stack spacing={2}>
                    {/* Header avec score et feedback */}
                    <Stack direction="row" alignItems="center" spacing={2}>
                      <Iconify 
                        icon={getStrengthIcon(passwordStrength.score)} 
                        color={getStrengthColor(passwordStrength.score)}
                        width={20}
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
                          height: 6, 
                          borderRadius: 3,
                          bgcolor: 'background.paper',
                          '& .MuiLinearProgress-bar': {
                            borderRadius: 3,
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
                              <ListItemIcon sx={{ minWidth: 28 }}>
                                <Iconify
                                  icon={isMet ? 'eva:checkmark-circle-2-fill' : 'eva:radio-button-off-fill'}
                                  color={isMet ? 'success.main' : 'text.disabled'}
                                  width={14}
                                />
                              </ListItemIcon>
                              <ListItemText
                                primary={criterion.label}
                                primaryTypographyProps={{
                                  variant: 'body2',
                                  color: isMet ? 'text.primary' : 'text.disabled',
                                  sx: {
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

            {/* Confirmation du mot de passe (seulement pour la création) */}
            {!isEdit && (
              <RHFTextField
                name="confirm_password"
                label="Confirmation du mot de passe"
                type={showConfirmPassword ? 'text' : 'password'}
                InputProps={{
                  endAdornment: (
                    <InputAdornment position="end">
                      <IconButton onClick={() => setShowConfirmPassword(!showConfirmPassword)} edge="end">
                        <Iconify icon={showConfirmPassword ? 'eva:eye-fill' : 'eva:eye-off-fill'} />
                      </IconButton>
                    </InputAdornment>
                  ),
                }}
              />
            )}

            {/* Indicateur de correspondance des mots de passe (seulement pour la création) */}
            {!isEdit && confirmPassword && (
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

            <RHFSelect name="role" label="Rôle">
              {roleOptions.map((option) => (
                <MenuItem key={option.value} value={option.value}>
                  {option.label}
                </MenuItem>
              ))}
            </RHFSelect>
          </Stack>
        </DialogContent>

        <DialogActions>
          <Button onClick={handleClose} color="inherit">
            Annuler
          </Button>
          <LoadingButton 
            type="submit" 
            variant="contained" 
            loading={isSubmitting || isLoading}
            disabled={
              !isEdit && (
                (passwordStrength && !passwordStrength.is_valid) || 
                !isPasswordMatch ||
                !password ||
                !confirmPassword
              )
            }
            startIcon={<Iconify icon={isEdit ? 'eva:edit-fill' : 'eva:plus-fill'} />}
          >
            {isEdit ? 'Modifier' : 'Créer'}
          </LoadingButton>
        </DialogActions>
      </FormProvider>
    </Dialog>
  );
}

// Helper function to get role options
function getRoleOptions(creatableRoles) {
  const roleLabels = {
    super_admin: 'Super Administrateur',
    admin: 'Administrateur',
    professeur: 'Professeur',
    etudiant: 'Étudiant',
  };

  return creatableRoles.map(role => ({
    value: role,
    label: roleLabels[role] || role,
  }));
} 