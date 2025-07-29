import PropTypes from 'prop-types';
import {
  Stack,
  Typography,
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
} from '@mui/material';
import Iconify from '../iconify';

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

// Validation côté client en temps réel
export const validatePasswordClient = (password) => {
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

export default function PasswordStrengthIndicator({ 
  password, 
  passwordStrength, 
  isValidating = false, 
  compact = false,
  showCriteria = true,
  showAlert = true,
}) {
  if (!password || !passwordStrength) {
    return null;
  }

  return (
    <Fade in={true} timeout={300}>
      <Paper elevation={1} sx={{ p: compact ? 1.5 : 2, bgcolor: 'background.neutral' }}>
        <Stack spacing={compact ? 1.5 : 2}>
          {/* Header avec score et feedback */}
          <Stack direction="row" alignItems="center" spacing={2}>
            <Iconify 
              icon={getStrengthIcon(passwordStrength.score)} 
              color={getStrengthColor(passwordStrength.score)}
              width={compact ? 18 : 24}
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
                height: compact ? 6 : 8, 
                borderRadius: compact ? 3 : 4,
                bgcolor: 'background.paper',
                '& .MuiLinearProgress-bar': {
                  borderRadius: compact ? 3 : 4,
                }
              }}
            />
          </Box>

          {/* Message de feedback */}
          <Typography variant="body2" color="text.secondary">
            {passwordStrength.feedback}
          </Typography>

          {/* Critères de validation */}
          {showCriteria && (
            <Box>
              <Typography variant="subtitle2" color="text.primary" sx={{ mb: 1 }}>
                Critères de sécurité:
              </Typography>
              <List dense sx={{ py: 0 }}>
                {PASSWORD_CRITERIA.map((criterion) => {
                  const isMet = passwordStrength.criteria[criterion.key];
                  return (
                    <ListItem key={criterion.key} sx={{ py: 0.5, px: 0 }}>
                      <ListItemIcon sx={{ minWidth: compact ? 24 : 32 }}>
                        <Iconify
                          icon={isMet ? 'eva:checkmark-circle-2-fill' : 'eva:radio-button-off-fill'}
                          color={isMet ? 'success.main' : 'text.disabled'}
                          width={compact ? 12 : 16}
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
          )}

          {/* Alerte si le mot de passe n'est pas assez fort */}
          {showAlert && !passwordStrength.is_valid && (
            <Alert severity="warning" sx={{ mt: 1 }}>
              Le mot de passe doit respecter au moins 4 critères de sécurité sur 5.
            </Alert>
          )}
        </Stack>
      </Paper>
    </Fade>
  );
}

PasswordStrengthIndicator.propTypes = {
  password: PropTypes.string,
  passwordStrength: PropTypes.object,
  isValidating: PropTypes.bool,
  compact: PropTypes.bool,
  showCriteria: PropTypes.bool,
  showAlert: PropTypes.bool,
}; 