import { LinksGroup } from './LinksGroup';
import { Code, Group, ScrollArea, } from '@mantine/core';
import classes from './sidebar.module.css';
import { NavItem } from "@/navigation";

export function Sidebar({ navItems }: { navItems: NavItem[] }) {
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