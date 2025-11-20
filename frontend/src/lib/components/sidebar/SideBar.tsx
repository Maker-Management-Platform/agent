import React from 'react';
import { Box } from '@mantine/core';

import { useAssetPage } from '../../contexts/AssetPageContext';

import classes from './SideBar.module.css';

/**
 *
 * @returns
 */
function SideBar() {
  const pageState = useAssetPage();
  return (
    <Box className={classes.sidebar}>
      {pageState.asset?.Label}
    </Box>
  );
}

export default SideBar;
