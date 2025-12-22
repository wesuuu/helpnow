/// <reference types="svelte" />
/// <reference types="vite/client" />

declare module '*.svelte' {
    import { Component } from 'svelte';
    const component: Component;
    export default component;
}
