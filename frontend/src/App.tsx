import '@mantine/core/styles.css';

import classes from './App.module.css';
import { navigationItems } from './navigation';
import { Outlet } from 'react-router';
import { SettingsProvider } from './core/providers/settings/settingsProvider';
import { Box, LoadingOverlay, useComputedColorScheme, useMantineTheme } from '@mantine/core';
import { NavBar } from './core/components/navbar/NavBar';

export default function App() {
  const computedColorScheme = useComputedColorScheme('light', { getInitialValueInEffect: true });
  const theme = useMantineTheme();

  return <Box className={classes.wrapper}
    style={{
      backgroundColor: computedColorScheme === 'dark' ? theme.colors.dark[9] : theme.colors.gray[3],
    }}>
    <SettingsProvider
      loading={<LoadingOverlay visible={true} zIndex={1000} overlayProps={{ blur: 2 }} />}
    >
      <NavBar navItems={navigationItems} />
      <Outlet />

    </SettingsProvider>
  </Box>
}