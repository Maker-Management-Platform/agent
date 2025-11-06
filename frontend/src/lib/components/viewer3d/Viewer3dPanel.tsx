import { ActionIcon, Avatar, Box, Button, Flex, Group, Paper, Stack } from "@mantine/core";
import { useRef, useState, useEffect, useContext } from "react";
import { useAssetPage } from "../../contexts/AssetPageContext";
import { Viewer3D, createViewer3D } from "./Viewer3d";
import { Asset } from "@/lib/models";
import { IconArrowsMaximize, IconArrowsMinimize, IconMinus, IconX } from "@tabler/icons-react";
import { SettingsContext } from "@/core/providers/settings/settingsContext";

export function Viewer3dPanel() {
    const assetPage = useAssetPage();
    const parent = useRef<HTMLElement>();
    const [viewer3D, setViewer3D] = useState<Viewer3D>();

    useEffect(() => {
        if (assetPage.viewer3dAssets.length !== 0) {
            console.log("Viewer3dPanel useEffect", assetPage.viewer3dAssets);
            assetPage.minimizeSidebar();
        }
        if (!parent.current) return;
        if (!viewer3D) return;
        viewer3D.setModels(assetPage.viewer3dAssets);

    }, [assetPage.viewer3dAssets, viewer3D]);

    /*useEffect(() => {
        if (!parent.current) return;

        let viewer = createViewer3D(parent.current);
        setViewer3D(viewer);

        return () => {
            viewer.destroy();
        }
    }, [])*/

    useEffect(() => {
        if (viewer3D) return;
        if (assetPage.sidebarState === "closed") return;
        if (!parent.current) return;

        let viewer = createViewer3D(parent.current);
        setViewer3D(viewer);

        return () => {
            if (viewer3D) {
                viewer.destroy();
                setViewer3D(undefined);
            }
        }
    }, [assetPage.sidebarState])


    return (
        <Paper display={assetPage.sidebarState === "closed" ? "none" : "flex"
        } m="sm" ml={0} flex={1} >

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
                        {assetPage.sidebarState !== "maximized" && <ActionIcon
                            variant="transparent"
                            aria-label="Maximize"
                            onClick={() => assetPage.maximizeSidebar()}>
                            <IconArrowsMaximize style={{ width: '60%', height: '60%' }} stroke={1.5} />
                        </ActionIcon>}
                        {assetPage.sidebarState !== "minimized" && <ActionIcon
                            variant="transparent"
                            aria-label="Minimize"
                            onClick={() => assetPage.minimizeSidebar()}>
                            <IconArrowsMinimize style={{ width: '60%', height: '60%' }} stroke={1.5} />
                        </ActionIcon>}
                        <ActionIcon
                            variant="transparent"
                            aria-label="Close"
                            onClick={() => assetPage.closeSidebar()}>
                            <IconX style={{ width: '70%', height: '70%' }} stroke={1.5} />
                        </ActionIcon>
                    </Group>
                </Group>
                <Box ref={parent} flex={2} />
                <Group gap="xs" justify="space-between">
                    <Box px="md">Models</Box>
                    <Button variant="light">Clear</Button>
                </Group>
                <Stack align="stretch" flex={1} style={{ overflowY: "auto" }}>
                    <Stack>
                        {assetPage.viewer3dAssets.map(asset => (
                            <ViewerListCard key={asset.ID} asset={asset} onRemove={(a) => { assetPage.removeFromViewer3d(a) }} />
                        ))}
                    </Stack>
                </Stack>
            </Stack>
        </Paper >
    );
}

function ViewerListCard({ asset, onRemove }: { asset: Asset, onRemove?: (asset: Asset) => void }) {
    const { settings } = useContext(SettingsContext)
    return (
        <Flex bg="dimmed"
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
                src={`${settings.localBackend}/api/lib/${asset.Thumbnail}/file`}
            />
            <Box flex={1}>{asset.Label}</Box>
            <Box>
                <ActionIcon
                    variant="filled"
                    aria-label="Settings"
                    size="xl">
                    <IconMinus
                        onClick={() => onRemove && onRemove(asset)}
                        style={{ width: '70%', height: '70%' }} stroke={1.5} />
                </ActionIcon>
            </Box>
        </Flex>
    );
}