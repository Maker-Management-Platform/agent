import { atom } from 'jotai';

import { type Settings } from 'types/Settings';

import jotaiStore from './jotai.store';

const settingsAtom = atom<Settings>({
  localBackend: 'http://localhost:8000',
  experimental: {
    dashboard: false,
  },
});

const setSettings = (settings: Settings) => {
  jotaiStore.set(settingsAtom, settings);
};

const getLocalBackend = () => (
  jotaiStore.get(settingsAtom).localBackend ?? 'http://localhost:8000'
);

export {
  settingsAtom,
  setSettings,
  getLocalBackend,
};
