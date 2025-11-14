import { useQuery, useQueryClient } from '@tanstack/react-query';

import { type Asset, type AssetPatch, type NewFolderPost } from 'types/Asset';

import backendAxiosInstance from './axios';

/**
 *
 * @param assetId
 * @returns
 */
const useGetAsset = (assetId?: string) => (
  useQuery<Asset>({
    queryKey: ['asset', assetId],
    queryFn: async () => {
      const { data } = await backendAxiosInstance.get(`lib/${assetId}`);
      return data;
    },
  })
);

/**
 *
 * @returns
 */
const useGetAssets = () => (
  useQuery<Asset>({
    queryKey: ['assets'],
    queryFn: async () => {
      const { data } = await backendAxiosInstance.get('lib');
      return data;
    },
  })
);

/**
 *
 * @returns
 */
const useNewFolderPost = () => {
  const queryClient = useQueryClient();

  return (input: NewFolderPost) => (
    backendAxiosInstance.post(`lib/${input.ParentID}/folder`)
      .then((response) => {
        queryClient.invalidateQueries({ queryKey: ['asset', input.ParentID] }).catch(console.error);
        queryClient.invalidateQueries({ queryKey: ['assets'] }).catch(console.error);
        return response;
      })
      .catch((error) => {
        queryClient.invalidateQueries({ queryKey: ['asset', input.ParentID] }).catch(console.error);
        queryClient.invalidateQueries({ queryKey: ['assets'] }).catch(console.error);
        return error;
      })
  );
};

/**
 *
 * @returns
 */
const usePatchAsset = () => {
  const queryClient = useQueryClient();

  return (asset: AssetPatch) => (
    backendAxiosInstance.patch(`lib/${asset.ID}`)
      .then((response) => {
        queryClient.invalidateQueries({ queryKey: ['asset', asset.ID] }).catch(console.error);
        return response;
      })
      .catch((error) => {
        queryClient.invalidateQueries({ queryKey: ['asset', asset.ID] }).catch(console.error);
        return error;
      })
  );
};

// /**
//  *
//  * @param assetId
//  * @returns
//  */
// async function getAsset(assetId?: string): Promise<Asset> {
//   const assetPath = assetId ? `/${assetId}` : '';
//   const response = await fetch(
//     `http://localhost:8000/api/lib${assetPath}`,
//   );
//   if (!response.ok) {
//     throw new Error('Network response was not ok');
//   }
//   return response.json();
// }

/**
 *
 * @param assetId
 * @returns
 */
// function useGetAsset(assetId?: string) {
//   return useQuery<Assets[]>({
//     queryKey: ['asset', assetId],
//     queryFn: async (): Promise<Asset> => {
//       const assetPath = assetId ? `/${assetId}` : '';
//       const response = await fetch(
//         `http://localhost:8000/api/lib${assetPath}`,
//       );
//       if (!response.ok) {
//         throw new Error('Network response was not ok');
//       }
//       return response.json();
//     },
//   });
// }

/**
 *
 * @param asset
 * @returns
 */
// async function patchAsset(asset: AssetPatch): Promise<Asset> {
//   const response = await fetch(`http://localhost:8000/api/lib/${asset.ID}`, {
//     method: 'PATCH',
//     headers: {
//       'Content-Type': 'application/json',
//     },
//     body: JSON.stringify(asset),
//   });
//   if (!response.ok) {
//     throw new Error('Network response was not ok');
//   }
//   return response.json();
// }

// /**
//  *
//  * @param input
//  * @returns
//  */
// async function newFolderPost(): Promise<Asset> {
//   const queryClient = useQueryClient();

//   return (input: NewFolderPost) => (
//     fetch(`http://localhost:8000/api/lib/${input.ParentID}/folder`,
//       {
//         method: 'POST',
//         headers: {
//           'Content-Type': 'application/json',
//         },
//         body: JSON.stringify(input),
//       }
//     ).then((response) = {
//       queryClient.invalidateQueries({ queryKey: ['asset', input.ParentID] })
// .catch(console.error);

//       if (!response.ok) {
//         throw new Error('Network response was not ok');
//       }

//       return response.json();
//     }).catch(() => {
//       queryClient.invalidateQueries({ queryKey: ['asset', input.ParentID] })
// .catch(console.error);

//       throw new Error('Network response was not ok');
//     })
//   );
// };

export {
  // getAsset,
  useNewFolderPost,
  usePatchAsset,
  useGetAsset,
  useGetAssets,
};
