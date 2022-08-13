<script setup>

const props = defineProps({
  id: Number,
  description: String,
  title: String,
  syntax: String,
  content: String,
  indent: Number,
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
      <div class="form mb-1">
        <div class="row mb-1">
          <label for="note-title" class="col-sm-3 col-form-label" aria-required="true">Title</label>
          <div class="col-sm-9">
            <input type="text" class="form-control" id="note-title"
              placeholder="Snippet title" :value="title" disabled>
          </div>
        </div>

        <div class="row mb-1">
          <label for="note-description" class="col-sm-3 col-form-label">Description</label>
          <div class="col-sm-9">
            <input type="text" class="form-control" id="note-description"
              placeholder="Snippet description" :value="description" disabled>
          </div>
        </div>
      </div>

      <!-- body -->
      <div class="card">
        <div class="card-body">
          <textarea name="note-content" id="note-content" class="form-control" disabled
            rows="8" :value="content"></textarea>
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
</style>
