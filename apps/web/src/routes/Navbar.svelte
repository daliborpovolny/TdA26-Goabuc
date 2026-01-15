<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import NavbarLink from './NavbarLink.svelte';

	let show = $state(false);

	function onclick(): void {
		show = false;
	}

	$inspect(show);
</script>

<div class="flex flex-col bg-[#0257A5] md:flex-row md:items-center">
	<div class="grid w-full grid-cols-3 md:flex md:w-fit md:flex-row">
		<button
			type="button"
			class="ml-10 justify-self-start text-4xl md:hidden"
			onclick={() => {
				show = !show;
			}}
			>☰
		</button>

		<img
			src="/resources/Think-different-Academy_LOGO_bily.svg"
			alt="Logo TdA"
			width="50"
			class="m-[calc(100px/3)] justify-self-center"
		/>
	</div>

	<div class={`flex flex-col items-center md:flex-row ${show ? 'flex' : 'hidden'} md:flex`}>
		<NavbarLink name="Home" url="/" {onclick} />
		<NavbarLink name="Courses" url="/courses" {onclick} />

		{#if auth.user}
			<NavbarLink name="Dashboard" url="/dashboard" {onclick} />
			<NavbarLink name="Profile" url="/profile" {onclick} />

			<button onclick={() => auth.logout()}>Logout </button>
		{:else}
			<NavbarLink name="Login" url="/login" {onclick} />
			<NavbarLink name="Register" url="/register" {onclick} />
		{/if}
	</div>
</div>
