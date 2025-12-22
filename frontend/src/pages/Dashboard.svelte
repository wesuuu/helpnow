<script lang="ts">
  let stats = [
    { label: 'Active Agents', value: '3', change: '+12%', changeType: 'positive' },
    { label: 'Routines Run', value: '142', change: '+5.4%', changeType: 'positive' },
    { label: 'Pending Reviews', value: '4', change: 'Urgent', changeType: 'negative' },
  ];

  let recentActivity = [
    { agent: 'Marketing Bot', action: 'Created Email Campaign', time: '2 mins ago', status: 'Success' },
    { agent: 'SEO Optimizer', action: 'Analyzed Competitor Site', time: '1 hour ago', status: 'Success' },
    { agent: 'Content Gen', action: 'Drafted Blog Post', time: '4 hours ago', status: 'Pending Review' },
  ];
</script>

<div>
  <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">Dashboard</h2>

  <!-- Stats Grid -->
  <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
    {#each stats as stat}
      <div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700">
        <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">{stat.label}</h3>
        <div class="mt-2 flex items-baseline">
          <span class="text-3xl font-semibold text-gray-900 dark:text-white">{stat.value}</span>
          <span class={`ml-2 text-sm font-medium ${stat.changeType === 'positive' ? 'text-green-600' : 'text-red-600'}`}>
            {stat.change}
          </span>
        </div>
      </div>
    {/each}
  </div>

  <!-- Recent Activity -->
  <div class="bg-white dark:bg-gray-800 shadow-sm rounded-lg border border-gray-200 dark:border-gray-700">
    <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
      <h3 class="text-lg font-medium text-gray-900 dark:text-white">Recent Activity</h3>
    </div>
    <ul class="divide-y divide-gray-200 dark:divide-gray-700">
      {#each recentActivity as item}
        <li class="px-6 py-4 flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-900 dark:text-white">{item.action}</p>
            <p class="text-sm text-gray-500 dark:text-gray-400">by {item.agent}</p>
          </div>
          <div class="flex items-center">
            <span class={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
              item.status === 'Success' ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' : 
              item.status === 'Pending Review' ? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200' : 
              'bg-gray-100 text-gray-800'}`}>
              {item.status}
            </span>
            <span class="ml-4 text-sm text-gray-500 dark:text-gray-400">{item.time}</span>
          </div>
        </li>
      {/each}
    </ul>
  </div>
</div>
