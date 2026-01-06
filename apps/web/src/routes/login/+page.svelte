<title>Login</title>

<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/auth.svelte';

	let email = '';
	let password = '';

	let loginPromise: Promise<any> | null = null;

	function login() {
		loginPromise = fetch('/api/login', {
			method: 'POST',
			headers: { 'Content-type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({ email, password })
		})
			.then(async (res) => {
				if (!res.ok) {
					const err = await res.json();
					throw new Error(err.message || 'Login failed');
				}
				return res.json();
			})
			.then((data) => {
				auth.check();

				setTimeout(() => {
					goto('/dashboard');
				}, 1500);

				return data;
			});
	}
</script>

<title>Login</title>
<form on:submit={login} class="mx-auto mt-10 max-w-md space-y-4 rounded-xl bg-white p-6 shadow-lg">
	<h2 class="text-center text-2xl font-bold">Login</h2>
	<div>
		<label for="email" class="mb-1 block text-sm font-medium text-gray-700">Email</label>
		<input
			id="email"
			type="email"
			bind:value={email}
			required
			class="w-full rounded-md border border-gray-300 px-4 py-2 focus:border-stone-500 focus:ring-2 focus:ring-stone-500 focus:outline-none"
		/>
	</div>
	<div>
		<label for="password" class="mb-1 block text-sm font-medium text-gray-700">Password</label>
		<input
			id="password"
			type="password"
			bind:value={password}
			required
			class="w-full rounded-md border border-gray-300 px-4 py-2 focus:border-stone-500 focus:ring-2 focus:ring-stone-500 focus:outline-none"
		/>
	</div>
	<button
		type="submit"
		class="w-full rounded-md bg-stone-300 px-4 py-2 font-semibold text-white transition-colors hover:bg-stone-400"
	>
		Login
	</button>
	<p class="text-center text-sm text-gray-500">
		Don't have an account? <a href="/register" class="texthover:underline text-stone-600"
			>Register</a
		>
	</p>

	{#if loginPromise}
		{#await loginPromise}
			<p>Logging in...</p>
		{:then data}
			<p class=": text-green-500">Success! Welcome {data.first_name}</p>
		{:catch error}
			<p class="text-red-500 capitalize">{error.message}</p>
		{/await}
	{/if}
</form>
