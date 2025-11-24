import React from 'react';
import { IconGauge } from '@tabler/icons-react';

import { type NavItem } from 'navigation';

import { Dashboard } from './pages/dashboard/Dashboard';

const navItem: NavItem = {
  label: 'Dashboard',
  icon: IconGauge,
  path: '/dashboard',
  element: <Dashboard />,
};

export {
  navItem, // eslint-disable-line import/prefer-default-export
};
