<script setup>

const props = defineProps({
  notes: Number,
  alias: String,
  is_active: Boolean,
  id: Number
})

const emits = defineEmits(['tag:clicked'])

const emitClicked = () => {
  emits('tag:clicked', props.id)
}

</script>

<template>
  <li :class="{active: is_active}"
    type="button" @click="emitClicked"
    class="list-group-item d-flex justify-content-between align-items-start">
    <div>
      <span>{{ alias }}</span>
    </div>
    <span class="badge bg-secondary rounded-pill">
      {{ notes }}
    </span>
  </li>
</template>

<style scoped>
</style>
