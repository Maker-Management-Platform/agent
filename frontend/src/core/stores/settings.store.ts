import { atom } from 'jotai';

import { type Settings } from 'core/types/settings';

import jotaiStore from './jotai.store';

const settingsAtom = atom<Settings>({
  localBackend: 'http://localhost:8000',
  experimental: {
    dashboard: false,
  },
});

const isDashboardShown = () => (
  jotaiStore.get(settingsAtom).experimental.dashboard
);

const getLocalBackend = () => (
  jotaiStore.get(settingsAtom).localBackend
);

export {
  settingsAtom,
  isDashboardShown,
  getLocalBackend,
};
