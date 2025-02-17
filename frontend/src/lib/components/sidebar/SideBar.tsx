import { Box } from "@mantine/core";
import { useAssetPage } from "../../contexts/AssetPageContext";
import classes from "./SideBar.module.css";

export function SideBar() {
    const pageState = useAssetPage();
    return (
        <Box className={classes.sidebar}>
            {pageState.asset?.Label}
        </Box>
    );
}