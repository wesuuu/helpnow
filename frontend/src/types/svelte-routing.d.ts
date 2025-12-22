declare module 'svelte-routing' {
    import { SvelteComponent } from 'svelte';

    export interface RouterProps {
        url?: string;
        [key: string]: any;
    }
    export class Router extends SvelteComponent<RouterProps> { }

    export interface RouteProps {
        path?: string;
        component?: any;
        [key: string]: any;
    }
    export class Route extends SvelteComponent<RouteProps> { }

    export interface LinkProps {
        to: string;
        replace?: boolean;
        [key: string]: any;
    }
    export class Link extends SvelteComponent<LinkProps> { }

    export function navigate(to: string, { replace, state }?: { replace?: boolean; state?: any }): void;
    export function link(node: HTMLElement): { destroy(): void };
    export function useLocation(): { pathname: { subscribe: (run: (value: string) => void) => () => void } };
}
