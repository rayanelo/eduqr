import PropTypes from 'prop-types';
// @mui
import { Box, FormControlLabel, Radio, RadioGroup } from '@mui/material';

// ----------------------------------------------------------------------

ColorSinglePicker.propTypes = {
  colors: PropTypes.arrayOf(PropTypes.string),
  value: PropTypes.string,
  onChange: PropTypes.func,
};

export default function ColorSinglePicker({ colors, value, onChange }) {
  return (
    <RadioGroup value={value} onChange={onChange}>
      <Box
        gap={1}
        display="grid"
        gridTemplateColumns="repeat(7, 1fr)"
      >
        {colors.map((color) => {
          return (
            <FormControlLabel
              key={color}
              value={color}
              control={<Radio />}
              label=""
              sx={{
                m: 0,
                '& .MuiFormControlLabel-label': { width: 0 },
              }}
            />
          );
        })}
      </Box>
    </RadioGroup>
  );
} 