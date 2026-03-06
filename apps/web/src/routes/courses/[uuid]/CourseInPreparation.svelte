<script lang="ts">
	import { fade, scale } from 'svelte/transition';
</script>

<div
	in:fade
	class="relative my-4 overflow-hidden rounded-2xl border-4 border-s-black bg-white p-6 shadow-[6px_6px_0px_0px_rgba(245,158,11,1)] md:p-10"
>
	<div class="absolute -top-6 -right-6 rotate-12 text-[8rem] opacity-5 select-none md:text-[12rem]">
		🏗️
	</div>

	<div class="relative z-10 flex flex-col items-center text-center">
		<div
			in:scale={{ duration: 400, delay: 200 }}
			class="mb-4 flex h-16 w-16 items-center justify-center rounded-full border-4 border-s-black bg-amber-400 text-3xl shadow-[3px_3px_0px_0px_rgba(0,0,0,1)] md:h-20 md:w-20 md:text-4xl"
		>
			🛠️
		</div>

		<h2 class="mb-2 text-3xl font-black tracking-tighter text-amber-600 uppercase md:text-5xl">
			Under Construction
		</h2>

		<p class="max-w-sm text-base font-bold text-gray-600 italic md:text-lg">
			We're currently polishing the curriculum. This course will be available for enrollment
			shortly!
		</p>

		<div class="mt-6 flex flex-wrap justify-center gap-4">
			<a
				href="/courses"
				class="rounded-lg border-4 border-s-black bg-s-black px-6 py-2 text-sm font-black text-white transition-transform hover:-translate-y-1 active:translate-y-0 md:text-base"
			>
				BACK TO ACADEMY
			</a>
		</div>
	</div>
</div>
