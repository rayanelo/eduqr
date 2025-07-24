import { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Stack,
  Alert,
  CircularProgress,
} from '@mui/material';
import { useForm, Controller } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import { useSubjects } from '../../hooks/useSubjects';

const schema = yup.object().shape({
  name: yup.string().required('Le nom de la matière est requis'),
  code: yup.string().max(20, 'Le code ne peut pas dépasser 20 caractères'),
  description: yup.string().max(500, 'La description ne peut pas dépasser 500 caractères'),
});

export default function SubjectFormDialog({ open, onClose, subject, onSuccess }) {
  const { createSubject, updateSubject, error } = useSubjects();
  const [submitError, setSubmitError] = useState(null);

  const {
    control,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm({
    resolver: yupResolver(schema),
    defaultValues: {
      name: '',
      code: '',
      description: '',
    },
  });

  useEffect(() => {
    if (open) {
      setSubmitError(null);
      if (subject) {
        // Mode édition
        reset({
          name: subject.name,
          code: subject.code || '',
          description: subject.description || '',
        });
      } else {
        // Mode création
        reset({
          name: '',
          code: '',
          description: '',
        });
      }
    }
  }, [open, subject, reset]);

  const onSubmit = async (data) => {
    try {
      setSubmitError(null);
      
      if (subject) {
        // Mode édition
        await updateSubject(subject.id, data);
      } else {
        // Mode création
        await createSubject(data);
      }
      
      onSuccess();
    } catch (err) {
      setSubmitError(err.response?.data?.error || 'Une erreur est survenue');
    }
  };

  const handleClose = () => {
    if (!isSubmitting) {
      onClose();
    }
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        {subject ? 'Modifier la matière' : 'Nouvelle matière'}
      </DialogTitle>
      
      <form onSubmit={handleSubmit(onSubmit)}>
        <DialogContent>
          <Stack spacing={3}>
            {(error || submitError) && (
              <Alert severity="error" onClose={() => setSubmitError(null)}>
                {submitError || error}
              </Alert>
            )}

            <Controller
              name="name"
              control={control}
              render={({ field }) => (
                <TextField
                  {...field}
                  label="Nom de la matière *"
                  fullWidth
                  error={!!errors.name}
                  helperText={errors.name?.message}
                  disabled={isSubmitting}
                />
              )}
            />

            <Controller
              name="code"
              control={control}
              render={({ field }) => (
                <TextField
                  {...field}
                  label="Code matière"
                  fullWidth
                  placeholder="Ex: BIO101"
                  error={!!errors.code}
                  helperText={errors.code?.message || 'Facultatif mais utile pour l\'identification'}
                  disabled={isSubmitting}
                />
              )}
            />

            <Controller
              name="description"
              control={control}
              render={({ field }) => (
                <TextField
                  {...field}
                  label="Description"
                  fullWidth
                  multiline
                  rows={3}
                  placeholder="Courte description de la matière..."
                  error={!!errors.description}
                  helperText={errors.description?.message || 'Facultatif'}
                  disabled={isSubmitting}
                />
              )}
            />
          </Stack>
        </DialogContent>

        <DialogActions>
          <Button onClick={handleClose} disabled={isSubmitting}>
            Annuler
          </Button>
          <Button
            type="submit"
            variant="contained"
            disabled={isSubmitting}
            startIcon={isSubmitting ? <CircularProgress size={20} /> : null}
          >
            {isSubmitting ? 'Enregistrement...' : (subject ? 'Modifier' : 'Créer')}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
} 