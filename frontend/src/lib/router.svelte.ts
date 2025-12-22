export class Router {
    path = $state(window.location.pathname);

    constructor() {
        window.addEventListener('popstate', () => {
            this.path = window.location.pathname;
        });
    }

    navigate(to) {
        history.pushState(null, '', to);
        this.path = to;
    }
}

export const router = new Router();
