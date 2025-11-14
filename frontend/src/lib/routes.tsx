import React from 'react';
import { IconBook } from '@tabler/icons-react';

import { type INavItem } from 'types/Nav';

import LibMain from './pages/libMain/LibMain';
import LibSettings from './pages/libSettings/LibSettings';
import Asset from './components/Asset';

const libNavItem: INavItem = {
  label: 'Libraries',
  icon: IconBook,
  path: '/lib',
  element: <LibMain />,
  children: [
    {
      label: 'Libraries',
      path: ':id?',
      element: <Asset />,
    },
    {
      label: 'Library',
      path: '/lib/:id',
      element: <LibSettings />,
    },
  ],
};

export default libNavItem;
