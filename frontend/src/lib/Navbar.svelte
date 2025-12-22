<script lang="ts">
  import { toggleTheme, theme } from "../stores/theme.js";
  import { auth, logout } from "../stores/auth.js";
  import { router } from "./router.svelte.js";

  function handleLogout() {
    logout();
    router.navigate("/login");
  }
</script>

<header
  class="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700 h-16 flex items-center justify-between px-6 z-10"
>
  <div class="flex items-center">
    <!-- Mobile menu button could go here -->
    <h1
      class="text-xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent"
    >
      HelpNow AI
    </h1>
  </div>

  <div class="flex items-center space-x-4">
    <button
      onclick={toggleTheme}
      class="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
      aria-label="Toggle Theme"
    >
      {#if $theme === "dark"}
        <!-- Sun Icon -->
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5 text-gray-400"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"
          />
        </svg>
      {:else}
        <!-- Moon Icon -->
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5 text-gray-500"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
          />
        </svg>
      {/if}
    </button>

    {#if $auth.isAuthenticated}
      <div class="relative group">
        <button class="flex items-center space-x-2 focus:outline-none">
          <div
            class="w-8 h-8 rounded-full bg-indigo-500 flex items-center justify-center text-white font-medium"
          >
            {$auth.user?.full_name ? $auth.user.full_name[0] : "U"}
          </div>
          <span class="text-sm font-medium hidden md:block"
            >{$auth.user?.full_name || "User"}</span
          >
        </button>
        <!-- Dropdown -->
        <!-- Dropdown -->
        <div
          class="absolute right-0 top-full pt-2 hidden group-hover:block hover:block w-48 z-50"
        >
          <div
            class="bg-white dark:bg-gray-800 rounded-md shadow-lg py-1 border border-gray-200 dark:border-gray-700"
          >
            <button
              onclick={() => router.navigate("/organization-settings")}
              class="w-full text-left block px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
              >Organization</button
            >
            <button
              onclick={() => router.navigate("/team-settings")}
              class="w-full text-left block px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
              >Team</button
            >
            <button
              onclick={() => router.navigate("/user-settings")}
              class="w-full text-left block px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
              >User Settings</button
            >
            <div
              class="border-t border-gray-200 dark:border-gray-700 my-1"
            ></div>
            <button
              onclick={handleLogout}
              class="w-full text-left block px-4 py-2 text-sm text-red-600 hover:bg-gray-100 dark:hover:bg-gray-700"
              >Logout</button
            >
          </div>
        </div>
      </div>
    {/if}
  </div>
</header>
