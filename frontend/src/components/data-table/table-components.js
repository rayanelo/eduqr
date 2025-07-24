import PropTypes from 'prop-types';
// @mui
import {
  Box,
  Button,
  TableBody,
  TableHead,
  TableRow,
  TableCell,
  TableSortLabel,
  Typography,
  Checkbox,
  TablePagination,
  Switch,
  FormControlLabel,
} from '@mui/material';
// components
import Iconify from '../iconify';

// ----------------------------------------------------------------------

export function TableHeadCustom({
  order,
  orderBy,
  headLabel,
  rowCount = 0,
  numSelected = 0,
  onSort,
  onSelectAllRows,
  sx,
}) {
  return (
    <TableHead sx={sx}>
      <TableRow>
        {onSelectAllRows && (
          <TableCell padding="checkbox">
            <Checkbox
              indeterminate={numSelected > 0 && numSelected < rowCount}
              checked={rowCount > 0 && numSelected === rowCount}
              onChange={(event) => onSelectAllRows(event.target.checked)}
            />
          </TableCell>
        )}

        {headLabel.map((headCell) => (
          <TableCell
            key={headCell.id}
            align={headCell.align || 'left'}
            sortDirection={orderBy === headCell.id ? order : false}
            sx={{ width: headCell.width, minWidth: headCell.minWidth }}
          >
            {onSort ? (
              <TableSortLabel
                hideSortIcon
                active={orderBy === headCell.id}
                direction={orderBy === headCell.id ? order : 'asc'}
                onClick={() => onSort(headCell.id)}
              >
                {headCell.label}

                {orderBy === headCell.id ? (
                  <Box sx={{ typography: 'caption', ml: 1 }}>
                    {order === 'desc' ? '↓' : '↑'}
                  </Box>
                ) : null}
              </TableSortLabel>
            ) : (
              headCell.label
            )}
          </TableCell>
        ))}
      </TableRow>
    </TableHead>
  );
}

TableHeadCustom.propTypes = {
  sx: PropTypes.object,
  onSort: PropTypes.func,
  orderBy: PropTypes.string,
  headLabel: PropTypes.array,
  rowCount: PropTypes.number,
  numSelected: PropTypes.number,
  onSelectAllRows: PropTypes.func,
  order: PropTypes.oneOf(['asc', 'desc']),
};

// ----------------------------------------------------------------------

export function TableSelectedAction({
  dense,
  action,
  rowCount,
  numSelected,
  onSelectAllRows,
  sx,
  ...other
}) {
  return (
    <Box
      sx={{
        px: 2,
        top: 0,
        left: 0,
        right: 0,
        zIndex: 9,
        height: 58,
        display: 'flex',
        alignItems: 'center',
        position: 'absolute',
        bgcolor: 'primary.lighter',
        ...(dense && {
          pl: 1,
          height: 38,
        }),
        ...sx,
      }}
      {...other}
    >
      <Checkbox
        indeterminate={numSelected > 0 && numSelected < rowCount}
        checked={rowCount > 0 && numSelected === rowCount}
        onChange={(event) => onSelectAllRows(event.target.checked)}
      />

      <Typography
        variant="subtitle1"
        sx={{
          ml: 2,
          flexGrow: 1,
          color: 'primary.main',
          ...(dense && {
            ml: 3,
            variant: 'subtitle2',
          }),
        }}
      >
        {numSelected} selected
      </Typography>

      {action && action}
    </Box>
  );
}

TableSelectedAction.propTypes = {
  action: PropTypes.node,
  rowCount: PropTypes.number,
  numSelected: PropTypes.number,
  onSelectAllRows: PropTypes.func,
  dense: PropTypes.bool,
  sx: PropTypes.object,
};

// ----------------------------------------------------------------------

export function TableEmptyRows({ emptyRows, height }) {
  if (!emptyRows) {
    return null;
  }

  return (
    <TableRow
      sx={{
        height: height || 52,
      }}
    >
      <TableCell colSpan={9} />
    </TableRow>
  );
}

TableEmptyRows.propTypes = {
  emptyRows: PropTypes.number,
  height: PropTypes.number,
};

// ----------------------------------------------------------------------

export function TableNoData({ isNotFound, columns = [], onAddNew, ...other }) {
  return (
    <TableBody>
      <TableRow>
        <TableCell 
          colSpan={columns.length || 8} 
          align="center" 
          sx={{ 
            border: 'none',
            py: 8,
          }}
          {...other}
        >
          <Box
            sx={{
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              justifyContent: 'center',
              py: 4,
              px: 2,
            }}
          >
            <Box
              component="img"
              src="/assets/illustrations/illustration_empty_content.svg"
              sx={{
                width: 120,
                height: 120,
                mb: 3,
                opacity: 0.6,
              }}
            />
            
            <Typography 
              variant="h6" 
              sx={{ 
                mb: 1,
                color: 'text.secondary',
                fontWeight: 500,
              }}
            >
              {isNotFound ? 'Aucune donnée trouvée' : 'Aucune donnée'}
            </Typography>

            <Typography 
              variant="body2" 
              sx={{ 
                color: 'text.disabled',
                textAlign: 'center',
                maxWidth: 300,
                mb: 3,
              }}
            >
              {isNotFound
                ? 'Aucun résultat trouvé pour votre recherche. Essayez de modifier vos critères.'
                : 'Il n\'y a pas encore de données dans cette liste. Commencez par ajouter un nouvel élément.'}
            </Typography>

            {!isNotFound && onAddNew && (
              <Button
                variant="contained"
                startIcon={<Iconify icon="eva:plus-fill" />}
                onClick={onAddNew}
                sx={{
                  mt: 1,
                }}
              >
                Ajouter le premier élément
              </Button>
            )}
          </Box>
        </TableCell>
      </TableRow>
    </TableBody>
  );
}

TableNoData.propTypes = {
  isNotFound: PropTypes.bool,
};

// ----------------------------------------------------------------------

export function TablePaginationCustom({
  dense,
  onChangeDense,
  rowsPerPageOptions = [5, 10, 25],
  sx,
  ...other
}) {
  return (
    <Box sx={{ position: 'relative', ...sx }}>
      <TablePagination
        rowsPerPageOptions={rowsPerPageOptions}
        component="div"
        {...other}
        sx={{
          borderTop: (theme) => `solid 1px ${theme.palette.divider}`,
        }}
      />

      {onChangeDense && (
        <FormControlLabel
          label="Dense"
          control={<Switch checked={dense} onChange={onChangeDense} />}
          sx={{
            pl: 2,
            py: 1.5,
            top: 0,
            position: {
              sm: 'absolute',
            },
          }}
        />
      )}
    </Box>
  );
}

TablePaginationCustom.propTypes = {
  dense: PropTypes.bool,
  onChangeDense: PropTypes.func,
  rowsPerPageOptions: PropTypes.arrayOf(PropTypes.number),
  sx: PropTypes.object,
};

// ----------------------------------------------------------------------

export function TableToolbar({ filters, onFilters, onResetFilters, ...other }) {
  return (
    <Box
      sx={{
        p: 2.5,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        borderBottom: (theme) => `solid 1px ${theme.palette.divider}`,
      }}
      {...other}
    >
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
        {filters}
      </Box>

      {onResetFilters && (
        <Button
          color="error"
          onClick={onResetFilters}
          startIcon={<Iconify icon="solar:trash-bin-trash-bold" />}
        >
          Clear
        </Button>
      )}
    </Box>
  );
}

TableToolbar.propTypes = {
  filters: PropTypes.node,
  onFilters: PropTypes.func,
  onResetFilters: PropTypes.func,
}; 