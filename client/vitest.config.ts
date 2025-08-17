import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    include: [
      "src/**/*.{test,spec}.{js,ts}",
      "tests/**/*.ts"
    ],
    exclude: [
      'node_modules',
      'dist',
      'e2e/**',
      'playwright/**',
      '**/routes/**'
    ],
  },
})
