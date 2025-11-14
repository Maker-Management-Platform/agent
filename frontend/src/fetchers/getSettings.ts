import { useQuery } from '@tanstack/react-query';

import { type Settings } from 'types/Settings';
import { setSettings } from 'stores/Settings';

/**
 *
 * @returns
 */
const useGetSettings = () => (
  useQuery<Settings>({
    queryKey: ['settings'],
    queryFn: async () => {
      const response = await fetch('/settings.json');
      const settings = await response.json();

      setSettings(settings as Settings);

      return settings;
    },
  })
);

export {
  useGetSettings, // eslint-disable-line import/prefer-default-export
};
