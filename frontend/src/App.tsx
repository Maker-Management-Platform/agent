import React from 'react';
import { Outlet } from 'react-router';
import {
  Box,
  LoadingOverlay,
  useComputedColorScheme,
  useMantineTheme,
} from '@mantine/core';

import classes from 'App.module.css';
import { navigationItems } from 'navigation';
import SettingsProvider from 'core/providers/settings/settingsProvider';
import NavBar from 'core/components/navbar/NavBar';

import '@mantine/core/styles.css';

/**
 *
 * @returns
 */
function App() {
  const computedColorScheme = useComputedColorScheme('light', { getInitialValueInEffect: true });
  const theme = useMantineTheme();

  return (
    <Box
      className={classes.wrapper}
      style={{
        backgroundColor: computedColorScheme === 'dark' ? theme.colors.dark[9] : theme.colors.gray[3],
      }}
    >
      <SettingsProvider
        loading={<LoadingOverlay visible zIndex={1000} overlayProps={{ blur: 2 }} />}
      >
        <NavBar navItems={navigationItems} />
        <Outlet />

      </SettingsProvider>
    </Box>
  );
}

export default App;
