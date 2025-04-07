<script setup lang="ts">
</script>

<template>
  <div class="min-h-screen flex flex-col bg-gray-50 text-gray-800 w-full">
    <!-- Header -->
    <header class="bg-white shadow-sm sticky top-0 z-50 w-full">
      <div class="w-full px-6 py-4 flex justify-between items-center">
        <RouterLink to="/" class="text-2xl font-extrabold tracking-tight text-blue-600 hover:text-blue-700 transition">
          MyApp
        </RouterLink>

        <nav class="hidden md:flex space-x-6 text-sm font-medium">
          <RouterLink to="/" class="text-gray-600 hover:text-blue-500" active-class="text-blue-600 font-semibold" exact>
            Home
          </RouterLink>
          <RouterLink to="/about" class="text-gray-600 hover:text-blue-500" active-class="text-blue-600 font-semibold">
            About
          </RouterLink>
          <RouterLink to="/contact" class="text-gray-600 hover:text-blue-500" active-class="text-blue-600 font-semibold">
            Contact
          </RouterLink>
        </nav>
      </div>
    </header>

    <!-- Page content -->
    <main class="flex-1 px-6 py-10 w-full">
      <div class="w-full mx-auto">
        <RouterView />
      </div>
    </main>

    <!-- Footer -->
    <footer class="bg-white border-t text-sm text-gray-500 py-6 text-center w-full">
      © 2025 <span class="font-semibold text-gray-700">MyApp</span>. All rights reserved.
    </footer>
  </div>
</template>
