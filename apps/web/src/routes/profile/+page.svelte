<script lang="ts">
	import DataLoader from '$lib/components/DataLoader.svelte';
	import type { User } from '$lib/auth.svelte';

	let dataPromise = fetch('/api/me').then((r) => r.json());
</script>

<title>Profile</title>

<DataLoader promise={dataPromise}>
	{#snippet children(user: User)}
		<div class="ml-4 text-3xl">
			<br />
			<h1 class="text-5xl font-bold">This is you!</h1>
			<br />
			<p>First name: {user.firstName}</p>
			<br />
			<p>Last name: {user.lastName}</p>
			<br />
			<p>Email: {user.email}</p>
		</div>
	{/snippet}
</DataLoader>
