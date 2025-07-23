import PropTypes from 'prop-types';
// @mui
import { Box, Checkbox, FormControlLabel } from '@mui/material';

// ----------------------------------------------------------------------

ColorMultiPicker.propTypes = {
  colors: PropTypes.arrayOf(PropTypes.string),
  selected: PropTypes.arrayOf(PropTypes.string),
  onChangeColor: PropTypes.func,
  sx: PropTypes.object,
};

export default function ColorMultiPicker({ colors, selected, onChangeColor, sx }) {
  return (
    <Box
      gap={1}
      display="grid"
      gridTemplateColumns="repeat(7, 1fr)"
      sx={sx}
    >
      {colors.map((color) => {
        const isSelected = selected.includes(color);

        return (
          <FormControlLabel
            key={color}
            control={
              <Checkbox
                checked={isSelected}
                onChange={() => onChangeColor(color)}
                sx={{
                  color: color,
                  '&.Mui-checked': {
                    color: color,
                  },
                }}
              />
            }
            label=""
            sx={{
              m: 0,
              '& .MuiFormControlLabel-label': { width: 0 },
            }}
          />
        );
      })}
    </Box>
  );
} 