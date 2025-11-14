export interface Asset {
  ID: string;
  Label: string;
  Description: string;
  Path: string;
  Root: string;
  FSKind: string;
  FSName: string;
  Extension: string;
  Kind: string;
  NodeKind: string;
  ParentID: string;
  Parent: Asset;
  NestedAssets: Asset[];
  Thumbnail: string;
  SeenOnScan: boolean;
  Properties: Map<string, string>;
  Tags: string[];
  CreatedAt: Date;
  UpdatedAt: Date;
}

export interface AssetPatch {
  ID: string;
  Label?: string;
  Description?: string;
  Tags?: string[];
}

export interface NewFolderPost {
  ParentID: string;
  FolderName: string;
}
