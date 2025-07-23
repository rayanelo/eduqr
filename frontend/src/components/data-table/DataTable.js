import PropTypes from 'prop-types';
import { useState } from 'react';
// @mui
import {
  Card,
  Table,
  TableBody,
  IconButton,
  TableContainer,
  TableRow,
  TableCell,
  Chip,
  Avatar,
  Menu,
  MenuItem,
  ListItemIcon,
  ListItemText,
  Divider,
} from '@mui/material';
// components
import Iconify from '../iconify';
import Scrollbar from '../scrollbar';
import { TableEmptyRows, TableHeadCustom, TableNoData, TablePaginationCustom } from './table-components';
import { useTable, getComparator, emptyRows, applyFilter } from './use-table';
import { usePermissions, ROLES } from '../../hooks/usePermissions';

// ----------------------------------------------------------------------

DataTable.propTypes = {
  data: PropTypes.array,
  columns: PropTypes.array,
  table: PropTypes.object,
  tableData: PropTypes.array,
  dense: PropTypes.bool,
  isFiltered: PropTypes.bool,
  onDeleteRow: PropTypes.func,
  onEditRow: PropTypes.func,
  onViewRow: PropTypes.func,
  onUpdateRole: PropTypes.func,
  onAddNew: PropTypes.func,
  deleteRows: PropTypes.func,
  sx: PropTypes.object,
};

export default function DataTable({
  data,
  columns,
  table,
  tableData,
  dense,
  isFiltered,
  onDeleteRow,
  onEditRow,
  onViewRow,
  onUpdateRole,
  onAddNew,
  deleteRows,
  sx,
  ...other
}) {
  const {
    dense: denseTable,
    page,
    order,
    orderBy,
    rowsPerPage,
    selected,
    onSelectRow,
    onSort,
    onChangePage,
    onChangeRowsPerPage,
    onChangeDense,
  } = useTable();

  const denseHeight = denseTable ? 52 : 72;

  const dataFiltered = applyFilter({
    inputData: tableData,
    comparator: getComparator(order, orderBy),
  });

  const dataInPage = dataFiltered.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage);

  const isNotFound = false; // Ne jamais afficher "Aucune donnée" si on a des données



  const handleDeleteRow = (id) => {
    const deleteRow = tableData.filter((row) => row.id !== id);
    table.setData(deleteRow);

    if (selected.includes(id)) {
      const newSelected = selected.filter((selectedId) => selectedId !== id);
      table.setSelected(newSelected);
    }
  };

  const handleEditRow = (id) => {
    onEditRow(id);
  };

  const handleViewRow = (id) => {
    onViewRow(id);
  };

  const handleUpdateRole = (id, role) => {
    onUpdateRole(id, role);
  };

  return (
    <Card sx={{ ...sx }}>
      <TableContainer sx={{ position: 'relative', overflow: 'unset' }}>
        <Scrollbar>
          <Table size={dense ? 'small' : 'medium'} sx={{ minWidth: 960 }}>
            <TableHeadCustom
              order={order}
              orderBy={orderBy}
              headLabel={columns}
              onSort={onSort}
            />

            <TableBody>
              {dataInPage.map((row) => (
                <DataTableRow
                  key={row.id}
                  row={row}
                  columns={columns}
                  selected={selected.includes(row.id)}
                  onSelectRow={() => onSelectRow(row.id)}
                  onEditRow={() => handleEditRow(row.id)}
                  onViewRow={() => handleViewRow(row.id)}
                  onUpdateRole={(role) => handleUpdateRole(row.id, role)}
                  onDeleteRow={() => handleDeleteRow(row.id)}
                />
              ))}

              <TableEmptyRows
                height={denseHeight}
                emptyRows={emptyRows(page, rowsPerPage, tableData.length)}
              />

              {dataFiltered.length === 0 && <TableNoData isNotFound={isNotFound} columns={columns} onAddNew={onAddNew} />}
            </TableBody>
          </Table>
        </Scrollbar>
      </TableContainer>

      <TablePaginationCustom
        count={dataFiltered.length}
        page={page}
        rowsPerPage={rowsPerPage}
        onPageChange={onChangePage}
        onRowsPerPageChange={onChangeRowsPerPage}
        dense={denseTable}
        onChangeDense={onChangeDense}
      />
    </Card>
  );
}

// ----------------------------------------------------------------------

