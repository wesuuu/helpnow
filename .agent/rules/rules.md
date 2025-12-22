---
trigger: always_on
---

- For the frontend, I want you to use Svelte 5 runes and not use deprecated patterns like

```
$: ...
```

to reactivately update variables

https://svelte.dev/docs/svelte/$effect

- If you open any processes like `npm run dev` or `go run .`, make sure you close them after you're done.