import { defineConfig, type PluginOption } from 'vite';
import react from '@vitejs/plugin-react-swc';
import checker from 'vite-plugin-checker';
import tsconfigPaths from 'vite-tsconfig-paths';
import tailwindcss from '@tailwindcss/vite';
import postCssNested from 'postcss-nested';
import postCssSimpleVars from 'postcss-simple-vars';
import postCssMantinePreset from 'postcss-preset-mantine';
import { visualizer } from 'rollup-plugin-visualizer';

// https://vitejs.dev/config/
export default defineConfig({
  root: 'src',
  plugins: [
    tsconfigPaths(),
    tailwindcss(),
    react(),
    checker({
      typescript: true,
      stylelint: {
        lintCommand: 'stylelint **/*.css',
      },
      eslint: {
        lintCommand: 'eslint "./**/*.{ts,tsx}"',
      },
      overlay: {
        initialIsOpen: false,
      },
    }),
    visualizer() as PluginOption,

  ],
  css: {
    postcss: {
      plugins: [
        postCssNested(),
        postCssMantinePreset(),
        postCssSimpleVars({
          variables: {
            'mantine-breakpoint-xs': '36em',
            'mantine-breakpoint-sm': '48em',
            'mantine-breakpoint-md': '62em',
            'mantine-breakpoint-lg': '75em',
            'mantine-breakpoint-xl': '88em',
          },
        }),
      ],
    },
  },
  envPrefix: 'MMP_',
  build: {
    outDir: './public/',
    sourcemap: true,
    emptyOutDir: true,
    assetsInlineLimit: 0,
  },
  server: {
    open: true,
    cors: true,
  },
});
