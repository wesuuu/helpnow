<script lang="ts">
  let teamMembers = [
    { name: 'Alice Smith', role: 'Admin', email: 'alice@example.com' },
    { name: 'Bob Jones', role: 'Editor', email: 'bob@example.com' },
  ];
</script>

<div>
  <div class="flex justify-between items-center mb-6">
    <h2 class="text-2xl font-bold text-gray-900 dark:text-white">Team Settings</h2>
    <button class="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 text-sm">Invite Member</button>
  </div>

  <div class="bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-md">
    <ul class="divide-y divide-gray-200 dark:divide-gray-700">
      {#each teamMembers as member}
        <li class="px-6 py-4 flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-900 dark:text-white">{member.name}</p>
            <p class="text-sm text-gray-500 dark:text-gray-400">{member.email}</p>
          </div>
          <div class="flex items-center space-x-4">
            <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">
              {member.role}
            </span>
            <button class="text-red-600 hover:text-red-900 text-sm">Remove</button>
          </div>
        </li>
      {/each}
    </ul>
  </div>
</div>
