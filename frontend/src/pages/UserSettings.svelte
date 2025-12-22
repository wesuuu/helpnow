<script lang="ts">
  import { auth } from '../stores/auth.js';
  let fullName = $state($auth.user?.full_name || '');
  let email = $state($auth.user?.email || '');
</script>

<div>
  <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">User Settings</h2>

  <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6 max-w-2xl">
    <div class="grid grid-cols-1 gap-6">
      <div class="flex items-center space-x-4 mb-4">
        <div class="h-16 w-16 bg-indigo-500 rounded-full flex items-center justify-center text-white text-2xl font-bold">
          {fullName ? fullName[0] : 'U'}
        </div>
        <button class="text-sm text-indigo-600 hover:text-indigo-500">Change Avatar</button>
      </div>

      <div>
        <label for="fullName" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Full Name</label>
        <input type="text" id="fullName" bind:value={fullName} class="mt-1 block w-full rounded-md border-gray-300 dark:border-gray-600 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 dark:bg-gray-700 dark:text-white sm:text-sm px-3 py-2 border">
      </div>
      
      <div>
        <label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Email Address</label>
        <input type="email" id="email" bind:value={email} class="mt-1 block w-full rounded-md border-gray-300 dark:border-gray-600 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 dark:bg-gray-700 dark:text-white sm:text-sm px-3 py-2 border">
      </div>

      <div>
        <button class="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700">Update Profile</button>
      </div>
    </div>
  </div>
</div>
