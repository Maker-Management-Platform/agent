import React from 'react';
import ReactDOM from 'react-dom/client';
import { createTheme, MantineProvider } from '@mantine/core';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { RouterProvider, createBrowserRouter } from 'react-router';

import { navigationItems } from 'navigation';
import App from 'App';

const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: navigationItems,
  },
]);

const queryClient = new QueryClient();
const theme = createTheme({
  autoContrast: true,
  cursorType: 'pointer',
  headings: {
    fontFamily: 'Roboto, sans-serif',
  },
});
ReactDOM.createRoot(document.querySelector('#root') as HTMLElement).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <MantineProvider defaultColorScheme="dark" theme={theme}>
        <RouterProvider router={router} />
      </MantineProvider>
    </QueryClientProvider>
  </React.StrictMode>,
);
