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


<form on:submit={login} class="mx-auto mt-10 max-w-md space-y-4 rounded-xl bg-[#0070BB] p-6 shadow-lg">
	<h2 class="text-center text-5xl font-bold">Login</h2>
	<div>
		<label for="email" class="mb-1 block text-xl font-semibold">Email</label>
		<input
			id="email"
			type="email"
			bind:value={email}
			required
			class="w-full rounded-md border-[4px] border-[#1A1A1A] px-4 py-2 focus:border-[#91F5AD] focus:outline-[#91F5AD] focus:outline-2 outline-offset-none font-semibold"
		/>
	</div>
	<div>
		<label for="password" class="mb-1 block text-xl font-semibold">Password</label>
		<input
			id="password"
			type="password"
			bind:value={password}
			required
			class="w-full rounded-md border-[4px] border-[#1A1A1A] px-4 py-2 focus:border-[#91F5AD] focus:outline-[#91F5AD] focus:outline-2"
		/>
	</div>
	<button
		type="submit"
		class="w-full rounded-md bg-[#91F5AD] px-4 py-2 font-semibold text-[#1A1A1A] transition-colors hover:bg-[#6DD4B1] text-2xl cursor-pointer transition colors font-semibold"
	>
		Login
	</button>
	<p class="text-center text-xl text-[#1A1A1A] font-semibold">
		Don't have an account? <a href="/register" class="hover:underline text-[#91F5AD]"
			>Register</a
		>
	</p>

	{#if loginPromise}
		{#await loginPromise}
			<p class="text-xl font-semibold">Logging in...</p>
		{:then data}
			<p class="text-[#91F5AD] text-xl font-semibold">Success! Welcome {data.first_name}</p>
		{:catch error}
			<p class="text-red-500 capitalize text-xl font-semibold">{error.message}</p>
		{/await}
	{/if}
</form>
