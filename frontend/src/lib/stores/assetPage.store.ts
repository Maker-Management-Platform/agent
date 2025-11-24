import { atom } from 'jotai';

import { type Asset } from 'lib/models';
import jotaiStore from 'core/stores/jotai.store';

const currentAssetAtom = atom<Asset | undefined>();

const sidebarStateAtom = atom<'maximized' | 'minimized' | 'closed'>('closed');

const currentViewer3dAssetsAtom = atom<Asset[]>([]);

/**
 *
 * @param value
 * @returns
 */
const setAsset = (value: Asset) => (
  jotaiStore.set(currentAssetAtom, value)
);

/**
 *
 * @returns
 */
const maximizeSidebar = () => (
  jotaiStore.set(sidebarStateAtom, 'maximized')
);

/**
 *
 * @returns
 */
const minimizeSidebar = () => (
  jotaiStore.set(sidebarStateAtom, 'minimized')
);

/**
 *
 * @returns
 */
const closeSidebar = () => (
  jotaiStore.set(sidebarStateAtom, 'closed')
);

/**
 *
 * @returns
 */
const getViewer3dAssets = () => (
  jotaiStore.get(currentViewer3dAssetsAtom)
);

/**
 *
 * @param newAsset
 */
const addToViewer3d = (newAsset: Asset) => {
  const { ID } = newAsset;
  const viewer3dAssets = getViewer3dAssets();
  if (!newAsset && !viewer3dAssets.some((a) => a.ID === ID)) {
    jotaiStore.set(currentViewer3dAssetsAtom, [...viewer3dAssets, newAsset]);
  }
};

/**
 *
 * @param assetToRemove
 */
const removeFromViewer3d = (assetToRemove: Asset) => {
  const viewer3dAssets = getViewer3dAssets();
  const assetsToRemove = viewer3dAssets.filter((a) => a.ID !== assetToRemove.ID);
  jotaiStore.set(currentViewer3dAssetsAtom, assetsToRemove);
};

export {
  currentViewer3dAssetsAtom,
  sidebarStateAtom,
  currentAssetAtom,
  setAsset,
  maximizeSidebar,
  minimizeSidebar,
  closeSidebar,
  addToViewer3d,
  removeFromViewer3d,
};
