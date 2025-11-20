import React from 'react';
import { Box } from '@mantine/core';

import { type Asset } from '../../models';

/**
 *
 * @param param0
 * @param param0.asset
 * @returns
 */
function AssetDetails({ asset }: { readonly asset: Asset }) {
  return (
    <Box>
      <Box>{asset.Label}</Box>
    </Box>
  );
}

export default AssetDetails;
