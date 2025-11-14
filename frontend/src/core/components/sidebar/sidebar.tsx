import React from 'react';
import { Code, Group, ScrollArea } from '@mantine/core';

import { type INavItem } from 'types/Nav';

import LinksGroup from './LinksGroup';
import classes from './sidebar.module.css';

/**
 *
 * @param param0
 * @param param0.navItems
 * @returns
 */
function Sidebar({ navItems }: { readonly navItems: INavItem[] }) {
  /* eslint-disable-next-line react/jsx-props-no-spreading */
  const links = navItems.map((item) => <LinksGroup {...item} key={item.label} />);

  return (
    <nav className={classes.navbar}>
      <div className={classes.header}>
        <Group justify="space-between">
          somethings
          <Code fw={700}>v2.0.0</Code>
        </Group>
      </div>

      <ScrollArea className={classes.links}>
        <div className={classes.linksInner}>{links}</div>
      </ScrollArea>
    </nav>
  );
}

export default Sidebar;
