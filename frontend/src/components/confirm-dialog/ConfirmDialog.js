import PropTypes from 'prop-types';
// @mui
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
} from '@mui/material';

// ----------------------------------------------------------------------

ConfirmDialog.propTypes = {
  open: PropTypes.bool,
  title: PropTypes.string,
  content: PropTypes.string,
  action: PropTypes.node,
  onClose: PropTypes.func,
  onConfirm: PropTypes.func,
  confirmText: PropTypes.string,
  cancelText: PropTypes.string,
};

export default function ConfirmDialog({ 
  title, 
  content, 
  action, 
  open, 
  onClose, 
  onConfirm,
  confirmText = 'Confirmer',
  cancelText = 'Annuler',
  ...other 
}) {
  return (
    <Dialog fullWidth maxWidth="xs" open={open} onClose={onClose} {...other}>
      <DialogTitle sx={{ pb: 2 }}>
        {title}
      </DialogTitle>

      {content && (
        <DialogContent sx={{ typography: 'body2' }}>
          {content}
        </DialogContent>
      )}

      <DialogActions>
        {action}
        {onConfirm && (
          <Button variant="contained" color="primary" onClick={onConfirm}>
            {confirmText}
          </Button>
        )}
        <Button variant="outlined" color="inherit" onClick={onClose}>
          {cancelText}
        </Button>
      </DialogActions>
    </Dialog>
  );
} 