<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import NavbarLink from './NavbarLink.svelte';
	import tda_logo from '../../static/resources/Think-different-Academy_LOGO_bily.svg';
</script>

<nav class="flex gap-4 bg-[#0257A5] p-4 flex items-center text-4xl">
	<img src={tda_logo} alt="Logo TdA" width="100" class="m-[calc(100px/3)]"/>
	
	<NavbarLink name="Home" url="/" />
	<NavbarLink name="Courses" url="/courses" />

	{#if auth.user}
		<NavbarLink name="Dashboard" url="/dashboard" />
		<NavbarLink name="Profile" url="/profile" />

		<button class="rounded-md px-3 py-2 text-[#FFF] hover:bg-[#91F5AD] hover:text-[#000] cursor-pointer h-[66px] flex items-center" onclick={() => auth.logout()}
			>Logout
		</button>
	{:else}
		<NavbarLink name="Login" url="/login" />
		<NavbarLink name="Register" url="/register" />
	{/if}
</nav>
