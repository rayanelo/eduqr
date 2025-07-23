import PropTypes from 'prop-types';
import { forwardRef } from 'react';
// @mui
import { TextField } from '@mui/material';

// ----------------------------------------------------------------------

const DateTimePicker = forwardRef(({ label, value, onChange, error, helperText, size = 'small', ...other }, ref) => {
  const handleChange = (event) => {
    if (onChange) {
      onChange(event.target.value);
    }
  };

  return (
    <TextField
      ref={ref}
      type="datetime-local"
      label={label}
      value={value || ''}
      onChange={handleChange}
      error={error}
      helperText={helperText}
      size={size}
      fullWidth
      InputLabelProps={{
        shrink: true,
      }}
      {...other}
    />
  );
});

DateTimePicker.propTypes = {
  label: PropTypes.string,
  value: PropTypes.string,
  onChange: PropTypes.func,
  error: PropTypes.bool,
  helperText: PropTypes.string,
  size: PropTypes.oneOf(['small', 'medium']),
};

export default DateTimePicker; 