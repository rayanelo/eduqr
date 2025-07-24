import { useState, useEffect } from 'react';
import * as Yup from 'yup';
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
  Stepper,
  Step,
  StepLabel,
  Box,
  Typography,
  Alert,
} from '@mui/material';

// components
import FormProvider, { RHFTextField, RHFCheckbox } from '../../components/hook-form';
import { useSnackbar } from '../../components/snackbar';

// hooks
import { useRooms } from '../../hooks/useRooms';

// ----------------------------------------------------------------------

const RoomSchema = Yup.object().shape({
  name: Yup.string().required('Le nom de la salle est requis'),
  building: Yup.string(),
  floor: Yup.string(),
  is_modular: Yup.boolean(),
  sub_rooms_count: Yup.number().when('is_modular', {
    is: true,
    then: (schema) => schema.min(2, 'Au moins 2 sous-salles requises').max(10, 'Maximum 10 sous-salles'),
    otherwise: (schema) => schema.optional(),
  }),
});

const steps = ['Informations de base', 'Configuration modulaire'];

export default function RoomFormDialog({ open, onClose, onSuccess, room, isEdit }) {
  const { enqueueSnackbar } = useSnackbar();
  const { createRoom, updateRoom, isLoading } = useRooms();
  const [activeStep, setActiveStep] = useState(0);
  const [isModular, setIsModular] = useState(false);

  const methods = useForm({
    resolver: yupResolver(RoomSchema),
    defaultValues: {
      name: '',
      building: '',
      floor: '',
      is_modular: false,
      sub_rooms_count: 2,
    },
  });

  const { reset, handleSubmit, watch } = methods;

  // Watch for modular checkbox changes
  const watchedIsModular = watch('is_modular');

  useEffect(() => {
    setIsModular(watchedIsModular);
  }, [watchedIsModular]);

  // Reset form when dialog opens/closes or room changes
  useEffect(() => {
    if (open) {
      if (room && isEdit) {
        reset({
          name: room.name,
          building: room.building || '',
          floor: room.floor || '',
          is_modular: room.is_modular,
          sub_rooms_count: room.children?.length || 2,
        });
        setIsModular(room.is_modular);
      } else {
        reset({
          name: '',
          building: '',
          floor: '',
          is_modular: false,
          sub_rooms_count: 2,
        });
        setIsModular(false);
      }
      setActiveStep(0);
    }
  }, [open, room, isEdit, reset]);

  const handleNext = () => {
    setActiveStep((prevStep) => prevStep + 1);
  };

  const handleBack = () => {
    setActiveStep((prevStep) => prevStep - 1);
  };

  const handleClose = () => {
    setActiveStep(0);
    onClose();
  };

  const onSubmit = async (data) => {
    try {
      if (isEdit) {
        await updateRoom(room.id, {
          name: data.name,
          building: data.building,
          floor: data.floor,
          is_modular: data.is_modular,
        });
      } else {
        await createRoom({
          name: data.name,
          building: data.building,
          floor: data.floor,
          is_modular: data.is_modular,
          sub_rooms_count: data.is_modular ? data.sub_rooms_count : undefined,
        });
      }
      onSuccess();
    } catch (error) {
      enqueueSnackbar(error.message || 'Erreur lors de l\'opération', { variant: 'error' });
    }
  };

  const renderStepContent = (step) => {
    switch (step) {
      case 0:
        return (
          <Stack spacing={3}>
            <Typography variant="h6" sx={{ mb: 2 }}>
              Informations de base
            </Typography>

            <RHFTextField
              name="name"
              label="Nom de la salle *"
              placeholder="Ex: Salle 101, Amphi A, etc."
            />

            <RHFTextField
              name="building"
              label="Bâtiment"
              placeholder="Ex: Bâtiment A, Sciences, etc."
            />

            <RHFTextField
              name="floor"
              label="Étage"
              placeholder="Ex: RDC, 1er étage, etc."
            />

            <RHFCheckbox
              name="is_modular"
              label="Salle modulable"
              helperText="Une salle modulable peut être divisée en plusieurs sous-salles"
            />

            {isModular && (
              <Alert severity="info">
                Cette salle sera configurée comme modulable. Vous pourrez définir le nombre de sous-salles à l'étape suivante.
              </Alert>
            )}
          </Stack>
        );

      case 1:
        return (
          <Stack spacing={3}>
            <Typography variant="h6" sx={{ mb: 2 }}>
              Configuration modulaire
            </Typography>

            <Typography variant="body2" color="text.secondary">
              Définissez le nombre de sous-salles à créer pour "{watch('name')}".
              Les sous-salles seront nommées automatiquement : {watch('name')} A, {watch('name')} B, etc.
            </Typography>

            <RHFTextField
              name="sub_rooms_count"
              label="Nombre de sous-salles *"
              type="number"
              inputProps={{ min: 2, max: 10 }}
              helperText="Entre 2 et 10 sous-salles"
            />

            <Alert severity="info">
              <Typography variant="body2">
                <strong>Sous-salles qui seront créées :</strong>
                <br />
                {Array.from({ length: Math.min(watch('sub_rooms_count') || 2, 10) }, (_, i) => 
                  `${watch('name')} ${String.fromCharCode(65 + i)}`
                ).join(', ')}
              </Typography>
            </Alert>
          </Stack>
        );

      default:
        return null;
    }
  };

  const canProceedToNext = () => {
    const values = methods.getValues();
    if (activeStep === 0) {
      return values.name && values.name.trim() !== '';
    }
    return true;
  };

  const canSubmit = () => {
    const values = methods.getValues();
    if (activeStep === 0) {
      return values.name && values.name.trim() !== '';
    }
    if (activeStep === 1) {
      return values.is_modular ? values.sub_rooms_count >= 2 : true;
    }
    return false;
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
      <DialogTitle>
        {isEdit ? 'Modifier la salle' : 'Nouvelle salle'}
      </DialogTitle>

      <DialogContent>
        <Box sx={{ mt: 2, mb: 3 }}>
          <Stepper activeStep={activeStep} alternativeLabel>
            {steps.map((label) => (
              <Step key={label}>
                <StepLabel>{label}</StepLabel>
              </Step>
            ))}
          </Stepper>
        </Box>

        <FormProvider methods={methods} onSubmit={handleSubmit(onSubmit)}>
          {renderStepContent(activeStep)}
        </FormProvider>
      </DialogContent>

      <DialogActions>
        <Button onClick={handleClose} disabled={isLoading}>
          Annuler
        </Button>

        {activeStep > 0 && (
          <Button onClick={handleBack} disabled={isLoading}>
            Retour
          </Button>
        )}

        {activeStep < steps.length - 1 ? (
          <Button
            variant="contained"
            onClick={handleNext}
            disabled={!canProceedToNext() || isLoading}
          >
            Suivant
          </Button>
        ) : (
          <Button
            variant="contained"
            onClick={handleSubmit(onSubmit)}
            disabled={!canSubmit() || isLoading}
          >
            {isLoading ? 'Enregistrement...' : (isEdit ? 'Mettre à jour' : 'Créer')}
          </Button>
        )}
      </DialogActions>
    </Dialog>
  );
} 