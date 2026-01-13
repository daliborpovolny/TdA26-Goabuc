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
	class="mx-auto mt-10 max-w-md space-y-4 rounded-xl bg-[#0070BB] p-6 shadow-lg"
>
	<h2 class="text-center text-5xl font-bold">Login</h2>
	<div>
		<label for="email" class="mb-1 block text-xl font-semibold">Email</label>
		<input
			id="email"
			type="email"
			bind:value={email}
			required
			class="outline-offset-none w-full rounded-md border-[4px] border-[#1A1A1A] p-2 font-semibold focus:outline-2 focus:outline-[#91F5AD] focus:bg-white"
		/>
	</div>
	<div>
		<label for="password" class="mb-1 block text-xl font-semibold">Password</label>
		<input
			id="password"
			type="password"
			bind:value={password}
			required
			class="w-full rounded-md border-[4px] border-[#1A1A1A] p-2 focus:outline-2 focus:outline-[#91F5AD] focus:bg-white"
		/>
	</div>
	<button
		type="submit"
		class="colors w-full cursor-pointer rounded-md bg-[#91F5AD] px-4 py-2 text-2xl font-semibold text-[#1A1A1A] transition transition-colors hover:bg-[#6DD4B1]"
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
