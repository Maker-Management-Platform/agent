import React from 'react';
import { IconBook } from '@tabler/icons-react';

import { type NavItem } from 'navigation';

import LibMain from './pages/libMain/LibMain';
import LibSettings from './pages/libSettings/LibSettings';
import Asset from './components/asset/Asset';

const navItem: NavItem = {
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

export {
  navItem, // eslint-disable-line import/prefer-default-export
};
