import { useState } from 'react';
//
import Iconify from '../../iconify';
import { StyledCard } from '../styles';

// ----------------------------------------------------------------------

export default function FullScreenOptions() {
  const [fullscreen, setFullscreen] = useState(false);

  const onToggleFullScreen = () => {
    if (!document.fullscreenElement) {
      document.documentElement.requestFullscreen();
      setFullscreen(true);
    } else if (document.exitFullscreen) {
      document.exitFullscreen();
      setFullscreen(false);
    }
  };

  return (
    <StyledCard
      selected={fullscreen}
      onClick={onToggleFullScreen}
      sx={{
        height: 48,
        typography: 'subtitle2',
        '& .svg-color': {
          ml: 1,
          width: 16,
          height: 16,
        },
      }}
    >
      {fullscreen ? 'Exit Fullscreen' : 'Fullscreen'}

      <Iconify
        icon={fullscreen ? 'solar:exit-full-screen-linear' : 'solar:full-screen-linear'}
      />
    </StyledCard>
  );
}
