<script setup>
import { inject, onBeforeMount, ref } from 'vue'
import TagList from './TagList.vue';
import NoteList from './NoteList.vue';
import NoteContent from './NoteContent.vue';
import { computed } from '@vue/reactivity';

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

const loadTags = async () => {
  const resp = await http.get('/tag')
  tags.value = resp.data
}

const loadNotes = async (categoryId) => {
  let resp
  if (categoryId) {
    resp = await http.get(`/note/category/${categoryId}`)
  } else {
    resp = await http.get('/note')
  }
  notes.value = resp.data
}

const loadNote = async (id) => {
  if (id) {
    const resp = await http.get(`/note/${id}`)
    curNote.value = resp.data
  }
}

const getTagStats = () => {
  loadTags()
}

const setActiveTag = (id) => {
  tags.value.forEach((el) => {
    el.active = el.ID === id
  })
  tagId.value = id
  if (id) {
    loadNotes(id)
  } else {
    notes.value = []
  }
  setActiveNote(undefined)
  curNote.value = {}
}

const setActiveNote = (id) => {
  notes.value.forEach((el) => {
    el.active = el.ID === id
  })
  loadNote(id)
}

onBeforeMount(getTagStats)

</script>

<template>
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
        <header class="h5">
          Note Content
        </header>
        <article>
          <NoteContent
            :tag="curNote.TagID??''"
            :title="curNote.Title"
            :description="curNote.Description"
            :content="curNote.Content"
            :indent="curNote.Indent"
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
