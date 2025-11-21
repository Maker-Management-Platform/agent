import React, {
  useRef,
  useState,
  useEffect,
} from 'react';
import {
  ActionIcon,
  Avatar,
  Box,
  Button,
  Flex,
  Group,
  Paper,
  Stack,
} from '@mantine/core';
import {
  IconArrowsMaximize,
  IconArrowsMinimize,
  IconMinus,
  IconX,
} from '@tabler/icons-react';
import { useAtomValue } from 'jotai';

import { type Asset } from 'lib/models';
import { getLocalBackend } from 'core/stores/settings.store';
import {
  currentViewer3dAssetsAtom,
  sidebarStateAtom,
  minimizeSidebar,
  maximizeSidebar,
  closeSidebar,
  removeFromViewer3d,
} from 'lib/stores/assetPage.store';

import { type Viewer3D, createViewer3D } from './Viewer3d';

interface IViewerListCard {
  readonly asset: Asset;
  readonly onRemove?: (asset: Asset) => void;
}

/**
 *
 * @param param0
 * @param param0.asset
 * @param param0.onRemove
 * @returns
 */
function ViewerListCard({ asset, onRemove }: IViewerListCard) {
  return (
    <Flex
      bg="dimmed"
      style={{ cursor: 'pointer' }}
      p="xs"
      gap="xs"
      justify="center"
      align="center"
      direction="row"
      wrap="nowrap"
    >
      <Avatar
        size="lg"
        radius="sm"
        src={`${getLocalBackend()}/api/lib/${asset.Thumbnail}/file`}
      />
      <Box flex={1}>
        {asset.Label}
      </Box>
      <Box>
        <ActionIcon
          variant="filled"
          aria-label="Settings"
          size="xl"
        >
          <IconMinus
            onClick={() => onRemove && onRemove(asset)}
            style={{ width: '70%', height: '70%' }}
            stroke={1.5}
          />
        </ActionIcon>
      </Box>
    </Flex>
  );
}

/**
 *
 * @returns
 */
export function Viewer3dPanel() {
  const viewer3dAssets = useAtomValue(currentViewer3dAssetsAtom);
  const sidebarState = useAtomValue(sidebarStateAtom);
  const parent = useRef<HTMLDivElement>(null);
  const [viewer3D, setViewer3D] = useState<Viewer3D>();

  useEffect(() => {
    if (viewer3dAssets.length > 0) {
      // eslint-disable-next-line no-console
      console.log('Viewer3dPanel useEffect', viewer3dAssets);
      minimizeSidebar();
    }
    if (!parent.current) {
      return;
    }
    if (!viewer3D) {
      return;
    }
    viewer3D.setModels(viewer3dAssets);
  }, [viewer3dAssets, viewer3D]);

  /*
  useEffect(() => {
  if (!parent.current) return;

  let viewer = createViewer3D(parent.current);
  setViewer3D(viewer);

  return () => {
  viewer.destroy();
  }
  }, [])
  */

  useEffect(() => {
    if (viewer3D) {
      return;
    }
    if (sidebarState === 'closed') {
      return;
    }
    if (!parent.current) {
      return;
    }

    const viewer = createViewer3D(parent.current);
    setViewer3D(viewer);

    // eslint-disable-next-line consistent-return
    return () => {
      if (viewer3D) {
        viewer.destroy();
        setViewer3D(undefined);
      }
    };
  }, [sidebarState, viewer3D]);

  return (
    <Paper
      display={sidebarState === 'closed' ? 'none' : 'flex'}
      m="sm"
      ml={0}
      flex={1}
    >
      <Stack
        display="flex"
        flex={1}
        dir="column"
        justify="flex-start"
        p="1rem"
      >
        <Group justify="space-between" gap="xs">
          <Box>3D Viewer</Box>
          <Group gap="xs" justify="flex-end">
            {
              sidebarState === 'maximized'
                ? null
                : (
                  <ActionIcon
                    variant="transparent"
                    aria-label="Maximize"
                    onClick={() => maximizeSidebar()}
                  >
                    <IconArrowsMaximize style={{ width: '60%', height: '60%' }} stroke={1.5} />
                  </ActionIcon>
                )
            }
            {
              sidebarState === 'minimized'
                ? null
                : (
                  <ActionIcon
                    variant="transparent"
                    aria-label="Minimize"
                    onClick={() => minimizeSidebar()}
                  >
                    <IconArrowsMinimize style={{ width: '60%', height: '60%' }} stroke={1.5} />
                  </ActionIcon>
                )
            }
            <ActionIcon
              variant="transparent"
              aria-label="Close"
              onClick={() => closeSidebar()}
            >
              <IconX style={{ width: '70%', height: '70%' }} stroke={1.5} />
            </ActionIcon>
          </Group>
        </Group>
        <Box ref={parent} flex={2} />
        <Group gap="xs" justify="space-between">
          <Box px="md">Models</Box>
          <Button variant="light">Clear</Button>
        </Group>
        <Stack align="stretch" flex={1} style={{ overflowY: 'auto' }}>
          <Stack>
            {viewer3dAssets.map((asset: Asset) => (
              <ViewerListCard
                key={asset.ID}
                asset={asset}
                onRemove={(a) => { removeFromViewer3d(a); }}
              />
            ))}
          </Stack>
        </Stack>
      </Stack>
    </Paper>
  );
}

export default Viewer3dPanel;
