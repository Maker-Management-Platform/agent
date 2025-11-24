import React from 'react';
import { Outlet } from 'react-router';
import { Box } from '@mantine/core';

import { Viewer3dPanel } from 'lib/components/viewer3d/Viewer3dPanel';

import classes from './LibMain.module.css';

/**
 *
 * @returns
 */
function LibMain() {
  return (
    <Box className={classes.main}>
      <Outlet />
      <Viewer3dPanel />
    </Box>
  );
}

export default LibMain;
