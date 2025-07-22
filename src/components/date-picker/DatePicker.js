import PropTypes from 'prop-types';
import { forwardRef } from 'react';
// @mui
import { TextField } from '@mui/material';

// ----------------------------------------------------------------------

const DatePicker = forwardRef(({ label, value, onChange, error, helperText, size = 'small', ...other }, ref) => {
  const handleChange = (event) => {
    if (onChange) {
      onChange(event.target.value);
    }
  };

  return (
    <TextField
      ref={ref}
      type="date"
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

DatePicker.propTypes = {
  label: PropTypes.string,
  value: PropTypes.string,
  onChange: PropTypes.func,
  error: PropTypes.bool,
  helperText: PropTypes.string,
  size: PropTypes.oneOf(['small', 'medium']),
};

export default DatePicker; 