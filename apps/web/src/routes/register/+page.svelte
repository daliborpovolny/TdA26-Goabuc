<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/auth.svelte';
	import { fade } from 'svelte/transition';

	let first_name = $state('');
	let last_name = $state('');
	let email = $state('');
	let password = $state('');

	let registerPromise: Promise<any> | null = $state(null);

	function handleRegister(e: Event) {
		e.preventDefault();

		registerPromise = fetch('/api/register', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({ first_name, last_name, email, password })
		})
			.then(async (res) => {
				const data = await res.json();
				if (!res.ok) throw new Error(data.message || 'Registration failed');
				return data;
			})
			.then((data) => {
				auth.check();
				setTimeout(() => goto('/dashboard'), 1500);
				return data;
			});
	}
</script>

<svelte:head>
	<title>Join the Academy | TdA</title>
</svelte:head>

<div class="flex min-h-screen items-center justify-center p-4 py-12">
	<div class="relative w-full max-w-2xl">
		<div class="absolute inset-0 translate-x-3 translate-y-3 rounded-2xl bg-s-black"></div>

		<form
			onsubmit={handleRegister}
			class="relative flex flex-col space-y-6 rounded-2xl border-4 border-s-black bg-p-blue p-6 shadow-none md:p-10"
		>
			<div class="text-center">
				<h2 class="text-5xl font-black tracking-tighter text-white uppercase md:text-6xl">
					Register
				</h2>
				<div class="mx-auto mt-2 h-2 w-20 bg-p-green"></div>
				<p class="mt-4 text-sm font-bold tracking-widest text-white/80 uppercase">
					Join Think Different Academy
				</p>
			</div>

			<div class="space-y-4">
				<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
					<div>
						<label for="first_name" class="mb-2 block text-lg font-bold text-white uppercase"
							>First Name</label
						>
						<input
							id="first_name"
							type="text"
							bind:value={first_name}
							required
							placeholder="Krtkus"
							class="w-full rounded-xl border-4 border-s-black bg-white p-3 font-bold text-s-black focus:ring-4 focus:ring-p-green focus:outline-none"
						/>
					</div>
					<div>
						<label for="last_name" class="mb-2 block text-lg font-bold text-white uppercase"
							>Last Name</label
						>
						<input
							id="last_name"
							type="text"
							bind:value={last_name}
							required
							placeholder="Veliky"
							class="w-full rounded-xl border-4 border-s-black bg-white p-3 font-bold text-s-black focus:ring-4 focus:ring-p-green focus:outline-none"
						/>
					</div>
				</div>

				<div>
					<label for="email" class="mb-2 block text-lg font-bold text-white uppercase"
						>Email Address</label
					>
					<input
						id="email"
						type="email"
						bind:value={email}
						required
						placeholder="krtkus@priklad.cz"
						class="w-full rounded-xl border-4 border-s-black bg-white p-3 font-bold text-s-black focus:ring-4 focus:ring-p-green focus:outline-none"
					/>
				</div>

				<div>
					<label for="password" class="mb-2 block text-lg font-bold text-white uppercase"
						>Password</label
					>
					<input
						id="password"
						type="password"
						bind:value={password}
						required
						placeholder="••••••••"
						class="w-full rounded-xl border-4 border-s-black bg-white p-3 font-bold text-s-black focus:ring-4 focus:ring-p-green focus:outline-none"
					/>
				</div>
			</div>

			<button
				type="submit"
				disabled={!!registerPromise}
				class="group relative mt-4 overflow-hidden rounded-xl border-4 border-s-black bg-p-green py-4 text-2xl font-black tracking-widest text-s-black uppercase transition-all hover:translate-x-0.5 hover:translate-y-0.5 active:translate-x-1 active:translate-y-1 disabled:opacity-50"
			>
				{#if !registerPromise}
					<span>Start Learning</span>
				{:else}
					<span class="inline-block animate-pulse">Creating Account...</span>
				{/if}
			</button>

			<p class="text-center font-bold text-white">
				Already a member?
				<a
					href="/login"
					class="text-p-green underline decoration-4 underline-offset-4 hover:text-white"
				>
					Log In
				</a>
			</p>

			{#if registerPromise}
				<div class="mt-4 text-center font-black tracking-tight uppercase" transition:fade>
					{#await registerPromise}
						<p class="text-white">Processing your application...</p>
					{:then data}
						<p class="text-p-green">Account Created! Welcome, {data.first_name}!</p>
					{:catch error}
						<div class="rounded-lg border-2 border-s-black bg-red-500 p-3 text-white">
							{error.message}
						</div>
					{/await}
				</div>
			{/if}
		</form>
	</div>
</div>
