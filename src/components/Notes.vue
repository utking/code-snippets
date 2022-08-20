<script setup>
import { inject, onMounted, ref } from 'vue'
import TagList from './TagList.vue'
import NoteList from './NoteList.vue'
import NoteView from './NoteView.vue'
import NoteEditor from './NoteEditor.vue'

function injectStrict(key, fallback) {
  const resolved = inject(key, fallback)
  if (!resolved) {
    throw new Error(`Could not resolve ${key.description}`)
  }
  return resolved
}

const UNTAGGED = '[untagged]'
const http = injectStrict('axios')
const tags = ref([])
const notes = ref([])
const tagId = ref(UNTAGGED)
const curTag = ref({})
const curNote = ref({})
const noteMode = ref('create')
const error = ref(null)
const editTag = ref(false)

const loadTags = async () => {
  try {
    const resp = await http.get('/tag')
    tags.value = resp.data
    setActiveTag(tagId.value)
  } catch (err) {
    error.value = (err.response.data && err.response.data.Error) ? err.response.data.Error : err
  }
}

const updateTagAlias = async () => {
  error.value = ''
  try {
    await http.put(`/tag/${tagId.value}`, {ID: curTag.value.ID, Alias: curTag.value.Alias})
    editTag.value = false
    loadTags()
  } catch (err) {
    error.value = (err.response && err.response.data && err.response.data.Error) ? err.response.data.Error : err
  }
}

const createNote = async (data) => {
  error.value = ''
  let resp

  try {
    if (data.ID) {
      resp = await http.put(`/note/${data.ID}`, data)
    } else {
      resp = await http.post(`/note/`, data)
    }
    curNote.value = resp.data
    loadTags()
    loadNotes(tagId)
    setActiveNote(curNote.value.ID)
    noteMode.value = 'view'
    
  } catch (err) {
    error.value = (err.response && err.response.data && err.response.data.Error) ? err.response.data.Error : err
  }
}

const deleteNote = async (id) => {
  if (confirm("Are you sure you want to delete the snippet?")) {
    const resp = await http.delete(`/note/${id}`)
    try {
      loadTags()
      loadNotes(tagId)
      noteMode.value = 'create'
      error.value = ''
    } catch (err) {
      error.value = (err.response.data && err.response.data.Error) ? err.response.data.Error : err
    }
  }
}

const closeNoteEditor = () => {
  noteMode.value = (curNote.value.ID ? 'view' : 'create')
}

const toggleNoteEditor = async (id) => {
  if (id) {
    loadNote(id)
  }
  noteMode.value = 'edit'
}

const loadNotes = async (tagID) => {
  let resp
  try {
    if (tagID) {
      resp = await http.get(`/note/tag/${tagID}`)
      notes.value = resp.data
    } else {
      resp = await http.get('/note')
      notes.value = []
    }
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
  noteMode.value = 'create'
  error.value = ''

  tags.value.forEach((el) => {
    el.active = el.Alias === id
    if (el.active) {
      curTag.value = el
    }
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
  notes.value.forEach((el) => {
    el.active = el.ID === id
    if (el.active) {
      noteMode.value = 'view'
    }
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
        <h5 class="title h5 mt-2">Tags</h5>
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
        <h5 class="title h5 mt-2" v-if="!editTag">
          {{curTag.Alias}}
          <span title="Edit tag" @click="editTag = true" v-if="curTag.Alias">&#9998;</span>
        </h5>
        <div v-else class="input-group my-2">
          <input class="form-control" type="text" v-model="curTag.Alias">
          <button class="btn btn-outline-secondary" @click="updateTagAlias" title="Save">
            &check;
          </button>
          <button class="btn btn-outline-secondary" @click="editTag = false" title="Cancel">
            &cross;
          </button>
        </div>
        <article>
          <NoteList
            :notes="notes"

            @note:clicked="setActiveNote" />
        </article>
      </section>
    </div>
    <div class="col-sm-12 col-md-12 col-lg-6 col-xl-8">
      <section>
        <article v-if="noteMode === 'create'">
          <NoteEditor
            :title="curTag.Alias"
            
            @note:close="closeNoteEditor"
            @note:save="createNote" />
        </article>
        <article v-else-if="noteMode === 'edit'">
          <NoteEditor
            :orig-note="curNote"
            :title="curTag.Alias"
            
            @note:close="closeNoteEditor"
            @note:save="createNote" />
        </article>
        <article v-else>
          <NoteView
            :id="curNote.ID"
            :title="curNote.Title"
            :content="curNote.Content"
            :indent="curNote.Indent"
            :tag-alias="curTag.Alias"

            @note:delete="deleteNote"
            @note:edit="toggleNoteEditor" />
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
  max-height: 90vh;
}
</style>
