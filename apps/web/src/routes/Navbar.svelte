<script lang="ts">
	import { auth } from '$lib/auth.svelte';
	import { page } from '$app/state';
	import NavbarLink from './NavbarLink.svelte';
	import { slide } from 'svelte/transition';

	let show = $state(false);

	function closeMenu(): void {
		show = false;
	}
</script>

<nav class="border-b-4 border-s-black bg-s-blue-darker text-white shadow-md">
	<div class="mx-auto flex max-w-7xl flex-col md:flex-row md:items-center">
		<div class="grid w-full grid-cols-3 items-center px-4 py-2 md:flex md:w-fit md:px-6">
			<button
				type="button"
				class="justify-self-start text-4xl transition-transform active:scale-90 md:hidden"
				onclick={() => (show = !show)}
			>
				{show ? '✕' : '☰'}
			</button>

			<a href="/" class="justify-self-center py-2 md:pr-8" onclick={closeMenu}>
				<img
					src="/resources/Think-different-Academy_LOGO_bily.svg"
					alt="Logo TdA"
					width="50"
					class="transition-transform hover:rotate-3"
				/>
			</a>

			<div class="md:hidden"></div>
		</div>

		<div
			class={`flex flex-col items-center md:flex md:flex-1 md:flex-row md:justify-start ${show ? 'flex' : 'hidden'}`}
			transition:slide={{ duration: 300 }}
		>
			<NavbarLink name="Home" url="/" onclick={closeMenu} />
			<NavbarLink name="Courses" url="/courses" onclick={closeMenu} />

			{#if auth.user}
				<NavbarLink name="Dashboard" url="/dashboard" onclick={closeMenu} />
				<NavbarLink name="Profile" url="/profile" onclick={closeMenu} />

				<button
					onclick={() => {
						auth.logout();
						closeMenu();
					}}
					class="w-full px-6 py-4 text-left font-black tracking-widest text-red-300 uppercase hover:bg-red-500 hover:text-white md:ml-auto md:w-auto md:py-2 md:text-sm"
				>
					Logout
				</button>
			{:else}
				<div class="flex w-full flex-col md:ml-auto md:w-auto md:flex-row">
					<NavbarLink name="Login" url="/login" onclick={closeMenu} />
					<NavbarLink name="Register" url="/register" onclick={closeMenu} />
				</div>
			{/if}
		</div>
	</div>
</nav>
