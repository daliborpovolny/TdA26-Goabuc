<!-- UGlY ASS PAGE EXAMPLE OF HOW NOT TO DO IT -->

<script lang="ts">
	import { page } from '$app/state';

	let course = $state<any | null>(null);
	$inspect(course);

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

	async function createFileMaterial(e: SubmitEvent) {
		e.preventDefault();

		const form = e.target as HTMLFormElement;
		const formData = new FormData(form);

		const res = await fetch(`/api/courses/${course.uuid}/materials`, {
			method: 'POST',
			body: formData
		});

		if (!res.ok) {
			console.log(await res.text());
		} else {
			console.log('uploaded');
		}

		loadCourseDetail();
	}

	async function updateFileMaterial(e: SubmitEvent) {
		e.preventDefault();

		const form = e.target as HTMLFormElement;
		const formData = new FormData(form);

		const res = await fetch(`/api/courses/${course.uuid}/materials/uuidTODO`, {
			method: 'PUT',
			body: formData
		});

		if (!res.ok) {
			console.log(await res.text());
		} else {
			console.log('updated');
		}

		loadCourseDetail();
	}

	async function createUrlMaterial(e: SubmitEvent) {
		e.preventDefault();

		const form = e.target as HTMLFormElement;
		const formData = Object.fromEntries(new FormData(form).entries());

		const res = await fetch(`/api/courses/${course.uuid}/materials`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(formData)
		});

		if (!res.ok) {
			console.log(await res.text());
		} else {
			console.log('uploaded');
		}

		loadCourseDetail();
	}

	async function updateUrlMaterial(e: SubmitEvent) {
		e.preventDefault();

		const form = e.target as HTMLFormElement;
		const formData = Object.fromEntries(new FormData(form).entries());

		const res = await fetch(`/api/courses/${course.uuid}/materials`, {
			method: 'PUT',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(formData)
		});

		if (!res.ok) {
			console.log(await res.text());
		} else {
			console.log('updated');
		}

		loadCourseDetail();
	}
</script>

<div class="mx-auto max-w-2xl space-y-8 p-6">
	{#if loadingCourse}
		<p class="text-gray-600">Loading course detail...</p>
	{:else if course != null}
		<h2 class="mb-4 text-4xl font-semibold text-gray-800">Edit Course</h2>

		<div class="rounded-lg border border-stone-300 bg-stone-50 p-6">
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
			<h1 class="mb-3 text-2xl font-semibold text-gray-800">Add Materials</h1>
			<div class="flex gap-4">
				<button
					onclick={() =>
						(createNewMaterial =
							createNewMaterial == 'url' || createNewMaterial == 'idle' ? 'file' : 'idle')}
					class="rounded-md border border-stone-400 bg-stone-100 px-4 py-2 text-gray-800 hover:bg-stone-200 {createNewMaterial ===
					'file'
						? 'bg-stone-200'
						: ''}"
				>
					New File Material
				</button>
				<button
					onclick={() =>
						(createNewMaterial =
							createNewMaterial == 'file' || createNewMaterial == 'idle' ? 'url' : 'idle')}
					class="rounded-md border border-stone-400 bg-stone-100 px-4 py-2 text-gray-800 hover:bg-stone-200 {createNewMaterial ===
					'url'
						? 'bg-stone-200'
						: ''}"
				>
					New Url Material
				</button>
			</div>

			<br />

			{#if createNewMaterial === 'file'}
				<form
					class="space-y-4 rounded-lg border border-stone-300 bg-stone-50 p-4"
					method="POST"
					enctype="multipart/form-data"
					onsubmit={createFileMaterial}
				>
					<!-- required by API -->
					<input type="hidden" name="type" value="file" />

					<div class="flex flex-col">
						<label class="text-sm font-medium text-gray-700" for="name">Name</label>
						<input
							type="text"
							name="name"
							required
							class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
						/>
					</div>

					<div class="flex flex-col">
						<label class="text-sm font-medium text-gray-700" for="description">Description</label>
						<textarea
							name="description"
							class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
						></textarea>
					</div>

					<div class="flex flex-col">
						<label class="text-sm font-medium text-gray-700" for="file">File</label>
						<input type="file" name="file" required class="mt-1 text-gray-900" />
					</div>

					<button class="rounded-md bg-stone-800 px-4 py-2 text-white hover:bg-stone-700">
						Create
					</button>
				</form>
			{:else if createNewMaterial == 'url'}
				<form
					method="POST"
					onsubmit={createUrlMaterial}
					class="space-y-4 rounded-lg border border-stone-300 bg-stone-50 p-4"
				>
					<input type="hidden" name="type" value="url" />

					<div class="flex flex-col">
						<label class="text-sm font-medium text-gray-700" for="name">Name</label>
						<input
							type="text"
							name="name"
							class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
						/>
					</div>

					<div class="flex flex-col">
						<label class="text-sm font-medium text-gray-700" for="description">Description</label>
						<textarea
							name="description"
							class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
						></textarea>
					</div>

					<div class="flex flex-col">
						<label class="text-sm font-medium text-gray-700" for="url">URL</label>
						<input
							type="url"
							name="url"
							class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
						/>
					</div>

					<button class="rounded-md bg-stone-800 px-4 py-2 text-white hover:bg-stone-700">
						Create
					</button>
				</form>
			{/if}

			<br />

			{#if course.materials && course.materials.length > 0}
				<h1 class="mb-3 text-2xl font-semibold text-gray-800">Edit Materials</h1>

				<br />

				<ul>
					{#each course.materials as material}
						{#if material.Type === 'file'}
							<form
								class="space-y-4 rounded-lg border border-stone-300 bg-stone-50 p-4"
								method="POST"
								enctype="multipart/form-data"
								onsubmit={updateFileMaterial}
							>
								<!-- required by API -->
								<input type="hidden" name="type" value="file" />

								<div class="flex flex-col">
									<label class="text-sm font-medium text-gray-700" for="name">Name</label>
									<input
										type="text"
										name="name"
										value={material.Name}
										required
										class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
									/>
								</div>

								<div class="flex flex-col">
									<label class="text-sm font-medium text-gray-700" for="description"
										>Description</label
									>
									<textarea
										name="description"
										value={material.Description}
										class="mt-1 rounded-md border border-stone-300 bg-white px-3 py-2 text-gray-900 focus:ring-2 focus:ring-stone-400 focus:outline-none"
									></textarea>
								</div>

								<div class="flex flex-col">
									<label class="text-sm font-medium text-gray-700" for="file">File</label>
									<input type="file" name="file" class="mt-1 text-gray-900" />
								</div>

								<button class="rounded-md bg-stone-800 px-4 py-2 text-white hover:bg-stone-700">
									Update
								</button>

								<button class="rounded-md bg-red-800 px-4 py-2 text-white hover:bg-red-700">
									Remove
								</button>
							</form>
							<br />
						{:else if material.Type === 'url'}
							<p>URL</p>
						{:else}
							<p>Invalid material type</p>
						{/if}
					{/each}
				</ul>
			{/if}

			<br />
			<br />
		</div>
	{:else}
		<p class="text-red-600">Course not found.</p>
	{/if}
</div>
