import React from 'react';
import ReactDOM from 'react-dom/client';
import { RouterProvider, createBrowserRouter } from 'react-router';
import { createTheme, MantineProvider } from '@mantine/core';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Provider as JotaiProvider } from 'jotai';

import jotaiStore from 'stores/jotai.store';
import dashNavItem from 'dashboard/routes';
import libNavItem from 'lib/routes';

import App from './App';

const queryClient = new QueryClient();

const theme = createTheme({
  autoContrast: true,
  cursorType: 'pointer',
  headings: {
    fontFamily: 'Roboto, sans-serif',
  },
});

const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      dashNavItem,
      libNavItem,
    ],
  },
]);

ReactDOM.createRoot(document.querySelector('#root') as HTMLElement).render(
  <React.StrictMode>
    <JotaiProvider store={jotaiStore}>
      <QueryClientProvider client={queryClient}>
        <MantineProvider defaultColorScheme="dark" theme={theme}>
          <RouterProvider router={router} />
        </MantineProvider>
      </QueryClientProvider>
    </JotaiProvider>
  </React.StrictMode>,
);
