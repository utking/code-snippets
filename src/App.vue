<script setup>
import Navbar from './components/Navbar.vue';
// This starter template is using Vue 3 <script setup> SFCs
// Check out https://vuejs.org/api/sfc-script-setup.html#script-setup
</script>

<template>

  <Navbar />

  <div class="container-fluid notes">
    <router-view />
  </div>

</template>

<style scoped>
.notes {
  margin-top: 1rem;
  padding-right: 1rem;
  padding-left: 1rem;
}
</style>
