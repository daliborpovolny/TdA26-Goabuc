// import type { Actions, PageServerLoad } from './$types';
// import { fail } from '@sveltejs/kit';

// export const load: PageServerLoad = async ({ params, fetch }) => {
// 	// Initial fetch to populate page before SSE kicks in
// 	const res = await fetch(`/api/courses/${params.uuid}/feed`);
// 	const initialPosts = await res.json();
// 	return { initialPosts };
// };

// export const actions: Actions = {
// 	create: async ({ request, params, fetch }) => {
// 		const data = await request.formData();
// 		const message = data.get('message');

// 		const res = await fetch(`/api/courses/${params.uuid}/feed`, {
// 			method: 'POST',
// 			headers: { 'Content-Type': 'application/json' },
// 			body: JSON.stringify({ message })
// 		});

// 		if (!res.ok) return fail(res.status, { error: 'Failed to create' });
// 	},

// 	edit: async ({ request, params, fetch }) => {
// 		const data = await request.formData();
// 		const message = data.get('message');
// 		const uuid = data.get('uuid'); // Getting UUID from hidden input

// 		const res = await fetch(`/api/courses/${params.uuid}/feed/${uuid}`, {
// 			method: 'PUT',
// 			headers: { 'Content-Type': 'application/json' },
// 			body: JSON.stringify({ message, edited: true })
// 		});

// 		if (!res.ok) return fail(res.status, { error: 'Failed to edit' });
// 	},

// 	delete: async ({ request, params, fetch }) => {
// 		const data = await request.formData();
// 		const uuid = data.get('uuid');

// 		const res = await fetch(`/api/courses/${params.uuid}/feed/${uuid}`, {
// 			method: 'DELETE'
// 		});

// 		if (!res.ok) return fail(res.status, { error: 'Failed to delete' });
// 	}
// };
