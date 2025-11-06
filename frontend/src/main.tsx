import React from "react";
import ReactDOM from "react-dom/client";
import { router } from "./navigation.tsx";
import { RouterProvider } from "react-router";
import { createTheme, MantineProvider, rem } from "@mantine/core";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";


const queryClient = new QueryClient()
const theme = createTheme({
    autoContrast: true,
    cursorType: 'pointer',
    headings: {
        fontFamily: 'Roboto, sans-serif',
    },
});
ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
    <React.StrictMode>
        <QueryClientProvider client={queryClient}>

            <MantineProvider defaultColorScheme="dark" theme={theme}>
                <RouterProvider router={router} />
            </MantineProvider>
        </QueryClientProvider>
    </React.StrictMode>
);