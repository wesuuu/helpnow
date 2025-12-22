<script lang="ts">
  import { onMount } from "svelte";

  // State
  let orgName = $state("");
  let systemPrompt = $state("");
  let plan = $state("Enterprise");
  let message = $state("");

  onMount(async () => {
    // Fetch Org Settings (ID 1 hardcoded for MVP)
    const res = await fetch("/api/organizations/1");
    if (res.ok) {
      const data = await res.json();
      orgName = data.name;
      systemPrompt = data.system_prompt || "";
    }
  });

  async function saveSettings() {
    const res = await fetch("/api/organizations/1", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name: orgName,
        system_prompt: systemPrompt,
      }),
    });

    if (res.ok) {
      message = "Settings saved successfully!";
      setTimeout(() => (message = ""), 3000);
    } else {
      message = "Failed to save settings.";
    }
  }
</script>

<div>
  <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
    Organization Settings
  </h2>

  <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6 max-w-2xl">
    <div class="grid grid-cols-1 gap-6">
      <div>
        <label
          for="orgName"
          class="block text-sm font-medium text-gray-700 dark:text-gray-300"
          >Organization Name</label
        >
        <input
          type="text"
          id="orgName"
          bind:value={orgName}
          class="mt-1 block w-full rounded-md border-gray-300 dark:border-gray-600 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 dark:bg-gray-700 dark:text-white sm:text-sm px-3 py-2 border"
        />
      </div>

      <div>
        <label
          for="systemPrompt"
          class="block text-sm font-medium text-gray-700 dark:text-gray-300"
        >
          System Prompt
          <span
            class="block text-xs font-normal text-gray-500 dark:text-gray-400 mt-1"
            >This context is shared with all AI agents in your organization.
            Define your business, tone, and key information here.</span
          >
        </label>
        <textarea
          id="systemPrompt"
          bind:value={systemPrompt}
          rows="6"
          class="mt-1 block w-full rounded-md border-gray-300 dark:border-gray-600 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 dark:bg-gray-700 dark:text-white sm:text-sm px-3 py-2 border"
          placeholder="e.g., We are HelpNow, a customer support platform. Our tone is friendly and professional..."
        ></textarea>
      </div>

      <div>
        <label
          for="plan"
          class="block text-sm font-medium text-gray-700 dark:text-gray-300"
          >Current Plan</label
        >
        <select
          id="plan"
          bind:value={plan}
          class="mt-1 block w-full rounded-md border-gray-300 dark:border-gray-600 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 dark:bg-gray-700 dark:text-white sm:text-sm px-3 py-2 border"
        >
          <option>Free</option>
          <option>Pro</option>
          <option>Enterprise</option>
        </select>
      </div>

      <div class="flex items-center justify-between">
        <button
          onclick={saveSettings}
          class="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700"
          >Save Changes</button
        >
        {#if message}
          <span class="text-sm font-medium text-green-600 dark:text-green-400"
            >{message}</span
          >
        {/if}
      </div>
    </div>
  </div>
</div>
