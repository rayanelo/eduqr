// @mui
import { RadioGroup } from '@mui/material';
//
import Iconify from '../../iconify';
import { useSettingsContext } from '../SettingsContext';
import { StyledCard, StyledWrap, MaskControl } from '../styles';

// ----------------------------------------------------------------------

const OPTIONS = ['default', 'bold'];

export default function ContrastOptions() {
  const { themeContrast, onChangeContrast } = useSettingsContext();

  return (
    <RadioGroup name="themeContrast" value={themeContrast} onChange={onChangeContrast}>
      <StyledWrap>
        {OPTIONS.map((contrast) => (
          <StyledCard key={contrast} selected={themeContrast === contrast}>
            <Iconify
              icon={
                contrast === 'bold' ? 'solar:contrast-bold-linear' : 'solar:contrast-linear'
              }
            />

            <MaskControl value={contrast} />
          </StyledCard>
        ))}
      </StyledWrap>
    </RadioGroup>
  );
}
