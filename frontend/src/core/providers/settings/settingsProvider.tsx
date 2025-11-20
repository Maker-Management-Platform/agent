import React, { useEffect, useState } from 'react';
import { useLocalStorage } from '@mantine/hooks';

import SettingsContext, { type ExperimentalFeatures, type Settings } from './settingsContext';

interface ISettingProvider {
  readonly loading: React.ReactNode,
  readonly children: React.ReactNode,
}

/**
 *
 * @param param0
 * @param param0.loading
 * @param param0.children
 * @returns
 */
function SettingsProvider({ loading, children }: ISettingProvider) {
  const [settings, setSettings] = useState<Settings>({} as Settings);
  const [ready, setReady] = useState(false);

  const [experimental, setExperimental] = useLocalStorage<ExperimentalFeatures>({
    key: 'experimental',
    defaultValue: {
      dashboard: false,
    },
  });

  useEffect(() => {
    setSettings((prev) => ({ ...prev, experimental }));
  }, [experimental]);

  useEffect(() => {
    fetch('/settings.json')
      .then((response) => response.json())
      .then((data) => setSettings((prev) => ({ ...prev, ...data })))
      .then(() => setReady(true))
      .catch(console.error);
  }, []);

  return (
    // eslint-disable-next-line react/jsx-no-constructed-context-values
    <SettingsContext.Provider value={{ settings, setExperimental }}>
      {ready ? children : loading}
    </SettingsContext.Provider>
  );
}

export default SettingsProvider;
