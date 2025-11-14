import { atom } from 'jotai';

import { type Asset } from 'types/Asset';
import { type SideBarPosition } from 'types/AssetPage';

import jotaiStore from './jotai.store';

const currentAssetAtom = atom<Asset | undefined>();
const sidebarStateAtom = atom<SideBarPosition>('closed');

const viewer3dAssetsAtom = atom<Asset[]>([]);

const maximizeSidebar = () => {
  jotaiStore.set(sidebarStateAtom, 'maximized');
};

const minimizeSidebar = () => {
  jotaiStore.set(sidebarStateAtom, 'minimized');
};

const closeSidebar = () => {
  jotaiStore.set(sidebarStateAtom, 'closed');
};

export {
  currentAssetAtom,
  sidebarStateAtom,
  viewer3dAssetsAtom,
  maximizeSidebar,
  minimizeSidebar,
  closeSidebar,
};
