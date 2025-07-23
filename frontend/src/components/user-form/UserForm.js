import PropTypes from 'prop-types';
import { useState } from 'react';
import * as Yup from 'yup';
// form
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
// @mui
import {
  Stack,
  Button,
  Dialog,
  MenuItem,
  Typography,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
  InputAdornment,
} from '@mui/material';
import { LoadingButton } from '@mui/lab';
// components
import Iconify from '../iconify';
import FormProvider, { RHFTextField, RHFSelect } from '../hook-form';
import { usePermissions } from '../../hooks/usePermissions';

// ----------------------------------------------------------------------

const ROLE_OPTIONS = [
  { value: 'super_admin', label: 'Super Administrateur' },
  { value: 'admin', label: 'Administrateur' },
  { value: 'professeur', label: 'Professeur' },
  { value: 'etudiant', label: 'Étudiant' },
];

// Fonction pour obtenir les options de rôle selon les permissions
const getRoleOptions = (creatableRoles) => {
  return ROLE_OPTIONS.filter(option => creatableRoles.includes(option.value));
};



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
    formState: { isSubmitting },
  } = methods;

  const handleClose = () => {
    reset();
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

            {!isEdit && (
              <RHFTextField
                name="confirm_password"
                label="Confirmation du mot de passe"
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
          <LoadingButton type="submit" variant="contained" loading={isSubmitting || isLoading}>
            {isEdit ? 'Modifier' : 'Créer'}
          </LoadingButton>
        </DialogActions>
      </FormProvider>
    </Dialog>
  );
} 