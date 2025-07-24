import { useState } from 'react';
import { Helmet } from 'react-helmet-async';
// @mui
import {
  Card,
  Container,
  Tab,
  Tabs,
  Box,
  Stack,
} from '@mui/material';
// hooks
import { useProfile } from '../../hooks/useProfile';
// routes
import { PATH_DASHBOARD } from '../../routes/paths';
// components
import { useSnackbar } from '../../components/snackbar';
import CustomBreadcrumbs from '../../components/custom-breadcrumbs';
import { useSettingsContext } from '../../components/settings';
import UserInfo from '../../components/user-info/UserInfo';
import ProfileInfoTab from '../../sections/profile/ProfileInfoTab';
import PasswordTab from '../../sections/profile/PasswordTab';

// ----------------------------------------------------------------------

function TabPanel({ children, value, index, ...other }) {
  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`profile-tabpanel-${index}`}
      aria-labelledby={`profile-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ py: 3 }}>{children}</Box>}
    </div>
  );
}

export default function ProfilePage() {
  const { themeStretch } = useSettingsContext();
  const { enqueueSnackbar } = useSnackbar();
  const { user, isLoading, updateProfile, changePassword, validatePassword } = useProfile();

  const [currentTab, setCurrentTab] = useState(0);

  const handleTabChange = (event, newValue) => {
    setCurrentTab(newValue);
  };

  const handleUpdateProfile = async (data) => {
    try {
      await updateProfile(data);
      enqueueSnackbar('Profil mis à jour avec succès!');
    } catch (error) {
      enqueueSnackbar(error.message, { variant: 'error' });
    }
  };

  const handleChangePassword = async (data) => {
    try {
      await changePassword(data);
      enqueueSnackbar('Mot de passe modifié avec succès!');
    } catch (error) {
      enqueueSnackbar(error.message, { variant: 'error' });
    }
  };

  if (!user) {
    return null;
  }

  return (
    <>
      <Helmet>
        <title> Mon profil | EduQR</title>
      </Helmet>

      <Container maxWidth={themeStretch ? false : 'lg'}>
        <CustomBreadcrumbs
          heading="Mon profil"
          links={[
            {
              name: 'Dashboard',
              href: PATH_DASHBOARD.root,
            },
            {
              name: 'Profil',
            },
          ]}
        />

        <Stack spacing={3}>
          {/* UserInfo en haut */}
          <UserInfo />

          {/* Onglets */}
          <Card>
            <Tabs
              value={currentTab}
              onChange={handleTabChange}
              sx={{
                px: 2,
                bgcolor: 'background.neutral',
              }}
            >
              <Tab label="Informations personnelles" />
              <Tab label="Sécurité" />
            </Tabs>

            <TabPanel value={currentTab} index={0}>
              <ProfileInfoTab
                user={user}
                onUpdate={handleUpdateProfile}
                isLoading={isLoading}
              />
            </TabPanel>

            <TabPanel value={currentTab} index={1}>
              <PasswordTab
                onChangePassword={handleChangePassword}
                onValidatePassword={validatePassword}
                isLoading={isLoading}
              />
            </TabPanel>
          </Card>
        </Stack>
      </Container>
    </>
  );
} 