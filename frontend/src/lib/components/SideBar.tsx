import React from 'react';
import { useAtomValue } from 'jotai';
import { Box } from '@mantine/core';

import { currentAssetAtom } from 'stores/AssetPage';
/**
 *
 * @returns
 */
function SideBar() {
  const pageState = useAtomValue(currentAssetAtom);
  return (
    <Box>
      {pageState?.Label}
    </Box>
  );
}

export default SideBar;
