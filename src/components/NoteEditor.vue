<script setup>
import { ref } from "@vue/reactivity"
import { onBeforeMount, onMounted } from "@vue/runtime-core"

const DEFAULT_SYNTAX = 'auto'
const DEFAULT_INDENT = 4

const props = defineProps({
  origNote: Object,
  tagId: Number
})

const note = ref({
  ID: Number,
  Title: String,
  Content: String,
  Description: String,
  Syntax: String,
  Indent: Number,
})

const emits = defineEmits(['note:save', 'note:close', 'note:delete'])

const resetModel = () => {
  // console.log(props.origNote, props.origNote.ID)
  if (props.origNote && props.origNote.ID) {
    note.value.ID = props.origNote.ID
    note.value.Title = props.origNote.Title
    note.value.Description = props.origNote.Description
    note.value.Content = props.origNote.Content
    note.value.Syntax = props.origNote.Syntax
    note.value.Indent = props.origNote.Indent
  } else {
    note.value.ID = undefined
    note.value.Title = ''
    note.value.Description = ''
    note.value.Content = ''
    note.value.Syntax = DEFAULT_SYNTAX
    note.value.Indent = DEFAULT_INDENT
  }
}

const emitSave = () => {
  emits("note:save", note.value)
}

const emitClose = () => {
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
      New snippet
    </div>
    <div class="card-body">
      <div class="form mb-1">
        <div class="row mb-1">
          <label for="note-title" class="col-sm-3 col-form-label" aria-required="true">Title</label>
          <div class="col-sm-9">
            <input type="text" class="form-control" id="note-title"
              placeholder="Snippet title" required v-model="note.Title">
          </div>
        </div>

        <div class="row mb-1">
          <label for="note-description" class="col-sm-3 col-form-label">Description</label>
          <div class="col-sm-9">
            <input type="text" class="form-control" id="note-description"
              placeholder="Snippet description" v-model="note.Description">
          </div>
        </div>
        
        <!-- <div class="row">
          <label for="note-tag" class="col-sm-3 col-form-label" aria-required="true">Tag</label>
          <div class="col-sm-9">
            <input type="text" class="form-control" id="note-tag"
              placeholder="Snippet tag" required>
          </div>
        </div> -->
      </div>

      <!-- body -->
      <div class="card">
        <div class="card-header">
          <div class="row">
            <div class="col-6 col-sm-4 col-md-3 col-lg-4 col-xl-2">
              <!-- Syntax -->
              <select class="form-select" aria-label="Snippet Syntax" v-model="note.Syntax">
                <option selected value="auto">Auto</option>
                <option value="mysql">MySQL</option>
                <option value="shell">Shell</option>
                <option value="go">GoLang</option>
              </select>
            </div>

            <div class="col-6 col-sm-4 col-md-3 col-lg-3 col-xl-2">
              <!-- Indentation -->
              <select class="form-select" aria-label="Code indentation" v-model="note.Indent">
                <option value="2">2</option>
                <option value="4" selected>4</option>
                <option value="8">8</option>
              </select>
            </div>
          </div>
        </div>
        <div class="card-body" id="note-content">
          <textarea name="note-content" id="note-content" class="form-control"
            rows="8" placeholder="Paste a snippet of code..." v-model="note.Content"></textarea>
        </div>
        <div class="card-footer">
          <div class="d-grid gap-2 d-md-flex justify-content-md-end">
            <button type="button" class="btn btn-sm btn-danger" @click="emitClose">
              Cancel
            </button>
            <button type="button" class="btn btn-sm btn-success float-right" @click="emitSave">
              Create
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
</style>
