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

<form
	on:submit={login}
	class="flex mx-auto flex-col max-w-xs md:max-w-md mt-10 space-y-4 rounded-xl bg-[#0070BB] p-6 shadow-lg"
>
	<h2 class="text-center text-5xl font-bold">Login</h2>
	<div>
		<label for="email" class="mb-1 block text-xl font-semibold">Email</label>
		<input
			id="email"
			type="email"
			bind:value={email}
			required
			class="w-full outline-offset-none rounded-md border-4 border-[#1A1A1A] p-2 font-semibold bg-white focus:outline-2 focus:outline-[#91F5AD]"
		/>
	</div>
	<div>
		<label for="password" class="mb-1 block text-xl font-semibold">Password</label>
		<input
			id="password"
			type="password"
			bind:value={password}
			required
			class="w-full rounded-md border-4 border-[#1A1A1A] p-2 bg-white focus:outline-2 focus:outline-[#91F5AD]"
		/>
	</div>
	<button
		type="submit"
		class="colors cursor-pointer rounded-md bg-[#91F5AD] px-4 py-2 text-xl font-semibold text-[#1A1A1A] transition hover:bg-[#6DD4B1]"
	>
		Login
	</button>
	<p class="text-center text-xl font-semibold text-[#1A1A1A]">
		Don't have an account? <a href="/register" class="text-[#91F5AD] hover:underline">Register</a>
	</p>

	{#if loginPromise}
		{#await loginPromise}
			<p class="text-xl font-semibold">Logging in...</p>
		{:then data}
			<p class="text-xl font-semibold text-[#91F5AD]">Success! Welcome {data.first_name}</p>
		{:catch error}
			<p class="text-xl font-semibold text-red-500 capitalize">{error.message}</p>
		{/await}
	{/if}
</form>
