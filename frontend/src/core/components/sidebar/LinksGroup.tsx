import { useState } from 'react';
import { IconChevronRight } from '@tabler/icons-react';
import { Box, Collapse, Group, Text, ThemeIcon, UnstyledButton } from '@mantine/core';
import classes from './NavbarLinksGroup.module.css';
import { Link } from 'react-router';

interface LinksGroupProps {
  icon: React.FC<any>;
  label: string;
  initiallyOpened?: boolean;
  path?: string;
  children?: { label: string; path: string }[];
}

export function LinksGroup({ icon: Icon, label, path, initiallyOpened, children }: LinksGroupProps) {
  const hasLinks = Array.isArray(children);
  const [opened, setOpened] = useState(initiallyOpened || false);


  const items = (hasLinks ? children : []).map((link) => (
    <Text
      className={classes.link}
      key={link.label}
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
            {path ? (
              <Box ml="md"
                renderRoot={(props) => <Link {...props} to={path} />}
              >{label}</Box>) : <Box ml="md">{label}</Box>}
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
