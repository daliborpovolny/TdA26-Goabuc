<script lang="ts">
    import DataLoader from '$lib/components/DataLoader.svelte';
    import type { Course } from '$lib/types';

    let coursesPromise: Promise<Course[]> = loadCourses();

    async function loadCourses() {
        const res = await fetch('/api/courses');
        if (!res.ok) throw new Error('Failed to get list of courses');
        return res.json();
    }
</script>

<svelte:head>
    <title>Courses | TdA</title>
</svelte:head>

<div class="min-h-screen bg-s-white p-6 md:p-12">
    <header class="mb-12 text-center">
        <h1 class="inline-block text-5xl font-black uppercase tracking-tighter md:text-7xl">
            Our <span class="text-p-blue">Courses</span>
        </h1>
        <div class="mx-auto mt-2 h-2 w-24 bg-p-green"></div>
    </header>

    <div class="mx-auto max-w-6xl">
        <DataLoader promise={coursesPromise}>
            {#snippet children(courses: Course[])}
                <div class="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3">
                    {#each courses as course}
                        <a 
                            href="/courses/{course.uuid}" 
                            class="group relative block transition-transform hover:-translate-x-1 hover:-translate-y-1"
                        >
                            <div class="absolute inset-0 translate-x-2 translate-y-2 rounded-2xl bg-s-black"></div>
                            
                            <div class="relative flex h-full flex-col border-2 border-s-black bg-white p-6 transition-colors group-hover:bg-p-green">
                                <div class="mb-4 flex items-start justify-between">
                                    <span class="rounded-lg border-2 border-s-black bg-p-blue px-3 py-1 text-xs font-bold text-white uppercase">
                                        OPEN
                                    </span>
                                    <span class="text-3xl">🎓</span>
                                </div>

                                <h2 class="mb-4 text-3xl font-black leading-tight text-s-black group-hover:text-s-black">
                                    {course.name}
                                </h2>

                                <div class="mt-auto pt-6">
                                    <div class="flex items-center gap-2 font-bold text-s-black uppercase tracking-widest">
                                        View Course 
                                        <span class="transition-transform group-hover:translate-x-2">→</span>
                                    </div>
                                </div>
                            </div>
                        </a>
                    {/each}
                </div>
            {/snippet}
        </DataLoader>
    </div>
</div>