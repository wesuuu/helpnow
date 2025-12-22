import { writable } from 'svelte/store';

const storedToken = localStorage.getItem('token');
const storedUser = localStorage.getItem('user');

export const auth = writable({
    isAuthenticated: !!storedToken,
    token: storedToken,
    user: storedUser ? JSON.parse(storedUser) : null
});

export const login = (token: string, user: any) => {
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(user));
    auth.set({ isAuthenticated: true, token, user });
};

export const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    auth.set({ isAuthenticated: false, token: null, user: null });
};
