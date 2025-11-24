import React from 'react';
import ReactDOM from 'react-dom/client';
import { createTheme, MantineProvider } from '@mantine/core';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Provider as JotaiProvider } from 'jotai';
import { RouterProvider, createBrowserRouter } from 'react-router';

import jotaiStore from 'core/stores/jotai.store';
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
    <JotaiProvider store={jotaiStore}>
      <QueryClientProvider client={queryClient}>
        <MantineProvider defaultColorScheme="dark" theme={theme}>
          <RouterProvider router={router} />
        </MantineProvider>
      </QueryClientProvider>
    </JotaiProvider>
  </React.StrictMode>,
);
