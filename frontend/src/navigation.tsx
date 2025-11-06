import { navItem as dashNavItem } from './dashboard/routes.tsx';
import { navItem as libNavItem } from './lib/routes.tsx';
import App from "./App.tsx";
import { createBrowserRouter } from 'react-router';
import React from 'react';

export interface NavItem {
    label: string;
    icon: React.FC<any>;
    path: string;
    element?: React.ReactNode;
    initiallyOpened?: boolean;
    children?: {
        label: string;
        path: string;
        element: React.ReactNode;
    }[];
}

export const navigationItems: NavItem[] = [
    dashNavItem,
    libNavItem
];

export const router = createBrowserRouter([
    {
        path: "/",
        element: <App />,
        children: navigationItems,
    },
]);