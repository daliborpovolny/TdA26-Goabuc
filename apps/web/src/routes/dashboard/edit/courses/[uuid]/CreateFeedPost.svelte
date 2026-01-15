<script lang="ts">
    import { fade, slide } from 'svelte/transition';

    let { courseId }: { courseId: string } = $props();

    let collapsed = $state(true);
    let isSaving = $state(false);
    let showSuccess = $state(false);

    async function newFeedPost(e: Event) {
        e.preventDefault();
        isSaving = true;

        const form = e.currentTarget as HTMLFormElement;
        const formData = new FormData(form);
        const data = Object.fromEntries(formData);

        let res = await fetch(`/api/courses/${courseId}/feed`, {
            method: 'POST',
            headers: { 'Content-type': 'application/json' },
            body: JSON.stringify(data)
        });

        if (res.ok) {
            showSuccess = true;
            // Clear form and collapse after a short delay
            setTimeout(() => {
                collapsed = true;
                showSuccess = false;
                form.reset();
            }, 1000);
        }
        isSaving = false;
    }
</script>

<div class="mb-6">
    <div class="flex items-center gap-4">
        <button
            type="button"
            class="group relative flex items-center gap-2 rounded-xl border-4 border-s-black px-6 py-3 font-black uppercase tracking-widest transition-all
            {collapsed ? 'bg-white shadow-[4px_4px_0px_0px_rgba(26,26,26,1)] hover:bg-p-green' : 'bg-p-blue text-white translate-x-1 translate-y-1 shadow-none'}"
            onclick={() => (collapsed = !collapsed)}
        >
            <span>📢</span> {collapsed ? 'New Announcement' : 'Cancel Post'}
        </button>

        {#if showSuccess}
            <span transition:fade class="font-bold text-p-green uppercase tracking-tighter">
                ✓ Post Broadcasted
            </span>
        {/if}
    </div>

    {#if !collapsed}
        <div transition:slide class="mt-6 relative max-w-2xl">
            <div class="absolute inset-0 translate-x-3 translate-y-3 rounded-2xl bg-s-black"></div>
            
            <div class="relative rounded-2xl border-4 border-s-black bg-white p-6">
                <form onsubmit={newFeedPost} class="space-y-4">
                    <div class="space-y-1">
                        <label for="message" class="text-xs font-black uppercase tracking-widest text-gray-500">
                            Message to Students
                        </label>
                        <textarea
                            id="message"
                            name="message"
                            class="min-h-[120px] w-full resize-none rounded-xl border-4 border-s-black bg-white p-4 text-xl font-bold text-s-black focus:outline-none focus:ring-4 focus:ring-p-green"
                            placeholder="Type your update here..."
                            required
                        ></textarea>
                    </div>

                    <div class="flex justify-end pt-2">
                        <button 
                            type="submit" 
                            disabled={isSaving}
                            class="group relative overflow-hidden rounded-xl border-4 border-s-black bg-p-green px-8 py-3 text-xl font-black uppercase tracking-widest text-s-black transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none active:translate-x-2 active:translate-y-2 disabled:opacity-50"
                        >
                            <div class="absolute inset-0 translate-x-1 translate-y-1 bg-s-black opacity-0 group-hover:opacity-10"></div>
                            {isSaving ? 'Sending...' : 'Post Update →'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    {/if}
</div>