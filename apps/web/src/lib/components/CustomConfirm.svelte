<script lang="ts">
	import { modal } from '$lib/modal.svelte';
	import { fade, scale } from 'svelte/transition';
</script>

{#if modal.show}
	<div
		transition:fade={{ duration: 150 }}
		class="fixed inset-0 z-100 flex items-center justify-center bg-s-black/60 p-4 backdrop-blur-sm"
	>
		<div
			transition:scale={{ start: 0.9, duration: 200 }}
			class="relative w-full max-w-sm rounded-3xl border-4 border-s-black bg-white p-8 shadow-[8px_8px_0px_0px_rgba(26,26,26,1)]"
		>
			<div class="mb-6 text-center">
				<div
					class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full border-4 border-s-black bg-red-500 text-3xl text-white"
				>
					⚠️
				</div>
				<h3 class="text-2xl font-black tracking-tighter text-s-black uppercase">Are you sure?</h3>
				<p class="mt-4 leading-tight font-bold text-gray-500">
					{modal.message}
				</p>
			</div>

			<div class="flex flex-col gap-3">
				<button
					onclick={() => modal.close(true)}
					class="w-full rounded-xl border-4 border-s-black bg-p-green py-3 text-lg font-black tracking-widest text-s-black uppercase transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none active:bg-s-1"
				>
					Confirm
				</button>

				<button
					onclick={() => modal.close(false)}
					class="w-full rounded-xl border-4 border-s-black bg-white py-3 text-lg font-black tracking-widest text-s-black uppercase transition-all hover:bg-gray-100"
				>
					Cancel
				</button>
			</div>
		</div>
	</div>
{/if}
