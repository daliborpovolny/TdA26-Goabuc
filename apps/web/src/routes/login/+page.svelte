<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/auth.svelte';
	import { fade } from 'svelte/transition';

	let email = $state('');
	let password = $state('');
	let loginPromise: Promise<any> | null = $state(null);

	function handleLogin(e: Event) {
		e.preventDefault();

		loginPromise = fetch('/api/login', {
			method: 'POST',
			headers: { 'Content-type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({ email, password })
		})
			.then(async (res) => {
				const data = await res.json();

				if (!res.ok) throw new Error(data.message || 'Login failed');
				return data;
			})
			.then((data) => {
				auth.check();
				setTimeout(() => goto('/dashboard'), 1500);
				return data;
			});

		// loginPromise = null
	}
</script>

<svelte:head>
	<title>Login | TdA</title>
</svelte:head>

<div class="flex min-h-[80vh] items-center justify-center p-4">
	<div class="relative w-full max-w-md">
		<div class="absolute inset-0 translate-x-3 translate-y-3 rounded-2xl bg-s-black"></div>

		<form
			onsubmit={handleLogin}
			class="relative flex flex-col space-y-6 rounded-2xl border-4 border-s-black bg-p-blue p-8 shadow-none"
		>
			<div class="text-center">
				<h2 class="text-5xl font-black tracking-tighter text-white uppercase">Login</h2>
				<div class="mx-auto mt-2 h-1.5 w-16 bg-p-green"></div>
			</div>

			<div class="space-y-4">
				<div>
					<label
						for="email"
						class="mb-2 block text-lg font-bold tracking-wide text-white uppercase"
					>
						Username
					</label>
					<input
						id="email"
						type="text"
						bind:value={email}
						required
						placeholder="you@example.com"
						class="w-full rounded-xl border-4 border-s-black bg-white p-3 font-bold text-s-black placeholder:text-gray-400 focus:ring-4 focus:ring-p-green focus:outline-none"
					/>
				</div>

				<div>
					<label
						for="password"
						class="mb-2 block text-lg font-bold tracking-wide text-white uppercase"
					>
						Password
					</label>
					<input
						id="password"
						type="password"
						bind:value={password}
						required
						placeholder="••••••••"
						class="w-full rounded-xl border-4 border-s-black bg-white p-3 font-bold text-s-black placeholder:text-gray-400 focus:ring-4 focus:ring-p-green focus:outline-none"
					/>
				</div>
			</div>

			<button
				type="submit"
				class="group relative mt-4 cursor-pointer overflow-hidden rounded-xl border-4 border-s-black bg-p-green py-4 text-2xl font-black tracking-widest text-s-black uppercase transition-all hover:translate-x-0.5 hover:translate-y-0.5 active:translate-x-1 active:translate-y-1 disabled:opacity-50"
			>
				<span>Enter Academy</span>
			</button>

			<p class="text-center font-bold text-white">
				New here?
				<a
					href="/register"
					class="text-p-green underline decoration-4 underline-offset-4 hover:text-white"
				>
					Create Account
				</a>
			</p>

			{#if loginPromise}
				<div class="mt-4 text-center font-black tracking-tight uppercase" transition:fade>
					{#await loginPromise}
						<p class="text-white">Checking credentials...</p>
					{:then data}
						<p class="text-p-green">Welcome back, {data.first_name}!</p>
					{:catch error}
						<p class="rounded-lg border-2 border-s-black bg-red-500 p-2 text-white">
							{error.message}
						</p>
					{/await}
				</div>
			{/if}
		</form>
	</div>
</div>
