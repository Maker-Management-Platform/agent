import { ActionIcon, Box, rem, Title } from "@mantine/core";
import classes from "./AssetCard.module.css";
import { Asset } from "../../models";
import { Link } from "react-router";
import { useHover } from "@mantine/hooks";
import { IconHeart, IconSettings } from "@tabler/icons-react";
import { useContext } from "react";
import { SettingsContext } from "@/core/providers/settings/settingsContext";

export function AssetCard({
  asset,
  onSelected,
  onAddTo3dViewer,
}: {
  asset: Asset;
  onSelected?: (asset: Asset) => void;
  onAddTo3dViewer?: (asset: Asset) => void;
}) {
  const { hovered, ref } = useHover();
  const { settings } = useContext(SettingsContext);
  const shouldNavigate = (a: Asset) => {
    return (
      a.NodeKind === "bundle" || a.NodeKind === "dir" || a.NodeKind === "root"
    );
  };

  return (
    <Box className={classes.card} ref={ref} m="0.5rem">
      <div className={classes.content}>
        {shouldNavigate(asset) ? (
          <Title order={5} component={Link} to={`/lib/${asset.ID}`}>
            {asset.Label}
          </Title>
        ) : (
          <Title order={5}>{asset.Label}</Title>
        )}
      </div>
      <div
        className={classes.controls}
        style={{ display: hovered ? "flex" : "none" }}
      >
        <ActionIcon.Group>
          {asset.Kind === "model" && (
            <ActionIcon
              variant="white"
              size="sm"
              aria-label="Settings"
              onClick={
                onAddTo3dViewer ? () => onAddTo3dViewer(asset) : undefined
              }
            >
              <IconSettings style={{ width: rem(20) }} stroke={1.5} />
            </ActionIcon>
          )}
          <ActionIcon
            variant="white"
            size="sm"
            aria-label="Likes"
            onClick={onSelected ? () => onSelected(asset) : undefined}
          >
            <IconHeart style={{ width: rem(20) }} stroke={1.5} />
          </ActionIcon>
        </ActionIcon.Group>
      </div>
      <div
        style={{
          backgroundImage: `url(${settings.localBackend}/api/lib/${asset.Thumbnail}/file)`,
        }}
        className={classes.thumbnail}
      >
        <div className={classes.overlay} />
      </div>
    </Box>
  );
}

