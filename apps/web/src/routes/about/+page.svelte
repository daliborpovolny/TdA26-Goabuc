<script lang="ts">
	import TermsOfService from './TermsOfService.svelte';
	import PrivacyPolicy from './PrivacyPolicy.svelte';
	import { fade } from 'svelte/transition';
	import PrimaryButton from '$lib/components/PrimaryButton.svelte';

	let activeTab = $state<'terms' | 'privacy'>('terms');
</script>

<svelte:head>
	<title>About | TdA</title>
</svelte:head>

<div class="min-h-screen bg-p-green/10 p-4 py-16 md:p-12">
	<div class="mx-auto max-w-4xl">
		<header class="mb-12 text-center">
			<h1 class="text-6xl font-black tracking-tighter text-s-black uppercase md:text-8xl">
				The Rules
			</h1>
		</header>

		<div class="mb-12 flex flex-col gap-4 sm:flex-row sm:justify-center">
			<!-- <button onclick={() => (activeTab = 'terms')} class="group relative px-8 py-4 transition-all">
				<div class="absolute inset-0 translate-x-1 translate-y-1 rounded-xl bg-s-black"></div>
				<div
					class="relative border-4 border-s-black px-8 py-4 text-xl font-black uppercase transition-transform
                    {activeTab === 'terms'
						? '-translate-x-1 -translate-y-1 bg-p-blue text-white'
						: 'bg-white text-s-black hover:-translate-y-0.5'}"
				>
					Terms of Service
				</div>
			</button> -->

			<PrimaryButton onclick={() => (activeTab = 'terms')} class={activeTab === 'terms' ? 'shadow-none translate-x-1 translate-y-1' : ''}>
				Terms of Service
			</PrimaryButton>

			<!-- <button
				onclick={() => (activeTab = 'privacy')}
				class="group relative px-8 py-4 transition-all"
			>
				<div class="absolute inset-0 translate-x-1 translate-y-1 rounded-xl bg-s-black"></div>
				<div
					class="relative border-4 border-s-black px-8 py-4 text-xl font-black uppercase transition-transform
                    {activeTab === 'privacy'
						? '-translate-x-1 -translate-y-1 bg-p-green text-s-black'
						: 'bg-white text-s-black hover:-translate-y-0.5'}"
				>
					Privacy Policy
				</div>
			</button> -->

			<PrimaryButton onclick={() => (activeTab = 'privacy')} class={activeTab === 'privacy' ? 'shadow-none translate-x-1 translate-y-1' : ''}>
				Privacy Policy
			</PrimaryButton>
		</div>

		<main class="relative">
			{#if activeTab === 'terms'}
				<div in:fade={{ duration: 200 }}>
					<TermsOfService />
				</div>
			{:else}
				<div in:fade={{ duration: 200 }}>
					<PrivacyPolicy />
				</div>
			{/if}
		</main>

		<div class="mt-16 text-center">
			<a
				href="/"
				class="font-black uppercase underline decoration-p-blue decoration-4 underline-offset-4 hover:text-p-blue"
			>
				Return Home
			</a>
		</div>
	</div>
</div>
