export interface INavItemChild {
  readonly label: string;
  readonly path: string;
  readonly element: React.ReactNode;
  readonly children?: [] // only needed for recursive route generation
}

export interface INavItem {
  readonly label: string;
  readonly icon: React.FC<any>; // eslint-disable-line @typescript-eslint/no-explicit-any
  readonly path: string;
  readonly element?: React.ReactNode;
  readonly initiallyOpened?: boolean;
  readonly children?: INavItemChild[];
}
