import React from 'react';

export interface Settings {
  localBackend: string,
  agent?: Record<string, string>,
  experimental: ExperimentalFeatures
}
export interface ExperimentalFeatures {
  dashboard: boolean
}
interface SettingsProviderType {
  settings: Settings,
  setExperimental: (mutator: (prev: ExperimentalFeatures) => ExperimentalFeatures) => void,
}

const SettingsContext = React.createContext<SettingsProviderType>({} as SettingsProviderType);

export default SettingsContext;
