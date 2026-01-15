<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/auth.svelte';

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
				auth.check();

				setTimeout(() => {
					goto('/dashboard');
				}, 1500);

				return data;
			});
	}
</script>

<title>Register</title>

<form
	class="flex mx-auto flex-col max-w-xs md:max-w-md mt-10 space-y-4 rounded-xl bg-[#0070BB] p-6 shadow-lg"
	on:submit={register}
>
	<h2 class="text-center text-5xl font-bold">Register</h2>

	<div>
		<label for="first_name" class="mb-1 block text-xl font-semibold text-[#1A1A1A]"
			>First name</label
		>
		<input
			id="first_name"
			type="text"
			bind:value={first_name}
			required
			class="w-full rounded-md border-4 border-[#1A1A1A] p-2 font-semibold bg-white focus:outline-2 focus:outline-[#91F5AD]"
		/>
	</div>

	<div>
		<label for="last_name" class="mb-1 block text-xl font-semibold text-[#1A1A1A]">Last name</label>
		<input
			id="last_name"
			type="text"
			bind:value={last_name}
			required
			class="w-full rounded-md border-4 border-[#1A1A1A] p-2 font-semibold bg-white focus:outline-2 focus:outline-[#91F5AD]"
		/>
	</div>

	<div>
		<label for="email" class="mb-1 block text-xl font-semibold text-[#1A1A1A]">Email</label>
		<input
			id="email"
			type="email"
			bind:value={email}
			required
			class="w-full rounded-md border-4 border-[#1A1A1A] p-2 font-semibold bg-white focus:outline-2 focus:outline-[#91F5AD]"
		/>
	</div>

	<div>
		<label for="password" class="mb-1 block text-xl font-semibold text-[#1A1A1A]">Password</label>
		<input
			id="password"
			type="password"
			bind:value={password}
			required
			class="w-full rounded-md border-4 border-[#1A1A1A] p-2 font-semibold bg-white focus:outline-2 focus:outline-[#91F5AD]"
		/>
	</div>

	<button
		type="submit"
		class="w-full cursor-pointer rounded-md bg-[#91F5AD] px-4 py-2 text-2xl font-semibold text-[#1A1A1A] transition-colors hover:bg-[#6DD4B1]"
	>
		Register
	</button>

	<p class="text-center text-xl font-semibold text-[#1A1A1A]">
		Already have an account? <a href="/login" class="text-[#91F5AD] hover:underline">Log in</a>
	</p>

	{#if registerPromise}
		{#await registerPromise}
			<p class="text-xl font-semibold">Registering...</p>
		{:then data}
			<p class="text-xl font-semibold text-[#91F5AD]">Success! Welcome {data.first_name}</p>
		{:catch error}
			<p class="text-xl font-semibold text-red-500 capitalize">{error.message}</p>
		{/await}
	{/if}
</form>
