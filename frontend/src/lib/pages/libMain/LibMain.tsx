import { Outlet } from "react-router";
import classes from "./LibMain.module.css";
import { AssetPageProvider } from "@/lib/contexts/AssetPageContext";
import { Viewer3dPanel } from "@/lib/components/viewer3d/Viewer3dPanel";
import { Box } from "@mantine/core";


export function LibMain() {
    return (
        <Box className={classes.main} >
            <AssetPageProvider>
                <Outlet />
                <Viewer3dPanel />
            </AssetPageProvider>
        </Box>
    );
}