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
  Typography,
  Divider,
  TextField,
} from '@mui/material';
import { LoadingButton } from '@mui/lab';
// components
import FormProvider, { RHFTextField } from '../../components/hook-form';

// ----------------------------------------------------------------------

const ProfileSchema = Yup.object().shape({
  contact_email: Yup.string().email('Email de contact invalide'),
  phone: Yup.string().required('Le téléphone est requis'),
  address: Yup.string().required('L\'adresse est requise'),
});

export default function ProfileInfoTab({ user, onUpdate, isLoading }) {
  const [isEditing, setIsEditing] = useState(false);

  const defaultValues = {
    contact_email: user?.contact_email || user?.email || '',
    phone: user?.phone || '',
    address: user?.address || '',
  };

  const methods = useForm({
    resolver: yupResolver(ProfileSchema),
    defaultValues,
  });

  const {
    reset,
    handleSubmit,
    formState: { isSubmitting },
  } = methods;

  const handleEdit = () => {
    setIsEditing(true);
  };

  const handleCancel = () => {
    setIsEditing(false);
    reset();
  };

  const onSubmit = async (data) => {
    try {
      await onUpdate(data);
      setIsEditing(false);
    } catch (error) {
      console.error('Error updating profile:', error);
    }
  };

  return (
    
      <Stack sx={{ p: 3 }} spacing={3}>
        {/* Header */}
        <Stack direction="row" alignItems="center" justifyContent="space-between">
          <Typography variant="h6">Informations personnelles</Typography>
          {!isEditing && (
            <Button variant="outlined" onClick={handleEdit}>
              Modifier
            </Button>
          )}
        </Stack>

        <Divider />

        {/* Informations en lecture seule */}
        <Stack spacing={2}>
          <Typography variant="subtitle2" color="text.secondary">
            Informations du compte (non modifiables)
          </Typography>
          
          <Stack direction={{ xs: 'column', md: 'row' }} spacing={2}>
            <TextField
              label="Prénom"
              value={user?.first_name || ''}
              disabled
              fullWidth
            />
            <TextField
              label="Nom"
              value={user?.last_name || ''}
              disabled
              fullWidth
            />
          </Stack>

          <TextField
            label="Email du compte"
            value={user?.email || ''}
            disabled
            fullWidth
            helperText="Email utilisé pour la connexion"
          />
        </Stack>

        <Divider />

        {/* Form */}
        <FormProvider methods={methods} onSubmit={handleSubmit(onSubmit)}>
          <Stack spacing={3}>
            <Typography variant="subtitle2" color="text.secondary">
              Informations modifiables
            </Typography>

            {/* Email de contact */}
            <RHFTextField
              name="contact_email"
              label="Email de contact"
              disabled={!isEditing}
              fullWidth
              helperText="Email pour les communications (par défaut: email du compte)"
            />

            {/* Téléphone et adresse */}
            <Stack direction={{ xs: 'column', md: 'row' }} spacing={2}>
              <RHFTextField
                name="phone"
                label="Téléphone"
                disabled={!isEditing}
                fullWidth
              />
              <RHFTextField
                name="address"
                label="Adresse postale"
                disabled={!isEditing}
                fullWidth
                multiline
                rows={2}
              />
            </Stack>

            {/* Actions */}
            {isEditing && (
              <Stack direction="row" spacing={2} justifyContent="flex-end">
                <Button variant="outlined" onClick={handleCancel}>
                  Annuler
                </Button>
                <LoadingButton
                  type="submit"
                  variant="contained"
                  loading={isSubmitting || isLoading}
                >
                  Enregistrer
                </LoadingButton>
              </Stack>
            )}
          </Stack>
        </FormProvider>
      </Stack>
  );
}

ProfileInfoTab.propTypes = {
  user: PropTypes.object,
  onUpdate: PropTypes.func,
  isLoading: PropTypes.bool,
}; 