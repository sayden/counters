// @ts-check
import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";
import react from "@astrojs/react";

import tailwindcss from "@tailwindcss/vite";

// https://astro.build/config
export default defineConfig({
  base: "/documentation/dist/",
  integrations: [
    react(),
    starlight({
      title: "Counters",
      components: {
        Header: "./src/overrides/Header.astro",
      },
      customCss: ["./src/styles/global.css", "./src/styles/overrides.css"],
      social: [
        {
          icon: "github",
          label: "GitHub",
          href: "https://github.com/sayden/counters",
        },
      ],
      sidebar: [
        {
          label: "Getting Started",
          autogenerate: { directory: "introduction" },
        },
        {
          label: "Counters",
          autogenerate: { directory: "counters" },
        },
        {
          label: "Cards",
          autogenerate: { directory: "cards" },
        },
        {
          label: "Prototypes",
          autogenerate: { directory: "prototypes" },
        },
        {
          label: "Settings",
          autogenerate: { directory: "settings" },
        },
        {
          label: "CLI",
          autogenerate: { directory: "cli" },
        },
      ],
    }),
  ],

  vite: {
    plugins: [tailwindcss()],
  },
});
