// tailwind.config.ts

import { nextui } from "@nextui-org/react";

/** @type {import('tailwindcss').Config} */
export default {
    content: [
        // ... your content paths
        "./node_modules/@nextui-org/theme/dist/**/*.{js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {},
    },
    darkMode: "class",
    plugins: [nextui()],
};
