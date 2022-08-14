<script setup>

const props = defineProps({
  id: Number,
  title: String,
  content: String,
  indent: Number,
  tagAlias: String
})

const emits = defineEmits(['note:edit', 'note:delete'])

const emitEdit = () => {
  emits("note:edit", props.id)
}

const emitDelete = () => {
  emits("note:delete", props.id)
}

</script>

<template>
  <div class="card">
    <div class="card-body">
      <h5 class="title h5">{{tagAlias}} | {{title}}</h5>
      
      <!-- body -->
      <div class="card">
        <div class="card-body">
          <pre v-highlightjs><code class="automatic">{{ content }}</code></pre>
        </div>
        <div class="card-footer">
          <div class="d-grid gap-2 d-md-flex justify-content-md-end">
            <button type="button" class="btn btn-sm btn-info" @click="emitEdit">
              Edit
            </button>
            <button type="button" class="btn btn-sm btn-danger float-right" @click="emitDelete">
              Delete
            </button>
          </div>
        </div>
      </div>
      <!-- end of body -->

    </div>
  </div>
</template>

<style scoped>
.hljs {
  min-height: 70vh;
}
</style>
