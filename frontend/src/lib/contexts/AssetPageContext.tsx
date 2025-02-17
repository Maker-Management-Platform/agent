import { createContext } from "react";
import { Asset } from "../models";
import React from "react";

export interface AssetPageContextType {
    asset: Asset;
    setAsset: (value: Asset) => void;
    sidebarState: "maximized" | "minimized" | "closed";
    viewer3dAssets: Asset[];
    maximizeSidebar: () => void;
    minimizeSidebar: () => void;
    closeSidebar: () => void;
    addToViewer3d: (asset: Asset) => void;
    removeFromViewer3d: (asset: Asset) => void;
}

export const AssetPageContext = createContext<AssetPageContextType | undefined>(undefined);

export function AssetPageProvider({ children }: { children: React.ReactNode }) {
    const [asset, setAsset] = React.useState<Asset | undefined>(undefined);
    const [sidebar, setSidebar] = React.useState<"maximized" | "minimized" | "closed">("closed");
    const [viewer3dAssets, setViewer3dAssets] = React.useState<Asset[]>([]);
    const state: AssetPageContextType = {
        asset: asset!,
        setAsset: setAsset,
        sidebarState: sidebar,
        maximizeSidebar: () => {
            setSidebar("maximized");
        },
        minimizeSidebar: () => {
            setSidebar("minimized");
        },
        closeSidebar: () => {
            setSidebar("closed");
        },
        viewer3dAssets: viewer3dAssets,
        addToViewer3d: (asset: Asset) => {
            if (viewer3dAssets.find((a) => a.ID === asset.ID)) return;
            setViewer3dAssets([...viewer3dAssets, asset]);
        },
        removeFromViewer3d: (asset: Asset) => {
            setViewer3dAssets(viewer3dAssets.filter((a) => a.ID !== asset.ID));
        }
    }
    return <AssetPageContext.Provider value={state}>{children}</AssetPageContext.Provider>;
}


export function useAssetPage() {
    const context = React.useContext(AssetPageContext);
    if (context === undefined) {
        throw new Error('useAssetPage must be used within a AssetPageProvider');
    }
    return context;
}