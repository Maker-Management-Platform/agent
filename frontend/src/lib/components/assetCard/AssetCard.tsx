import React from 'react';
import {
  ActionIcon,
  Box,
  rem,
  Title,
} from '@mantine/core';
import { Link } from 'react-router';
import { useHover } from '@mantine/hooks';
import { IconHeart, IconSettings } from '@tabler/icons-react';

import { getLocalBackend } from 'core/stores/settings.store';

import { type Asset } from '../../models';

import classes from './AssetCard.module.css';

/**
 *
 * @param param0
 * @param param0.asset
 * @param param0.onSelected
 * @param param0.onAddTo3dViewer
 * @returns
 */
function AssetCard({
  asset,
  onSelected,
  onAddTo3dViewer,
}: {
  readonly asset: Asset;
  readonly onSelected?: (asset: Asset) => void;
  readonly onAddTo3dViewer?: (asset: Asset) => void;
}) {
  const { hovered, ref } = useHover();
  const shouldNavigate = (a: Asset) => (
    (
      a.NodeKind === 'bundle' || a.NodeKind === 'dir' || a.NodeKind === 'root'
    )
  );

  return (
    <Box className={classes.card} ref={ref} m="0.5rem">
      <div className={classes.content}>
        {shouldNavigate(asset) ? (
          <Title order={5} component={Link} to={`/lib/${asset.ID}`}>
            {asset.Label}
          </Title>
        ) : (
          <Title order={5}>{asset.Label}</Title>
        )}
      </div>
      <div
        className={classes.controls}
        style={{ display: hovered ? 'flex' : 'none' }}
      >
        <ActionIcon.Group>
          {asset.Kind === 'model' && (
            <ActionIcon
              variant="white"
              size="sm"
              aria-label="Settings"
              onClick={
                onAddTo3dViewer ? () => onAddTo3dViewer(asset) : undefined
              }
            >
              <IconSettings style={{ width: rem(20) }} stroke={1.5} />
            </ActionIcon>
          )}
          <ActionIcon
            variant="white"
            size="sm"
            aria-label="Likes"
            onClick={onSelected ? () => onSelected(asset) : undefined}
          >
            <IconHeart style={{ width: rem(20) }} stroke={1.5} />
          </ActionIcon>
        </ActionIcon.Group>
      </div>
      <div
        style={{
          backgroundImage: `url(${getLocalBackend()}/api/lib/${asset.Thumbnail}/file)`,
        }}
        className={classes.thumbnail}
      >
        <div className={classes.overlay} />
      </div>
    </Box>
  );
}

export default AssetCard;
