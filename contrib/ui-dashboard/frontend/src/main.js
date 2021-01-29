import App from './App.svelte';

const app = new App({
	target: document.body,
	props: {
        channels: [
            "nixos-unstable",
            "nixos-20.09",
        ],
        evaluations: [
            "nixos-21.03pre257339.83cbad92d73",
            "nixos-21.03pre257780.e9158eca70a",
        ],
    }
});

export default app;
