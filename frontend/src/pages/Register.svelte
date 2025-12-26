<script lang="ts">
  import Link from "../lib/Link.svelte";
  import { router } from "../lib/router.svelte.js";
  import { login } from "../stores/auth.js";

  let firstName = $state("");
  let middleName = $state("");
  let lastName = $state("");
  let email = $state("");
  let password = $state("");

  async function handleSubmit(event: Event) {
    event.preventDefault();
    login("mock-token", {
      first_name: firstName,
      middle_name: middleName,
      last_name: lastName,
      email,
    });
    router.navigate("/dashboard");
  }
</script>

<div
  class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8"
>
  <div
    class="max-w-md w-full space-y-8 bg-white dark:bg-gray-800 p-8 rounded-lg shadow-md"
  >
    <div>
      <h2
        class="mt-6 text-center text-3xl font-extrabold text-gray-900 dark:text-white"
      >
        Create your account
      </h2>
      <p class="mt-2 text-center text-sm text-gray-600 dark:text-gray-400">
        Already have an account? <Link
          to="/login"
          class="font-medium text-indigo-600 hover:text-indigo-500"
          >Sign in</Link
        >
      </p>
    </div>
    <form class="mt-8 space-y-6" onsubmit={handleSubmit}>
      <div class="rounded-md shadow-sm -space-y-px">
        <div>
          <label for="first-name" class="sr-only">First Name</label>
          <input
            id="first-name"
            name="firstName"
            type="text"
            autocomplete="given-name"
            required
            bind:value={firstName}
            class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-700 placeholder-gray-500 text-gray-900 dark:text-white dark:bg-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
            placeholder="First Name"
          />
        </div>
        <div>
          <label for="middle-name" class="sr-only">Middle Name</label>
          <input
            id="middle-name"
            name="middleName"
            type="text"
            autocomplete="additional-name"
            bind:value={middleName}
            class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-700 placeholder-gray-500 text-gray-900 dark:text-white dark:bg-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
            placeholder="Middle Name (Optional)"
          />
        </div>
        <div>
          <label for="last-name" class="sr-only">Last Name</label>
          <input
            id="last-name"
            name="lastName"
            type="text"
            autocomplete="family-name"
            required
            bind:value={lastName}
            class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-700 placeholder-gray-500 text-gray-900 dark:text-white dark:bg-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
            placeholder="Last Name"
          />
        </div>
        <div>
          <label for="email-address" class="sr-only">Email address</label>
          <input
            id="email-address"
            name="email"
            type="email"
            autocomplete="email"
            required
            bind:value={email}
            class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-700 placeholder-gray-500 text-gray-900 dark:text-white dark:bg-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
            placeholder="Email address"
          />
        </div>
        <div>
          <label for="password" class="sr-only">Password</label>
          <input
            id="password"
            name="password"
            type="password"
            autocomplete="new-password"
            required
            bind:value={password}
            class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-700 placeholder-gray-500 text-gray-900 dark:text-white dark:bg-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
            placeholder="Password"
          />
        </div>
      </div>

      <div>
        <button
          type="submit"
          class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          Sign up
        </button>
      </div>
    </form>
  </div>
</div>
