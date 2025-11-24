import { navItem as dashNavItem } from './dashboard/routes';
import { navItem as libNavItem } from './lib/routes';

export interface NavItem {
  readonly label: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  readonly icon: React.FC<any>;
  readonly path: string;
  readonly element?: React.ReactNode;
  readonly initiallyOpened?: boolean;
  readonly children?: {
    readonly label: string;
    readonly path: string;
    readonly element: React.ReactNode;
  }[];
}

export const navigationItems: NavItem[] = [
  dashNavItem,
  libNavItem,
];
