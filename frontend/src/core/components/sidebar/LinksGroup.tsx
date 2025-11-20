import React, { useState } from 'react';
import { Link } from 'react-router';
import { IconChevronRight } from '@tabler/icons-react';
import {
  Box,
  Collapse,
  Group,
  Text,
  ThemeIcon,
  UnstyledButton,
} from '@mantine/core';

import classes from './NavbarLinksGroup.module.css';

interface LinksGroupProps {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  readonly icon: React.FC<any>;
  readonly label: string;
  readonly initiallyOpened?: boolean;
  readonly path?: string;
  readonly children?: { label: string; path: string }[];
}

/**
 *
 * @param param0
 * @param param0.icon
 * @param param0.label
 * @param param0.path
 * @param param0.initiallyOpened
 * @param param0.children
 * @returns
 */
function LinksGroup({
  icon: Icon,
  label,
  path,
  initiallyOpened,
  children,
}: LinksGroupProps) {
  const hasLinks = Array.isArray(children);
  const [opened, setOpened] = useState(initiallyOpened || false);

  const items = (hasLinks ? children : []).map((link) => (
    <Text
      className={classes.link}
      key={link.label}
      // eslint-disable-next-line react/jsx-props-no-spreading
      renderRoot={(props) => <Link {...props} to={path + link.path} />}
    >
      {link.label}
    </Text>
  ));

  return (
    <>
      <UnstyledButton onClick={() => setOpened((o) => !o)} className={classes.control}>
        <Group justify="space-between" gap={0}>
          <Box style={{ display: 'flex', alignItems: 'center' }}>
            <ThemeIcon variant="light" size={30}>
              <Icon size={18} />
            </ThemeIcon>
            {path
              ? (
                <Box
                  ml="md"
                  // eslint-disable-next-line react/jsx-props-no-spreading
                  renderRoot={(props) => <Link {...props} to={path} />}
                >
                  {label}
                </Box>
              )
              : <Box ml="md">{label}</Box>}
          </Box>
          {hasLinks && (
            <IconChevronRight
              className={classes.chevron}
              stroke={1.5}
              size={16}
              style={{ transform: opened ? 'rotate(-90deg)' : 'none' }}
            />
          )}
        </Group>
      </UnstyledButton>
      {hasLinks ? <Collapse in={opened}>{items}</Collapse> : null}
    </>
  );
}

export default LinksGroup;
