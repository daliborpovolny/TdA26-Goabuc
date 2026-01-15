<script lang="ts">
    import { auth } from '$lib/auth.svelte';
    import { fade, fly } from 'svelte/transition';

    let dataPromise = fetch('/api/me').then((r) => {
        if (!r.ok) throw new Error('Not logged in');
        return r.json();
    });
</script>

<svelte:head>
    <title>Home | TdA</title>
</svelte:head>

<main class="min-h-screen bg-s-white">
    <section class="relative overflow-hidden border-b-8 border-s-black bg-p-blue px-6 py-20 text-white">
        <div class="mx-auto max-w-6xl">
            <div class="max-w-3xl space-y-6">
                <h1 class="text-6xl font-black uppercase tracking-tighter md:text-8xl" in:fly={{ y: 20 }}>
                    Tour de <span class="text-p-green">App!</span>
                </h1>
                
                {#await dataPromise}
                    <div class="h-20 w-20 animate-pulse rounded-lg bg-white/20"></div>
                {:then user}
                    <div in:fade class="space-y-6">
                        <h2 class="text-3xl font-bold md:text-4xl">
                            Welcome back, {user.firstName}!
                        </h2>
                        <div class="flex flex-wrap gap-4">
                            <a href="/dashboard" class="rounded-xl border-4 border-s-black bg-p-green px-8 py-4 text-2xl font-black uppercase text-s-black shadow-[6px_6px_0px_0px_rgba(26,26,26,1)] transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none">
                                Open Dashboard
                            </a>
                        </div>
                    </div>
                {:catch}
                    <div in:fade class="space-y-6">
                        <p class="text-2xl font-medium opacity-90">
                            The ultimate academy for modern developers. Master your craft with interactive courses and real-time feeds.
                        </p>
                        <div class="flex flex-wrap gap-4 pt-4">
                            <a href="/register" class="rounded-xl border-4 border-s-black bg-p-green px-8 py-4 text-2xl font-black uppercase text-s-black shadow-[6px_6px_0px_0px_rgba(26,26,26,1)] transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none">
                                Get Started
                            </a>
                            <a href="/login" class="rounded-xl border-4 border-s-black bg-white px-8 py-4 text-2xl font-black uppercase text-s-black shadow-[6px_6px_0px_0px_rgba(26,26,26,1)] transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none">
                                Member Login
                            </a>
                        </div>
                    </div>
                {/await}
            </div>
        </div>

        <div class="absolute -right-10 bottom-10 hidden rotate-12 md:block">
            <div class="rounded-full border-4 border-s-black bg-s-2 p-10 text-4xl font-black uppercase shadow-lg">
                V0.1
            </div>
        </div>
    </section>

    <section class="mx-auto max-w-6xl px-6 py-20">
        <div class="grid grid-cols-1 gap-12 md:grid-cols-3">
            <div class="space-y-4">
                <span class="text-5xl">🎓</span>
                <h3 class="text-2xl font-black uppercase">Learn</h3>
                <p class="font-bold text-gray-500">Access structured materials and files instantly.</p>
            </div>
            <div class="space-y-4">
                <span class="text-5xl">📝</span>
                <h3 class="text-2xl font-black uppercase">Quiz</h3>
                <p class="font-bold text-gray-500">Validate your skills in an engaging way.</p>
            </div>
            <div class="space-y-4">
                <span class="text-5xl">📣</span>
                <h3 class="text-2xl font-black uppercase">Feed</h3>
                <p class="font-bold text-gray-500">Stay updated with real-time teacher announcements.</p>
            </div>
        </div>
    </section>
</main>