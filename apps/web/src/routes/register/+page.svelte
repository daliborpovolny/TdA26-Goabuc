<script lang="ts">
	import { goto } from '$app/navigation';
	import { checkSession } from '$lib/auth.svelte';

	let first_name = '';
	let last_name = '';

	let email = '';
	let password = '';

	let registerPromise: Promise<any> | null = null;

	function register() {
		registerPromise = fetch('/api/register', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({ first_name, last_name, email, password })
		})
			.then(async (res) => {
				if (!res.ok) {
					const err = await res.json();
					throw new Error(err.message || 'Registration failed');
				}
				return res.json();
			})
			.then((data) => {
				checkSession(fetch);

				setTimeout(() => {
					goto('/dashboard');
				}, 1500);

				return data;
			});
	}
</script>

<title>Register</title>

<form
	class="mx-auto mt-10 max-w-md space-y-4 rounded-xl bg-white p-6 shadow-lg"
	on:submit={register}
>
	<h2 class="text-center text-2xl font-bold">Register</h2>

	<div>
		<label for="first_name" class="mb-1 block text-sm font-medium text-gray-700">First name</label>
		<input
			id="first_name"
			type="text"
			bind:value={first_name}
			required
			class="w-full rounded-md border border-gray-300 px-4 py-2 focus:border-stone-500 focus:ring-2 focus:ring-stone-500 focus:outline-none"
		/>
	</div>

	<div>
		<label for="last_name" class="mb-1 block text-sm font-medium text-gray-700">Last name</label>
		<input
			id="last_name"
			type="text"
			bind:value={last_name}
			required
			class="w-full rounded-md border border-gray-300 px-4 py-2 focus:border-stone-500 focus:ring-2 focus:ring-stone-500 focus:outline-none"
		/>
	</div>

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
		Register
	</button>

	<p class="text-center text-sm text-gray-500">
		Already have an account? <a href="/login" class="text-stone-600 hover:underline">Login</a>
	</p>

	{#if registerPromise}
		{#await registerPromise}
			<p>Registering...</p>
		{:then data}
			<p class="text-green-500">Success! Welcome {data.first_name}</p>
		{:catch error}
			<p class="text-red-500 capitalize">{error.message}</p>
		{/await}
	{/if}
</form>
