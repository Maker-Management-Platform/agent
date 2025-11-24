export interface ExperimentalFeatures {
  dashboard: boolean
}

export interface Settings {
  localBackend: string,
  agent?: Record<string, string>,
  experimental: ExperimentalFeatures
}
