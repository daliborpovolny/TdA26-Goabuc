<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/auth.svelte';
	import { fade } from 'svelte/transition';
	import UniButton from '../UniButton.svelte';

	let first_name = $state('');
	let last_name = $state('');
	let email = $state('');
	let password = $state('');

	let acceptedTerms = $state(false);

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
			.then(async (data) => {
				await auth.check();
				setTimeout(() => goto('/courses'), 500);
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
			class="relative flex flex-col space-y-6 border-4 border-s-black bg-p-blue p-6 shadow-none md:p-10"
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

				<div class="flex items-start space-x-3 py-2">
					<div class="relative flex items-center">
						<input
							id="terms"
							type="checkbox"
							bind:checked={acceptedTerms}
							required
							class="h-6 w-6 cursor-pointer appearance-none rounded-md border-4 border-s-black bg-white transition-all checked:bg-p-green focus:ring-2 focus:ring-p-green focus:outline-none"
						/>
						{#if acceptedTerms}
							<span
								class="pointer-events-none absolute inset-0 flex items-center justify-center text-s-black"
							>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									class="h-4 w-4"
									viewBox="0 0 20 20"
									fill="currentColor"
								>
									<path
										fill-rule="evenodd"
										d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
										clip-rule="evenodd"
									/>
								</svg>
							</span>
						{/if}
					</div>
					<label for="terms" class="text-sm font-bold text-white select-none">
						I AGREE TO THE
						<a href="/about" class="text-p-green underline hover:text-white">TERMS OF SERVICE</a>
						AND
						<a href="/about" class="text-p-green underline hover:text-white">PRIVACY POLICY</a>.
					</label>
				</div>
			</div>

			<UniButton
				type="submit"
				disabled={!!registerPromise || !acceptedTerms}
				uppercase
				more_style="disabled:opacity-50"
			>
				{#if !registerPromise}
					<span>Start Learning</span>
				{:else}
					<span class="inline-block animate-pulse">Creating Account...</span>
				{/if}
			</UniButton>

			<p class="text-center font-bold text-white">
				Already a member?
				<a href="/login" class="text-p-green underline hover:text-white"> Log In </a>
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
