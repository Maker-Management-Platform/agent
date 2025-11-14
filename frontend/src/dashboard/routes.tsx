import React from 'react';
import { IconGauge } from '@tabler/icons-react';

import { type INavItem } from 'types/Nav';

import Dashboard from './pages/dashboard/Dashboard';

const dashNavItem: INavItem = {
  label: 'Dashboard',
  icon: IconGauge,
  path: '/dashboard',
  element: <Dashboard />,
};

export default dashNavItem;
