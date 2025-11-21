import React from 'react';
import { Box } from '@mantine/core';
import { useAtomValue } from 'jotai';

import { currentAssetAtom } from 'lib/stores/assetPage.store';

import classes from './SideBar.module.css';

/**
 *
 * @returns
 */
function SideBar() {
  const asset = useAtomValue(currentAssetAtom);
  return (
    <Box className={classes.sidebar}>
      {asset?.Label}
    </Box>
  );
}

export default SideBar;
