export class Router {
    path = $state(window.location.pathname);

    guard: (() => boolean) | null = null;

    constructor() {
        window.addEventListener('popstate', () => {
            // Handling popstate blocking is harder, simple impl for now
            this.path = window.location.pathname;
        });
    }

    navigate(to) {
        if (this.guard) {
            if (!this.guard()) return;
        }
        history.pushState(null, '', to);
        this.path = to;
    }

    registerGuard(fn: () => boolean) {
        this.guard = fn;
    }
    unregisterGuard() {
        this.guard = null;
    }
}

export const router = new Router();
