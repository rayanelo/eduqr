// @mui
import { RadioGroup } from '@mui/material';
//
import Iconify from '../../iconify';
import { useSettingsContext } from '../SettingsContext';
import { StyledCard, StyledWrap, MaskControl } from '../styles';

// ----------------------------------------------------------------------

const OPTIONS = ['ltr', 'rtl'];

export default function DirectionOptions() {
  const { themeDirection, onChangeDirection } = useSettingsContext();

  return (
    <RadioGroup name="themeDirection" value={themeDirection} onChange={onChangeDirection}>
      <StyledWrap>
        {OPTIONS.map((direction) => (
          <StyledCard key={direction} selected={themeDirection === direction}>
            <Iconify
              icon={
                direction === 'rtl' ? 'solar:align-right-linear' : 'solar:align-left-linear'
              }
            />

            <MaskControl value={direction} />
          </StyledCard>
        ))}
      </StyledWrap>
    </RadioGroup>
  );
}
