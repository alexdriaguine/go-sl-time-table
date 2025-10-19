import { defineConfig } from "vite";
import solidPlugin from "vite-plugin-solid";

export default defineConfig({
	plugins: [solidPlugin()],
	build: {
		outDir: "internal/static",
		emptyOutDir: true,
		rollupOptions: {
			input: "src/index.html",
		},
	},
	server: {
		port: 3014,
		proxy: {
			"/api": {
				target: "http://localhost:3000",
				changeOrigin: true,
			},
		},
	},
});
