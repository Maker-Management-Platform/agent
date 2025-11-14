import React from 'react';
import {
  Box,
  LoadingOverlay,
  useComputedColorScheme,
  useMantineTheme,
} from '@mantine/core';
import '@mantine/core/styles.css';
import { Outlet } from 'react-router';

import { type INavItem } from 'types/Nav';
import NavBar from 'core/components/navbar/NavBar';
import dashNavItem from 'dashboard/routes';
import libNavItem from 'lib/routes';
import { useGetSettings } from 'fetchers/getSettings';

import classes from './App.module.css';

const navigationItems: INavItem[] = [
  dashNavItem,
  libNavItem,
];

/**
 *
 * @returns
 */
function App() {
  const { isLoading } = useGetSettings();
  const computedColorScheme = useComputedColorScheme('light', { getInitialValueInEffect: true });
  const theme = useMantineTheme();

  if (isLoading) {
    return <LoadingOverlay visible zIndex={1000} overlayProps={{ blur: 2 }} />;
  }

  return (
    <Box
      className={classes.wrapper}
      style={{
        backgroundColor: computedColorScheme === 'dark' ? theme.colors.dark[9] : theme.colors.gray[3],
      }}
    >
      <NavBar navItems={navigationItems} />
      <Outlet />
    </Box>
  );
}

export default App;
