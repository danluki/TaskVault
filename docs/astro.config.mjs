// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	redirects: {
		"/": "/en",
	},
	integrations: [
		starlight({
			title: 'Syncra',
			logo: {
				src: './src/assets/syncra.png'
			},
			social: {
				github: 'https://github.com/danluki/TaskVault',
			},
			sidebar: [
				{
					label: 'Basics',
					autogenerate: { directory: 'basics' },
				},
				{
					label: 'Guides',
					items: [
						// Each item here is one entry in the navigation menu.
						{ label: 'Kubernetes basics', slug: 'guides/minikube' },
						{ label: 'Getting started with docker', slug: 'guides/docker' },
					],
				},
				{
					label: 'Intro',
					autogenerate: { directory: 'intro' },
				},
				{
					label: 'CLI',
					autogenerate: { directory: 'cli' },
				},
				{
					label: 'Internals',
					autogenerate: { directory: 'internals' },
				},
				{
					label: 'Commercial',
					autogenerate: { directory: 'commercial' },
				},
			],
			defaultLocale: 'en',
			locales: {
				en: {
					label: 'English',
				},
				ru: {
					label: 'Russian'
				},
			}
		}),
	],
});
