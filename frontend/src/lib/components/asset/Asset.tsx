import { Link, useLocation, useParams } from "react-router";
import { useGetAsset } from "../../fetchers/getAsset";
import { AssetCard } from "../assetCard/AssetCard";
import {
  Anchor,
  Box,
  Breadcrumbs,
  Group,
  Collapse,
  ActionIcon,
  Paper,
} from "@mantine/core";
import classes from "./Asset.module.css";
import { Header } from "../header/Header";
import { useAssetPage } from "../../contexts/AssetPageContext";
import { useState } from "react";
import {
  IconCircleDashedPlus,
  IconEdit,
  IconInfoCircle,
} from "@tabler/icons-react";
import { AssetEditForm } from "../assetEditForm/AssetEditForm";
import { AssetNewForm } from "../assetNewForm/AssetNewForm";

export function Asset() {
  const loc = useLocation();
  const { id } = useParams();

  const assetPageState = useAssetPage();

  const [expandedView, setExpandedView] = useState(loc.hash.substring(1));

  const { data, isLoading } = useGetAsset(id);
  const breadcrumbs = (function calculateBreadcrumbs(asset): JSX.Element[] {
    if (asset && asset.Parent) {
      return calculateBreadcrumbs(asset.Parent).concat(
        <Anchor
          component={Link}
          to={`/lib/${asset.Parent.ID}`}
          key={asset.Parent.ID}
        >
          {asset.Parent.Label}
        </Anchor>,
      );
    } else {
      return [
        <Anchor key="home" component={Link} to={`/lib`}>
          Home
        </Anchor>,
      ];
    }
  })(data);

  const toggleExpandedView = (view: string) =>
    setExpandedView(() => {
      window.location.hash = view;
      if (expandedView === view) {
        window.location.hash = "";
        return "";
      }
      return view;
    });

  return (
    <Paper className={classes.asset} m="sm" display="flex">
      <Header
        title={data?.Label}
        description={data?.Description}
        loading={isLoading}
        imgID={data?.Thumbnail}
      />
      <Group
        justify="space-between"
        preventGrowOverflow={false}
        mt="xs"
        mx="sm"
      >
        <Breadcrumbs separator="→" separatorMargin="sm">
          {breadcrumbs}
        </Breadcrumbs>
        <Group justify="flex-end" wrap="nowrap">
          <ActionIcon
            variant={expandedView == "details" ? "light" : "filled"}
            aria-label="Details"
            onClick={() => toggleExpandedView("details")}
          >
            <IconInfoCircle
              style={{ width: "70%", height: "70%" }}
              stroke={1.5}
            />
          </ActionIcon>
          <ActionIcon
            variant={expandedView == "edit" ? "light" : "filled"}
            aria-label="Settings"
            onClick={() => toggleExpandedView("edit")}
          >
            <IconEdit style={{ width: "70%", height: "70%" }} stroke={1.5} />
          </ActionIcon>
          <ActionIcon
            variant={expandedView == "add" ? "light" : "filled"}
            aria-label="Settings"
            onClick={() => toggleExpandedView("add")}
          >
            <IconCircleDashedPlus
              style={{ width: "70%", height: "70%" }}
              stroke={1.5}
            />
          </ActionIcon>
        </Group>
      </Group>

      <Box className={classes.content} mt="sm">
        <Collapse in={expandedView === "details"}>
          <Box>a thing</Box>
        </Collapse>
        <Collapse in={expandedView === "edit"}>
          {data && (
            <AssetEditForm asset={data} onClose={() => setExpandedView("")} />
          )}
        </Collapse>
        <Collapse in={expandedView === "add"}>
          {data && (
            <AssetNewForm parent={data} onClose={() => setExpandedView("")} />
          )}
        </Collapse>
      </Box>

      <Box className={classes.list} m="sm">
        {data &&
          data?.NestedAssets.map((a) => (
            <AssetCard
              asset={a}
              key={a.ID}
              onSelected={(asset) => {
                assetPageState.setAsset(asset);
              }}
              onAddTo3dViewer={(asset) => {
                assetPageState.setAsset(asset);
                assetPageState.addToViewer3d(asset);
              }}
            />
          ))}
      </Box>
    </Paper>
  );
}
