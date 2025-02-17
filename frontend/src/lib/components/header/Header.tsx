import { ActionIcon, LoadingOverlay, Overlay, Text, Title, TypographyStylesProvider, rem } from '@mantine/core';
import classes from './Header.module.css';
import { IconExternalLink } from '@tabler/icons-react';

type HeaderProps = {
    loading?: boolean
    title?: string
    description?: string
    imgID?: string
    link?: string
}

export function Header({ title, description, loading, imgID, link }: HeaderProps) {
    const fallbackImage = 'https://images.unsplash.com/photo-1563520239648-a24e51d4b570?q=80&w=2000&h=400&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D';
    const img = imgID ? `http://localhost:8000/api/lib/${imgID}/file` : fallbackImage;

    return (
        <div
            className={classes.wrapper}
            style={{ backgroundImage: `url(${loading ? fallbackImage : img})` }}
        >
            <Overlay color="#000" opacity={0.65} zIndex={1} />

            <LoadingOverlay visible={loading} zIndex={1000} overlayProps={{ blur: 2 }} />
            <div className={classes.inner}>
                <Title className={classes.title}>
                    {title}
                </Title>

                <Text size="lg" className={classes.description} lineClamp={3} component="div">
                    <TypographyStylesProvider>
                        <p>{description}</p>
                    </TypographyStylesProvider>
                </Text>

                {link && <ActionIcon className={classes.link} color='white' variant="subtle" size="lg" aria-label="Link" component="a" href={link} target='_blank'>
                    <IconExternalLink style={{ width: rem(20) }} stroke={1.5} />
                </ActionIcon>}
            </div>
        </div>
    );
}