import React, { createContext } from 'react';

import { type Asset } from '../models';

export interface AssetPageContextType {
  asset: Asset;
  setAsset: (value: Asset) => void;
  sidebarState: 'maximized' | 'minimized' | 'closed';
  viewer3dAssets: Asset[];
  maximizeSidebar: () => void;
  minimizeSidebar: () => void;
  closeSidebar: () => void;
  addToViewer3d: (asset: Asset) => void;
  removeFromViewer3d: (asset: Asset) => void;
}

export const AssetPageContext = createContext<AssetPageContextType | undefined>(undefined);

/**
 *
 * @param param0
 * @param param0.children
 * @returns
 */
export function AssetPageProvider({ children }: { readonly children: React.ReactNode }) {
  const [asset, setAsset] = React.useState<Asset | undefined>(undefined);
  const [sidebar, setSidebar] = React.useState<'maximized' | 'minimized' | 'closed'>('closed');
  const [viewer3dAssets, setViewer3dAssets] = React.useState<Asset[]>([]);
  // eslint-disable-next-line react/jsx-no-constructed-context-values
  const state: AssetPageContextType = {
    asset: asset!,
    setAsset,
    sidebarState: sidebar,
    maximizeSidebar: () => {
      setSidebar('maximized');
    },
    minimizeSidebar: () => {
      setSidebar('minimized');
    },
    closeSidebar: () => {
      setSidebar('closed');
    },
    viewer3dAssets,
    addToViewer3d: (newAsset: Asset) => {
      if (!newAsset || viewer3dAssets.some((a) => a.ID === newAsset.ID)) {
        return;
      }
      setViewer3dAssets([...viewer3dAssets, newAsset]);
    },
    removeFromViewer3d: (assetToRemove: Asset) => {
      const assetsToRemove = viewer3dAssets.filter((a) => a.ID !== assetToRemove.ID);
      if (assetsToRemove.length === 0) {
        return;
      }

      setViewer3dAssets(assetsToRemove);
    },
  };
  return <AssetPageContext.Provider value={state}>{children}</AssetPageContext.Provider>;
}

/**
 *
 * @returns
 */
export function useAssetPage() {
  const context = React.useContext(AssetPageContext);
  if (context === undefined) {
    throw new Error('useAssetPage must be used within a AssetPageProvider');
  }
  return context;
}
