<script lang="ts">
    import { page } from '$app/state';
    
    let { name, url, onclick }: { name: string; url: string; onclick?: () => void } = $props();

    // Determine if this link is currently active
    let isActive = $derived(page.url.pathname === url);
</script>

<a
    href={url}
    {onclick}
    class="flex w-full items-center px-6 py-5 font-black uppercase tracking-widest transition-all md:h-[72px] md:w-auto md:px-5 md:py-0 md:text-sm
    {isActive 
        ? 'bg-p-green text-s-black' 
        : 'text-white hover:bg-p-green/20 hover:text-p-green'}"
>
    {name}
</a>