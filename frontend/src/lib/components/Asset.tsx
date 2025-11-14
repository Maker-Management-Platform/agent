import React, { useState } from 'react';
import { Link, useLocation, useParams } from 'react-router';
import { useAtom, useSetAtom } from 'jotai';
import {
  Anchor,
  Box,
  Breadcrumbs,
  Group,
  Collapse,
  ActionIcon,
  Paper,
} from '@mantine/core';
import {
  IconCircleDashedPlus,
  IconEdit,
  IconInfoCircle,
} from '@tabler/icons-react';

import { useGetAsset } from 'fetchers/getAsset';
import { currentAssetAtom, viewer3dAssetsAtom } from 'stores/AssetPage';
import { type Asset as AssetType } from 'types/Asset';

import AssetCard from './AssetCard';
import classes from './Asset.module.css';
import Header from './Header';
import AssetEditForm from './AssetEditForm';
import AssetNewForm from './AssetNewForm';

/**
 *
 * @returns
 */
function Asset() {
  const loc = useLocation();
  const { id } = useParams();

  const setAsset = useSetAtom(currentAssetAtom);
  const [viewer3dAssets, setViewer3dAssets] = useAtom(viewer3dAssetsAtom);
  const [expandedView, setExpandedView] = useState(loc.hash.slice(1));

  const { data, isLoading } = useGetAsset(id);

  /**
   *
   * @param asset
   * @returns
   */
  function calculateBreadcrumbs(asset: AssetType | undefined): JSX.Element[] {
    if (asset && asset.Parent) {
      return [
        ...calculateBreadcrumbs(asset.Parent),
        <Anchor
          component={Link}
          to={`/lib/${asset.Parent.ID}`}
          key={asset.Parent.ID}
        >
          {asset.Parent.Label}
        </Anchor>,
      ];
    }

    return [
      <Anchor key="home" component={Link} to="/lib">
        Home
      </Anchor>,
    ];
  }

  const breadcrumbs = calculateBreadcrumbs(data);

  const toggleExpandedView = (view: string) => (
    setExpandedView(() => {
      globalThis.location.hash = view;
      if (expandedView === view) {
        globalThis.location.hash = '';
        return '';
      }
      return view;
    }));

  return (
    <Paper className={classes.asset} m="sm" display="flex">
      <Header
        title={data?.Label ?? ''}
        description={data?.Description ?? ''}
        loading={isLoading}
        imgID={data?.Thumbnail ?? ''}
      />
      <Group
        justify="space-between"
        preventGrowOverflow={false}
        mt="xs"
        mx="sm"
      >
        <Breadcrumbs separator="→" separatorMargin="sm">
          {breadcrumbs}
        </Breadcrumbs>
        <Group justify="flex-end" wrap="nowrap">
          <ActionIcon
            variant={expandedView === 'details' ? 'light' : 'filled'}
            aria-label="Details"
            onClick={() => toggleExpandedView('details')}
          >
            <IconInfoCircle
              style={{ width: '70%', height: '70%' }}
              stroke={1.5}
            />
          </ActionIcon>
          <ActionIcon
            variant={expandedView === 'edit' ? 'light' : 'filled'}
            aria-label="Settings"
            onClick={() => toggleExpandedView('edit')}
          >
            <IconEdit style={{ width: '70%', height: '70%' }} stroke={1.5} />
          </ActionIcon>
          <ActionIcon
            variant={expandedView === 'add' ? 'light' : 'filled'}
            aria-label="Settings"
            onClick={() => toggleExpandedView('add')}
          >
            <IconCircleDashedPlus
              style={{ width: '70%', height: '70%' }}
              stroke={1.5}
            />
          </ActionIcon>
        </Group>
      </Group>

      <Box className={classes.content} mt="sm">
        <Collapse in={expandedView === 'details'}>
          <Box>a thing</Box>
        </Collapse>
        <Collapse in={expandedView === 'edit'}>
          {data && (
            <AssetEditForm asset={data} onClose={() => setExpandedView('')} />
          )}
        </Collapse>
        <Collapse in={expandedView === 'add'}>
          {data && (
            <AssetNewForm parent={data} onClose={() => setExpandedView('')} />
          )}
        </Collapse>
      </Box>

      <Box className={classes.list} m="sm">
        {data && data?.NestedAssets.map((a: AssetType) => (
          <AssetCard
            asset={a}
            key={a.ID}
            onSelected={(asset: AssetType) => {
              setAsset(asset);
            }}
            onAddTo3dViewer={(asset: AssetType) => {
              setAsset(asset);

              setViewer3dAssets([...viewer3dAssets, asset]);
            }}
          />
        ))}
      </Box>
    </Paper>
  );
}

export default Asset;
