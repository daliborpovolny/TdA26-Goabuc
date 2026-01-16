<script lang="ts">
	import { fade } from 'svelte/transition';

	let dataPromise = fetch('/api/me').then((r) => {
		if (!r.ok) throw new Error('Failed to fetch profile');
		return r.json();
	});
</script>

<svelte:head>
	<title>My Profile | TdA</title>
</svelte:head>

<div class="flex min-h-[80vh] items-center justify-center p-6">
	{#await dataPromise}
		<p class="animate-pulse text-2xl font-black uppercase">Scanning Bio-Data...</p>
	{:then user}
		<div class="relative w-full max-w-md" in:fade>
			<div class="absolute inset-0 translate-x-4 translate-y-4 rounded-3xl bg-s-black"></div>

			<div class="relative overflow-hidden rounded-3xl border-4 border-s-black bg-white">
				<div class="bg-p-blue p-8 text-center text-white">
					<div
						class="mx-auto mb-4 flex h-24 w-24 items-center justify-center rounded-full border-4 border-s-black bg-p-green text-4xl font-black text-s-black shadow-md"
					>
						{user.firstName[0]}{user.lastName[0]}
					</div>
					<h1 class="text-4xl font-black tracking-tight uppercase">Profile</h1>
					<p class="text-xs font-bold tracking-widest uppercase opacity-70">
						Academy ID: {user.uuid || 'Active'}
					</p>
				</div>

				<div class="space-y-6 p-8">
					<div class="space-y-1">
						<span class="text-[10px] font-black tracking-widest text-gray-400 uppercase"
							>Full Name</span
						>
						<p class="text-2xl font-bold text-s-black">{user.firstName} {user.lastName}</p>
					</div>

					<div class="space-y-1 border-t-2 border-dashed border-gray-100 pt-4">
						<span class="text-[10px] font-black tracking-widest text-gray-400 uppercase"
							>Email Address</span
						>
						<p class="text-xl font-bold text-p-blue">{user.email}</p>
					</div>

					<button
						class="mt-4 w-full cursor-pointer rounded-xl border-4 border-s-black bg-s-black py-4 text-sm font-black text-white uppercase transition-all hover:bg-p-green hover:text-s-black active:translate-y-1 active:shadow-none"
					>
						Edit Account Settings
					</button>
				</div>
			</div>
		</div>
	{:catch error}
		<div
			class="rounded-xl border-4 border-s-black bg-red-500 p-8 text-white shadow-[4px_4px_0px_0px_rgba(26,26,26,1)]"
		>
			<h2 class="text-3xl font-black uppercase">Unauthorized</h2>
			<p class="mt-2 font-bold">{error.message}</p>
			<a href="/login" class="mt-4 inline-block font-black underline">Back to Login</a>
		</div>
	{/await}
</div>
