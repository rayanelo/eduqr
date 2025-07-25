import PropTypes from 'prop-types';
// @mui
import { Paper, Typography } from '@mui/material';

// ----------------------------------------------------------------------

SearchNotFound.propTypes = {
  query: PropTypes.string,
  sx: PropTypes.object,
};

export default function SearchNotFound({ query, sx, ...other }) {
  return query ? (
    <Paper
      sx={{
        textAlign: 'center',
        ...sx,
      }}
      {...other}
    >
      <Typography variant="h6" paragraph>
        Aucun résultat trouvé
      </Typography>

      <Typography variant="body2">
        Aucun résultat trouvé pour &nbsp;
        <strong>&quot;{query}&quot;</strong>.
        <br /> Vérifiez l'orthographe ou utilisez des mots complets.
      </Typography>
    </Paper>
  ) : (
          <Typography variant="body2" sx={sx}>
        Veuillez saisir des mots-clés
      </Typography>
  );
}
