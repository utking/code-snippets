<script setup>
import { inject, onBeforeMount, onMounted, ref } from 'vue'
import TagList from './TagList.vue';
import NoteList from './NoteList.vue';
import { computed } from '@vue/reactivity';
import NoteView from './NoteView.vue';
import NoteEditor from './NoteEditor.vue';

function injectStrict(key, fallback) {
  const resolved = inject(key, fallback)
  if (!resolved) {
    throw new Error(`Could not resolve ${key.description}`)
  }
  return resolved
}

const http = injectStrict('axios')
const tags = ref([])
const notes = ref([])
const tagId = ref(0)
const curNote = ref({})
const editNote = ref(false)
const error = ref(null)

const loadTags = async () => {
  try {
    const resp = await http.get('/tag')
    tags.value = resp.data
    setActiveTag(tagId.value)
  } catch (err) {
    error.value = (err.response.data && err.response.data.Error) ? err.response.data.Error : err
  }
}

const createNote = async (data) => {
  data.TagID = tagId
  data.Indent = parseInt(data.Indent)
  let resp
  try {
    if (data.ID) {
      resp = await http.put(`/note/${data.ID}`, data)
    } else {
      resp = await http.post(`/note/`, data)
    }
    loadNotes(tagId.value)
    error.value = ''
  } catch (err) {
    error.value = (err.response.data && err.response.data.Error) ? err.response.data.Error : err
  }
}

const deleteNote = async (id) => {
  const resp = await http.delete(`/note/${id}`)
  try {
    loadNotes(tagId.value)
    closeNoteEditor()
    error.value = ''
  } catch (err) {
    error.value = (err.response.data && err.response.data.Error) ? err.response.data.Error : err
  }
}

const closeNoteEditor = () => {
  editNote.value = false
  curNote.value = {ID: 0}
}

const toggleNoteEditor = async (id) => {
  if (id) {
    loadNote(id)
  }
  editNote.value = !!id
}

const loadNotes = async (tagID) => {
  let resp
  try {
    if (tagID) {
      resp = await http.get(`/note/category/${tagID}`)
    } else {
      resp = await http.get('/note')
    }
    notes.value = resp.data
    error.value = ''
    
  } catch (err) {
    error.value = (err.response.data && err.response.data.Error) ? err.response.data.Error : err
  }
}

const loadNote = async (id) => {
  if (id) {
    try {
      const resp = await http.get(`/note/${id}`)
      curNote.value = resp.data
      error.value = ''
    } catch (err) {
      error.value = (err.response.data && err.response.data.Error) ? err.response.data.Error : err
    }
  }
}

const setActiveTag = (id) => {
  // curNote.value = 0
  closeNoteEditor()
  tags.value.forEach((el) => {
    el.active = el.ID === id
  })
  tagId.value = id
  if (id) {
    loadNotes(id)
  } else {
    notes.value = []
  }
  setActiveNote()
  curNote.value = {}
}

const setActiveNote = (id) => {
  closeNoteEditor()
  notes.value.forEach((el) => {
    el.active = el.ID === id
  })
  loadNote(id)
}

onMounted(loadTags)

</script>

<template>
  <div v-if="error">
    <p class="alert alert-danger">{{ error }}</p>
  </div>
  <div class="row">
    <div class="col-sm-12 col-md-4 col-lg-3 col-xl-2 scrollable">
      <section>
        <header class="h5">
          Tags
        </header>
        <article>
          <TagList
            :tags="tags"
            @tag:clicked="setActiveTag"
            />
        </article>
      </section>
    </div>
    <div class="col-sm-12 col-md-8 col-lg-3 col-xl-2 scrollable">
      <section>
        <header class="h5">
          Notes
        </header>
        <article>
          <NoteList
            :notes="notes"

            @note:clicked="setActiveNote" />
        </article>
      </section>
    </div>
    <div class="col-sm-12 col-md-12 col-lg-6 col-xl-8">
      <section>
        <article v-if="editNote || !curNote.ID">
          <NoteEditor
            :orig-note="curNote"
            :tag-id="tagId"
            
            @note:close="closeNoteEditor"
            @note:save="createNote" />
        </article>
        <article v-else>
          <NoteView
            :id="curNote.ID"
            :title="curNote.Title"
            :description="curNote.Description"
            :content="curNote.Content"
            :syntax="curNote.Syntax"
            :indent="curNote.Indent"

            @note:delete="deleteNote"
            @note:edit="toggleNoteEditor"
          />
        </article>
      </section>
    </div>
  </div>
</template>

<style scoped>
.row > div.scrollable {
  overflow: auto;
  border-radius: 3px;
  border: 1px solid #f0f0f0;
}

.row > div {
  min-height: 50vh;
  max-height: 100vh;
}
</style>
