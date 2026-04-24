// 2025 Nicholas Triska. All rights reserved.
// The SSTPA Tools software and all associated modules, binaries, and source code
// are proprietary intellectual property of Nicholas Triska. Unauthorized
// reproduction, modification, or distribution is strictly prohibited. Licensed
// copies may be used under specific contractual terms provided by the author.
import { resolve } from "node:path"
import { fileURLToPath } from "node:url"

import react from "@vitejs/plugin-react"
import { defineConfig } from "vite"

const root = fileURLToPath(new URL("../../", import.meta.url))

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@sstpa/addon-sdk": resolve(root, "packages/addon-sdk/src/index.ts"),
      "@sstpa/domain": resolve(root, "packages/domain/src/index.ts"),
      "@sstpa/ui": resolve(root, "packages/ui/src/index.ts"),
      "@sstpa/navigator-tool": resolve(root, "addons/navigator/src/index.ts"),
      "@sstpa/requirements-tool": resolve(root, "addons/requirements/src/index.ts"),
      "@sstpa/message-center-tool": resolve(root, "addons/message-center/src/index.ts"),
    },
  },
})
