<script lang="ts" setup>
import DashNav from '@/components/dashboard/DashNav.vue'
</script>

<template>
  <slot name="modal" />
  <div class="drawer drawer-mobile">
    <input id="drawer" type="checkbox" class="drawer-toggle" />
    <div class="drawer-content">
      <DashNav />
      <div class="mx-4">
        <slot />
      </div>
    </div>
    <div class="drawer-side">
      <label for="drawer" class="drawer-overlay" />
      <ul class="menu menu-compact p-4 w-80 bg-base-200 text-base-content">
        <!-- Sidebar content here -->
        <li>
          <RouterLink to="/dashboard" active-class="active">Dashboard</RouterLink>
        </li>
        <li>
          <RouterLink to="/dashboard/boards" active-class="active">Boards</RouterLink>
        </li>
      </ul>
    </div>
  </div>
</template>
