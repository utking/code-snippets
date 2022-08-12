<script setup>
import TagItem from './TagItem.vue';

defineProps({
  tags: Array
})

const emits = defineEmits(['tag:clicked'])

const emitClicked = (id) => {
  emits('tag:clicked', id)
}

</script>

<template>
  <ul class="list-group list-group-flush">
    <TagItem
      v-for="tag in tags" v-bind:key="tag.id"
      :id="tag.id"
      :is_active="tag.active"
      :alias="tag.alias"
      :notes="tag.notes"
      @tag:clicked="emitClicked" />
  </ul>
</template>

<style scoped>
</style>
