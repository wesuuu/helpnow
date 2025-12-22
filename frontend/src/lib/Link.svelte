<script>
    import { router } from "./router.svelte.js";

    let { to, class: className, children, ...rest } = $props();

    function handleClick(e) {
        e.preventDefault();
        router.navigate(to);
    }
</script>

<a href={to} onclick={handleClick} class={className} {...rest}>
    {@render children?.()}
</a>
