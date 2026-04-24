import type { Config } from "tailwindcss"

export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}", "../../addons/**/*.{ts,tsx}", "../../packages/**/*.{ts,tsx}"],
  theme: {
    extend: {},
  },
  plugins: [],
} satisfies Config
