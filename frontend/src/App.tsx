import React from 'react';
import { Outlet } from 'react-router';
import { useSetAtom } from 'jotai';
import {
  Box,
  LoadingOverlay,
  useComputedColorScheme,
  useMantineTheme,
} from '@mantine/core';

import classes from 'App.module.css';
import { navigationItems } from 'navigation';
import NavBar from 'core/components/navbar/NavBar';
import { settingsAtom } from 'core/stores/settings.store';

import '@mantine/core/styles.css';

/**
 *
 * @returns
 */
function App() {
  const computedColorScheme = useComputedColorScheme('light', { getInitialValueInEffect: true });
  const [ready, setReady] = React.useState(false);
  const theme = useMantineTheme();

  const setSettings = useSetAtom(settingsAtom);

  React.useEffect(() => {
    if (ready) {
      return;
    }

    fetch('/settings.json')
      .then((response) => response.json())
      .then((data) => setSettings((prev) => ({ ...prev, ...data })))
      .then(() => setReady(true))
      .catch(console.error);
  }, [ready, setSettings]);

  if (!ready) {
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
