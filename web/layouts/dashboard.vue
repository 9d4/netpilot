<template>
  <div class="drawer drawer-mobile">
    <input
      id="drawer"
      type="checkbox"
      class="drawer-toggle"
    >
    <div class="drawer-content">
      <DashboardNavbar />
      <div class="mx-4">
        <slot />
      </div>
    </div> 
    <div class="drawer-side">
      <label
        for="drawer"
        class="drawer-overlay"
      /> 
      <ul class="menu menu-compact p-4 w-80 bg-base-200 text-base-content">
        <!-- Sidebar content here -->
        <li>
          <NuxtLink
            to="/dashboard"
            active-class="active"
          >
            Dashboard
          </NuxtLink>
        </li>
        <li>
          <NuxtLink
            to="/dashboard/boards"
            active-class="active"
          >
            Boards
          </NuxtLink>
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">

</script>