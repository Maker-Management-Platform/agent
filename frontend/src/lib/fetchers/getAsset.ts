import { useQuery } from '@tanstack/react-query';

import { type Asset, type AssetPatch, type NewFolderPost } from '../models';

/**
 *
 * @param assetId
 * @returns
 */
export async function getAsset(assetId?: string): Promise<Asset> {
  const assetPath = assetId ? `/${assetId}` : '';
  const response = await fetch(
    `http://localhost:8000/api/lib${assetPath}`,
  );
  if (!response.ok) {
    throw new Error('Network response was not ok');
  }
  return response.json();
}

/**
 *
 * @param assetId
 * @returns
 */
export function useGetAsset(assetId?: string) {
  return useQuery({
    queryKey: ['asset', assetId],
    queryFn: async (): Promise<Asset> => {
      const assetPath = assetId ? `/${assetId}` : '';
      const response = await fetch(
        `http://localhost:8000/api/lib${assetPath}`,
      );
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      return response.json();
    },
  });
}

/**
 *
 * @param asset
 * @returns
 */
export async function patchAsset(asset: AssetPatch): Promise<Asset> {
  const response = await fetch(`http://localhost:8000/api/lib/${asset.ID}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(asset),
  });
  if (!response.ok) {
    throw new Error('Network response was not ok');
  }
  return response.json();
}

/**
 *
 * @param input
 * @returns
 */
export async function newFolderPost(input: NewFolderPost): Promise<Asset> {
  const response = await fetch(
    `http://localhost:8000/api/lib/${input.ParentID}/folder`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(input),
    },
  );
  if (!response.ok) {
    throw new Error('Network response was not ok');
  }
  return response.json();
}
