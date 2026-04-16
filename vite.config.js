import { defineConfig } from 'vite';

export default defineConfig({
    build: {
        outDir: 'static/dist',
        emptyOutDir: true,
        manifest: true,
        rollupOptions: {
            input: {
                main: '/frontend/main.js',
                fingerprint: '/frontend/fingerprint.js'
            }
        }
    }
});