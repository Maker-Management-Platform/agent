import { type Asset } from 'types/Asset';

export type Viewer3D = {
  destroy(): void;
  setModels(models: Asset[]): void;
};
