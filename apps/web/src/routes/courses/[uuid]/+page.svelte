<script lang="ts">
    import { page } from '$app/state';
    import ViewMaterial from './ViewMaterial.svelte';
    import DataLoader from '$lib/components/DataLoader.svelte';
    import type { Course } from '$lib/types';
    import TakeQuiz from './TakeQuiz.svelte';
    import ViewFeed from './ViewFeed.svelte';

    let coursesPromise: Promise<any> = loadCourseDetail();

    // Svelte 5 state for the tabs
    let activeTab = $state('materials'); 

    async function loadCourseDetail() {
        const res = await fetch('/api/courses/' + page.params.uuid);
        if (!res.ok) throw new Error('Course not found');
        return res.json();
    }
</script>

<!-- <svelte:head>
	<title>{course.name}</title>
</svelte:head> -->

<DataLoader promise={coursesPromise}>
    {#snippet children(course: Course)}

        <div class="min-h-screen bg-s-white p-4 md:p-10">
            <div class="mx-auto flex w-full flex-col gap-6 md:max-w-4xl">
                
                <div class="space-y-4">
                    <h1 class="text-4xl font-bold text-s-black md:text-6xl">{course.name}</h1>
                    <div class="rounded-xl bg-p-blue p-1.5 shadow-lg">
                        <div class="rounded-lg bg-p-green p-4">
                            <p class="text-xl text-s-black md:text-2xl">{course.description}</p>
                        </div>
                    </div>
                </div>

                <div class="sticky top-2 z-10 flex gap-2 rounded-xl bg-s-black p-1 shadow-md">
                    {#each ['materials', 'quizzes', 'feed'] as tab}
                        <button
                            onclick={() => activeTab = tab}
                            class="flex-1 rounded-lg py-2 text-lg font-bold transition-all capitalize
                            {activeTab === tab ? 'bg-p-green text-s-black' : 'text-s-white hover:bg-p-blue'}"
                        >
                            {tab}
                        </button>
                    {/each}
                </div>

                <div class="mt-2 min-h-[400px]">
                    {#if activeTab === 'materials'}
                        <div class="space-y-3">
                            <h2 class="text-2xl font-bold">Course Materials</h2>
                            {#if course.materials?.length}
                                <div class="rounded-xl bg-p-blue p-2 space-y-2">
                                    {#each course.materials as material}
                                        <ViewMaterial {material} />
                                    {/each}
                                </div>
                            {:else}
                                <p class="text-gray-500 italic">No materials available yet.</p>
                            {/if}
                        </div>

                    {:else if activeTab === 'quizzes'}
                        <div class="space-y-3">
                            <h2 class="text-2xl font-bold">Available Quizzes</h2>
                            {#if course.quizzes?.length}
                                <div class="rounded-xl bg-p-blue p-2 space-y-2">
                                    {#each course.quizzes as quiz}
                                        <TakeQuiz {quiz} courseId={course.uuid} />
                                    {/each}
                                </div>
                            {:else}
                                <p class="text-gray-500 italic">No quizzes assigned to this course.</p>
                            {/if}
                        </div>

                    {:else if activeTab === 'feed'}
                        <div class="space-y-3">
                            <h2 class="text-2xl font-bold">News Feed</h2>
                            <ViewFeed courseId={course.uuid} />

                        </div>
                    {/if}
                </div>
            </div>
        </div>
    {/snippet}
</DataLoader>