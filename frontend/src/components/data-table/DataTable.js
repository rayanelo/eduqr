import React, { useState, useEffect } from 'react';
import { useTheme } from '@mui/material/styles';
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  TablePagination,
  TableSortLabel,
} from '@mui/material';

export function DataTable({ data, columns, onAddNew, isFiltered = false }) {
  const theme = useTheme();
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [orderBy, setOrderBy] = useState('');
  const [order, setOrder] = useState('asc');
  const [themeKey, setThemeKey] = useState(0);

  // Force re-render when theme changes
  useEffect(() => {
    setThemeKey(prev => prev + 1);
  }, [theme.palette.primary.main]);

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleRequestSort = (property) => {
    const isAsc = orderBy === property && order === 'asc';
    setOrder(isAsc ? 'desc' : 'asc');
    setOrderBy(property);
  };

  const sortedData = React.useMemo(() => {
    if (!orderBy) return data;

    return [...data].sort((a, b) => {
      const aValue = a[orderBy];
      const bValue = b[orderBy];

      if (aValue < bValue) {
        return order === 'asc' ? -1 : 1;
      }
      if (aValue > bValue) {
        return order === 'asc' ? 1 : -1;
      }
      return 0;
    });
  }, [data, orderBy, order]);

  const paginatedData = sortedData.slice(
    page * rowsPerPage,
    page * rowsPerPage + rowsPerPage
  );

  const createSortHandler = (property) => () => {
    handleRequestSort(property);
  };

  return (
    <Paper 
      key={themeKey}
      sx={{ 
        width: '100%', 
        overflow: 'hidden', 
        boxShadow: theme.customShadows?.primary || theme.shadows[2],
        borderRadius: 2,
        border: `1px solid ${theme.palette.divider}`
      }}
    >
      <TableContainer sx={{ maxHeight: 440 }}>
        <Table stickyHeader>
          <TableHead>
            <TableRow>
              {columns.map((column) => (
                <TableCell
                  key={column.id}
                  align={column.align || 'left'}
                  style={{ 
                    minWidth: column.minWidth, 
                    width: column.width,
                  }}
                  sortDirection={orderBy === column.id ? order : false}
                  sx={{
                    backgroundColor: `${theme.palette.primary.main} !important`,
                    color: `${theme.palette.primary.contrastText} !important`,
                    fontWeight: 'bold',
                    fontSize: '0.875rem',
                    '&:hover': {
                      backgroundColor: `${theme.palette.primary.dark} !important`,
                    },
                  }}
                >
                  {column.sortable !== false ? (
                    <TableSortLabel
                      active={orderBy === column.id}
                      direction={orderBy === column.id ? order : 'asc'}
                      onClick={createSortHandler(column.id)}
                      sx={{
                        color: `${theme.palette.primary.contrastText} !important`,
                        '&:hover': {
                          color: `${theme.palette.primary.contrastText} !important`,
                          opacity: 0.8,
                        },
                        '&.Mui-active': {
                          color: `${theme.palette.primary.contrastText} !important`,
                          '& .MuiTableSortLabel-icon': {
                            color: `${theme.palette.primary.contrastText} !important`,
                          },
                        },
                        '& .MuiTableSortLabel-icon': {
                          color: `${theme.palette.primary.contrastText} !important`,
                        },
                      }}
                    >
                      {column.label}
                    </TableSortLabel>
                  ) : (
                    <span style={{ 
                      color: `${theme.palette.primary.contrastText} !important`,
                      fontWeight: 'bold',
                    }}>
                      {column.label}
                    </span>
                  )}
                </TableCell>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {paginatedData.map((row, index) => (
              <TableRow 
                hover 
                role="checkbox" 
                tabIndex={-1} 
                key={row.id || index}
                sx={{
                  '&:nth-of-type(odd)': {
                    backgroundColor: theme.palette.mode === 'light' 
                      ? theme.palette.grey[50] 
                      : theme.palette.grey[800],
                  },
                  '&:hover': {
                    backgroundColor: theme.palette.mode === 'light' 
                      ? theme.palette.grey[100] 
                      : theme.palette.grey[700],
                    '& .MuiTableCell-root': {
                      color: theme.palette.text.primary,
                    },
                  },
                  transition: 'all 0.2s ease-in-out',
                }}
              >
                {columns.map((column) => {
                  const value = row[column.id];
                  return (
                    <TableCell 
                      key={column.id} 
                      align={column.align || 'left'}
                      sx={{
                        borderBottom: `1px solid ${theme.palette.divider}`,
                        py: 1.5,
                        fontSize: '0.875rem',
                      }}
                    >
                      {column.render ? column.render(value, row) : value}
                    </TableCell>
                  );
                })}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <TablePagination
        rowsPerPageOptions={[10, 25, 100]}
        component="div"
        count={data.length}
        rowsPerPage={rowsPerPage}
        page={page}
        onPageChange={handleChangePage}
        onRowsPerPageChange={handleChangeRowsPerPage}
        labelRowsPerPage="Lignes par page:"
        labelDisplayedRows={({ from, to, count }) =>
          `${from}-${to} sur ${count !== -1 ? count : `plus de ${to}`}`
        }
        sx={{
          backgroundColor: theme.palette.mode === 'light' 
            ? theme.palette.grey[100] 
            : theme.palette.grey[800],
          borderTop: `1px solid ${theme.palette.divider}`,
          '& .MuiTablePagination-selectLabel, & .MuiTablePagination-displayedRows': {
            color: theme.palette.text.primary,
          },
        }}
      />
    </Paper>
  );
} 