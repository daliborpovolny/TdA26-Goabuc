<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="mx-3">
	<nav
		class="mx-3 flex gap-4
            [&>a]:rounded-md [&>a]:px-3
            [&>a]:py-2
            [&>a:hover]:bg-stone-200"
	>
		<a href="/">Home</a>
		<a href="/courses">Courses</a>
		<a href="/dashboard">Dashboard</a>
		<a href="/register">Register</a>
		<a href="/login">Login</a>
	</nav>

	<div class="mt-3">
		{@render children()}
	</div>
</div>
