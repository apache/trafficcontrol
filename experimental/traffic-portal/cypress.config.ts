import { defineConfig } from "cypress";

export default defineConfig({
	component: {
		devServer: {
			bundler: "webpack",
			framework: "angular",
		},
		specPattern: "**/*.cy.ts"
	},
	e2e: {
		baseUrl: "http://localhost:4200"
	},

});
