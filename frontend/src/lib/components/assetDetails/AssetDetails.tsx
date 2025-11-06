import { Asset } from "../../models";
import { Box } from "@mantine/core";

export function AssetDetails({ asset }: { asset: Asset }) {
    return (
        <Box>
            <Box>{asset.Label}</Box>
        </Box>
    );
}