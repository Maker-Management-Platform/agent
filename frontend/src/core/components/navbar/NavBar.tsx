import React from 'react';
import cx from 'clsx';
import {
  Center,
  Tooltip,
  UnstyledButton,
  Stack,
  useMantineColorScheme,
  useComputedColorScheme,
  Code,
} from '@mantine/core';
import {
  IconSun,
  IconMoon,
  IconBrandMantine,
} from '@tabler/icons-react';
import { NavLink } from 'react-router';

import { type NavItem } from 'navigation';

import classes from './NavBar.module.css';

/**
 *
 * @param param0
 * @param param0.icon
 * @param param0.label
 * @param param0.path
 * @returns
 */
// eslint-disable-next-line sonarjs/prefer-read-only-props
function NavbarLink({ icon: Icon, label, path }: NavItem) {
  return (
    <Tooltip label={label} position="right" transitionProps={{ duration: 0 }}>
      <UnstyledButton
        className={classes.link}
        renderRoot={({ className, ...others }) => (
          <NavLink
            to={path}
            className={({ isActive }) => cx(className, isActive && classes.active)}
            // eslint-disable-next-line react/jsx-props-no-spreading
            {...others}
          />
        )}
      >
        <Icon stroke={1.5} />
      </UnstyledButton>
    </Tooltip>
  );
}

// const operationalItems: NavItem[] = [];

/**
 *
 * @param param0
 * @param param0.navItems
 * @returns
 */
function NavBar({ navItems }: { readonly navItems: NavItem[] }) {
  const { setColorScheme } = useMantineColorScheme();
  const computedColorScheme = useComputedColorScheme('light', { getInitialValueInEffect: true });

  const featureLinks = navItems.map((link) => (
    <NavbarLink
      // eslint-disable-next-line react/jsx-props-no-spreading
      {...link}
      key={link.label}
    />
  ));

  // const opsLinks = operationalItems.map((link) => (
  //   <NavbarLink
  //     // eslint-disable-next-line react/jsx-props-no-spreading
  //     {...link}
  //     key={link.label}
  //   />
  // ));

  return (
    <nav className={classes.navbar}>
      <Center>
        <IconBrandMantine type="mark" size={30} />
      </Center>

      <div className={classes.navbarMain}>
        <Stack justify="center" gap={0}>
          {featureLinks}
        </Stack>
      </div>

      <Stack justify="center" gap={0}>
        {/* {opsLinks} */}
        <Tooltip label="Toggle color scheme" position="right" transitionProps={{ duration: 0 }}>
          <UnstyledButton className={classes.link} onClick={() => setColorScheme(computedColorScheme === 'light' ? 'dark' : 'light')}>
            {computedColorScheme === 'dark' && <IconSun stroke={1.5} />}
            {computedColorScheme === 'light' && <IconMoon stroke={1.5} />}
          </UnstyledButton>
        </Tooltip>
        {/* <NavbarLink icon={IconSwitchHorizontal} href={'change'} label="Change account" /> */}
        {/* <NavbarLink icon={IconLogout} href={'logout'} label="Logout" /> */}

      </Stack>
      <Code fw={700}>v2.0.0</Code>
    </nav>
  );
}

export default NavBar;
