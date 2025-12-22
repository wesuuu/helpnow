import { writable } from 'svelte/store';

const storedTheme = localStorage.getItem('theme') || 'light';

export const theme = writable(storedTheme);

theme.subscribe(value => {
    localStorage.setItem('theme', value);
    if (value === 'dark') {
        document.documentElement.classList.add('dark');
    } else {
        document.documentElement.classList.remove('dark');
    }
});

export const toggleTheme = () => {
    theme.update(current => current === 'dark' ? 'light' : 'dark');
};
