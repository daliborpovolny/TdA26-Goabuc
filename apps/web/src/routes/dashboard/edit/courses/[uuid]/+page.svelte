<script lang="ts">
	import { page } from '$app/state';

	let course = $state<any | null>(null);
	let loadingCourse = $state(true);

	let courseName = $state('');
	let courseDescription = $state('');

	let updateStatus = $state<'idle' | 'waiting' | 'success' | 'error'>('idle');

	let createNewMaterial = $state<'file' | 'url' | 'idle'>('idle');

	async function loadCourseDetail() {
		loadingCourse = true;

		try {
			const res = await fetch('/api/courses/' + page.params.uuid, {
				method: 'GET',
				headers: { 'Content-type': 'application/json' }
			});

			if (res.status === 404) throw new Error('Unknown course');
			if (!res.ok) throw new Error('Failed to get course info');

			const data = await res.json();
			course = data;

			courseName = data.name;
			courseDescription = data.description;
		} catch (err: any) {
			course = null;
		} finally {
			loadingCourse = false;
		}
	}

	loadCourseDetail();

	// --- Update course ---
	async function updateCourse(event: SubmitEvent) {
		event.preventDefault();
		updateStatus = 'waiting';

		try {
			const res = await fetch('/api/courses/' + page.params.uuid, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					name: courseName,
					description: courseDescription
				})
			});

			if (!res.ok) {
				const err = await res.json();
				throw new Error(err.message || 'Failed to update course');
			}

			const updated = await res.json();
			course = 'updated';
			updateStatus = 'success';

			setTimeout(() => (updateStatus = 'idle'), 3000);
		} catch (err: any) {
			console.error(err);
			updateStatus = 'error';
			setTimeout(() => (updateStatus = 'idle'), 4000);
		}
	}
</script>

<div class="mx-auto max-w-2xl space-y-8 p-6">
	{#if loadingCourse}
		<p class="text-gray-600">Loading course detail...</p>
	{:else if course != null}
		<div class="rounded-lg border border-stone-300 bg-stone-50 p-6">
			<h2 class="mb-4 text-xl font-semibold text-gray-800">Edit Course</h2>

			<form onsubmit={updateCourse} class="space-y-4">
				<div>
					<label class="text-sm font-medium text-gray-700" for="course_name">Name</label><br />
					<input
						type="text"
						name="course_name"
						bind:value={courseName}
						class="mt-1 w-full rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
					/>
				</div>

				<div>
					<label class="text-sm font-medium text-gray-700" for="course_description"
						>Description</label
					><br />
					<textarea
						name="course_description"
						bind:value={courseDescription}
						cols="60"
						rows="3"
						class="mt-1 w-full rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
					></textarea>
				</div>

				<div class="pt-2">
					<button
						type="submit"
						class="rounded-md border border-stone-400 bg-stone-200 px-4 py-2 text-gray-800 hover:bg-stone-300"
					>
						Update
					</button>

					<!-- Status message -->
					{#if updateStatus === 'waiting'}
						<span class="ml-3 text-gray-600">Updating...</span>
					{:else if updateStatus === 'success'}
						<span class="ml-3 text-stone-600">Updated!</span>
					{:else if updateStatus === 'error'}
						<span class="ml-3 text-red-600">Failed to update</span>
					{/if}
				</div>
			</form>
		</div>

		<!-- Add Materials Section -->
		<div>
			<h1 class="mb-3 text-2xl font-semibold text-gray-800">Materials</h1>
			<div class="flex gap-4">
				<button
					onclick={() => (createNewMaterial = createNewMaterial == 'url' ? 'file' : 'url')}
					class="rounded-md border border-stone-400 bg-stone-100 px-4 py-2 text-gray-800 hover:bg-stone-200"
				>
					New File Material
				</button>
				<button
					onclick={() => (createNewMaterial = createNewMaterial == 'file' ? 'url' : 'file')}
					class="rounded-md border border-stone-400 bg-stone-100 px-4 py-2 text-gray-800 hover:bg-stone-200"
				>
					New Link Material
				</button>
			</div>

			<br />

			{#if createNewMaterial == 'file'}
				<div class="space-y-4 rounded-lg border border-stone-300 bg-stone-50 p-4">
					<div class="flex flex-col">
						<label class="text-sm font-medium text-gray-700" for="file_material_name">Name</label>
						<input
							type="text"
							name="file_material_name"
							class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
						/>
					</div>

					<div class="flex flex-col">
						<label class="text-sm font-medium text-gray-700" for="file_material_description"
							>Description</label
						>
						<textarea
							name="file_material_description"
							class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
						></textarea>
					</div>

					<div class="flex flex-col">
						<label class="text-sm font-medium text-gray-700" for="file_file">File</label>
						<input type="file" name="file_file" class="mt-1 text-gray-900" />
					</div>
				</div>
			{:else if createNewMaterial == 'url'}
				<div class="space-y-4 rounded-lg border border-stone-300 bg-stone-50 p-4">
					<div class="flex flex-col">
						<label class="text-sm font-medium text-gray-700">Name</label>
						<input
							type="text"
							class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
						/>
					</div>

					<div class="flex flex-col">
						<label class="text-sm font-medium text-gray-700">Description</label>
						<textarea
							class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
						></textarea>
					</div>

					<div class="flex flex-col">
						<label class="text-sm font-medium text-gray-700">URL</label>
						<input
							type="url"
							class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
						/>
					</div>
				</div>
			{/if}

			<br /><br /><br /><br />

			<ul>
				<div>mat1</div>
				<div>mat2</div>
				<div>...</div>
			</ul>
		</div>
	{:else}
		<p class="text-red-600">Course not found.</p>
	{/if}
</div>