function DataTableRow({
  row,
  columns,
  selected,
  onSelectRow,
  onEditRow,
  onViewRow,
  onUpdateRole,
  onDeleteRow,
}) {
  const [anchorEl, setAnchorEl] = useState(null);
  const { canUpdateUser, canDeleteUser, canChangeUserRole, getPromotableRoles } = usePermissions();

  const handleOpenMenu = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleCloseMenu = () => {
    setAnchorEl(null);
  };

  const renderCell = (column, value) => {
    switch (column.type) {
      case 'avatar':
        return (
          <Avatar alt={value} src={value} sx={{ width: 40, height: 40 }}>
            {value?.charAt(0)?.toUpperCase()}
          </Avatar>
        );
      case 'status':
        return (
          <Chip
            label={value}
            color={value === 'active' ? 'success' : 'warning'}
            size="small"
          />
        );
      case 'role':
        const getRoleColor = (role) => {
          switch (role) {
            case 'super_admin':
              return 'error';
            case 'admin':
              return 'warning';
            case 'professeur':
              return 'info';
            case 'etudiant':
              return 'default';
            default:
              return 'default';
          }
        };

        const getRoleLabel = (role) => {
          switch (role) {
            case 'super_admin':
              return 'Super Admin';
            case 'admin':
              return 'Admin';
            case 'professeur':
              return 'Professeur';
            case 'etudiant':
              return 'Étudiant';
            default:
              return role;
          }
        };

        return (
          <Chip
            label={getRoleLabel(value)}
            color={getRoleColor(value)}
            size="small"
          />
        );
      case 'date':
        return new Date(value).toLocaleDateString();
      default:
        return value;
    }
  };

  return (
    <TableRow hover>
      {columns.map((column) => (
        <TableCell key={column.id} align={column.align || 'left'}>
          {column.id === 'actions' ? (
            <IconButton onClick={handleOpenMenu}>
              <Iconify icon="eva:more-vertical-fill" />
            </IconButton>
          ) : (
            renderCell(column, row[column.id])
          )}
        </TableCell>
      ))}

      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleCloseMenu}
        anchorOrigin={{ vertical: 'top', horizontal: 'left' }}
        transformOrigin={{ vertical: 'top', horizontal: 'right' }}
      >
        <MenuItem onClick={() => { onViewRow(); handleCloseMenu(); }}>
          <ListItemIcon>
            <Iconify icon="eva:eye-fill" />
          </ListItemIcon>
          <ListItemText>Voir</ListItemText>
        </MenuItem>

        {canUpdateUser(row.role) && (
          <MenuItem onClick={() => { onEditRow(); handleCloseMenu(); }}>
            <ListItemIcon>
              <Iconify icon="eva:edit-fill" />
            </ListItemIcon>
            <ListItemText>Modifier</ListItemText>
          </MenuItem>
        )}

        {canChangeUserRole(row.role) && (
          <>
            <Divider sx={{ borderStyle: 'dashed' }} />
            {getPromotableRoles(row.role).map((role) => (
              <MenuItem 
                key={role} 
                onClick={() => { onUpdateRole(role); handleCloseMenu(); }}
              >
                <ListItemIcon>
                  <Iconify icon={role === ROLES.ETUDIANT ? "eva:person-fill" : "eva:shield-fill"} />
                </ListItemIcon>
                <ListItemText>
                  {role === ROLES.SUPER_ADMIN && "Promouvoir Super Admin"}
                  {role === ROLES.ADMIN && "Promouvoir Admin"}
                  {role === ROLES.PROFESSEUR && "Promouvoir Professeur"}
                  {role === ROLES.ETUDIANT && "Rétrograder Étudiant"}
                </ListItemText>
              </MenuItem>
            ))}
          </>
        )}

        {canDeleteUser(row.role) && (
          <>
            <Divider sx={{ borderStyle: 'dashed' }} />
            <MenuItem onClick={() => { onDeleteRow(); handleCloseMenu(); }} sx={{ color: 'error.main' }}>
              <ListItemIcon>
                <Iconify icon="eva:trash-2-fill" sx={{ color: 'error.main' }} />
              </ListItemIcon>
              <ListItemText>Supprimer</ListItemText>
            </MenuItem>
          </>
        )}
      </Menu>
    </TableRow>
  );
}

DataTableRow.propTypes = {
  row: PropTypes.object,
  columns: PropTypes.array,
  selected: PropTypes.bool,
  onSelectRow: PropTypes.func,
  onEditRow: PropTypes.func,
  onViewRow: PropTypes.func,
  onUpdateRole: PropTypes.func,
  onDeleteRow: PropTypes.func,
}; 