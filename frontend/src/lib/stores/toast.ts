import { writable } from "svelte/store";

export interface Toast {
    id: number;
    message: string;
    type: "success" | "error" | "info";
    duration: number;
}

function createToastStore() {
    const { subscribe, update } = writable<Toast[]>([]);

    let nextId = 1;

    function add(
        message: string,
        type: "success" | "error" | "info" = "info",
        duration: number = 3000,
    ) {
        const id = nextId++;
        update((toasts) => [...toasts, { id, message, type, duration }]);

        if (duration > 0) {
            setTimeout(() => {
                remove(id);
            }, duration);
        }
    }

    function remove(id: number) {
        update((toasts) => toasts.filter((t) => t.id !== id));
    }

    return {
        subscribe,
        add,
        remove,
        success: (msg: string, duration?: number) =>
            add(msg, "success", duration),
        error: (msg: string, duration?: number) =>
            add(msg, "error", duration),
        info: (msg: string, duration?: number) =>
            add(msg, "info", duration),
    };
}

export const toast = createToastStore();
