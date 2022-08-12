<script setup>
import { ref } from "@vue/reactivity"


defineProps({
  notes: Array
})

const emits = defineEmits(['note:clicked'])

const emitClicked = (id) => {
  emits('note:clicked', id)
}

</script>

<template>
  <ul class="list-group list-group-flush">
    <li v-for="note in notes" v-bind:key="note.id"
      class="list-group-item"
      :class={active:note.active}>
      <div>
        <span type="button" @click="emitClicked(note.id)">
          {{ note.title }}
        </span>
      </div>
    </li>
  </ul>
</template>

<style scoped>
</style>
