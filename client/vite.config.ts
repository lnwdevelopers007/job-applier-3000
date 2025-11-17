import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, process.cwd(), '');

	return {
		plugins: [tailwindcss(), sveltekit()],
		optimizeDeps: {
			force: true,
			include: ['lucide-svelte']
		},
		server: {
			host: '0.0.0.0',
			port: 5173,
			strictPort: true,
			watch: {
				usePolling: true
			},
			hmr: {
				host: 'localhost'
			},
			proxy: {
				'/health': { target: env.VITE_BACKEND, changeOrigin: true, secure: false },
				'/jobs': { target: env.VITE_BACKEND, changeOrigin: true, secure: false },
				'/apply': { target: env.VITE_BACKEND, changeOrigin: true, secure: false },
				'/users': { target: env.VITE_BACKEND, changeOrigin: true, secure: false },
				'/files': { target: env.VITE_BACKEND, changeOrigin: true, secure: false },
				'/notes': { target: env.VITE_BACKEND, changeOrigin: true, secure: false },
			}
		},
		test: {
			expect: { requireAssertions: true },
			projects: [
				{
					extends: './vite.config.ts',
					test: {
						name: 'client',
						environment: 'browser',
						browser: {
							enabled: true,
							provider: 'playwright',
							instances: [{ browser: 'chromium' }]
						},
						include: ['src/**/*.svelte.{test,spec}.{js,ts}'],
						exclude: ['src/lib/server/**'],
						setupFiles: ['./vitest-setup-client.ts']
					}
				},
				{
					extends: './vite.config.ts',
					test: {
						name: 'server',
						environment: 'node',
						include: ['src/**/*.{test,spec}.{js,ts}'],
						exclude: ['src/**/*.svelte.{test,spec}.{js,ts}']
					}
				}
			]
		}
	};
});
