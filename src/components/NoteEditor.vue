<script setup>
import { ref } from "@vue/reactivity"
import { onBeforeMount } from "@vue/runtime-core"

// const DEFAULT_INDENT = 4

const props = defineProps({
  origNote: Object,
  title: String
})

const note = ref({
  ID: Number,
  Title: String,
  Content: String,
  Tag: String,
})

const emits = defineEmits(['note:save', 'note:close', 'note:delete'])

const resetModel = () => {
  if (props.origNote && props.origNote.ID) {
    note.value.ID = props.origNote.ID
    note.value.Title = props.origNote.Title
    note.value.Content = props.origNote.Content
    note.value.Tag = props.origNote.Tag
  } else {
    note.value.ID = undefined
    note.value.Title = ''
    note.value.Content = ''
    note.value.Tag = ''
  }
}

const emitSave = () => {
  emits("note:save", note.value)
}

const emitClose = () => {
  note.value.ID = undefined
  note.value.Title = ''
  note.value.Content = ''
  note.value.Tag = ''
  emits("note:close")
}

const emitDelete = () => {
  emits("note:delete")
}

onBeforeMount(() => {
  resetModel()
})

</script>

<template>
  <div class="card">
    <div class="card-header">
      <span>{{title}} | </span>
      <span v-if="note.ID">Update</span><span v-else>New</span> snippet
    </div>
    <div class="card-body">
      <div class="form mb-1">
        <div class="row mb-1">
          <label for="note-title" class="col-sm-2 col-form-label text-end" aria-required="true">
            Title
          </label>
          <div class="col-sm-10">
            <input type="text" class="form-control" id="note-title"
              placeholder="Snippet title" required v-model="note.Title">
          </div>
        </div>
        <div class="row mb-1">
          <label for="note-tag" class="col-sm-2 col-form-label text-end" aria-required="true">
            Tag
          </label>
          <div class="col-sm-10">
            <input type="text" class="form-control" id="note-tag"
              placeholder="Snippet tag" required v-model="note.Tag">
          </div>
        </div>
      </div>

      <!-- body -->
      <div class="card">
        <div class="card-body" id="note-content">
          <textarea name="note-content" id="note-content" class="form-control"
            placeholder="Paste a snippet of code..." v-model="note.Content"></textarea>
        </div>
        <div class="card-footer">
          <div class="d-grid gap-2 d-md-flex justify-content-md-end">
            <button type="button" class="btn btn-sm btn-danger" @click="emitClose">
              Cancel
            </button>
            <button type="button" class="btn btn-sm btn-success float-right" @click="emitSave">
              <span v-if="note.ID">Update</span><span v-else>Create</span>
            </button>
          </div>
        </div>
      </div>
      <!-- end of body -->

    </div>
  </div>
</template>

<style scoped>
#note-content {
  padding: 2px;
  border-radius: 0;
}
textarea {
  min-height: 70vh;
}
</style>
