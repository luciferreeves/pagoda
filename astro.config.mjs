// @ts-check
import { defineConfig } from 'astro/config';

import tailwindcss from '@tailwindcss/vite';
import neocitiesIntegration from './integrations/neocities.mjs';

// https://astro.build/config
export default defineConfig({
  integrations: [neocitiesIntegration()],
  vite: {
    plugins: [tailwindcss()]
  }
});