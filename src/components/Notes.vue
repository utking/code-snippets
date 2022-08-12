<script setup>
import { onBeforeMount, ref } from 'vue'
import TagList from './TagList.vue';
import NoteList from './NoteList.vue';
import NoteContent from './NoteContent.vue';
import { computed } from '@vue/reactivity';

defineProps({
  msg: String
})

const tags = ref([
  {id: undefined, alias: "[untagged]", notes: 0},
  {id: 0, alias: "mysql", notes: 1},
  {id: 1, alias: "nginx", notes: 2},
  {id: 2, alias: "elasticsearch", notes: 1},
])

const notes = ref([
  {id: 0, tag_id: 1, title: 'Some note', active: 0},
  {id: 1, tag_id: 0, title: 'Some important SELECT', active: 0},
  {id: 2, tag_id: 2, title: 'Drop ES index', active: 0},
  {id: 3, tag_id: 2, title: 'List indices', active: 0},
])

const tagId = ref(undefined)

const getTagStats = () => {
  tags.value.forEach((tag) => {
    const items = notes.value.filter((el) => {
      return el.tag_id === tag.id
    })
    tag.notes = items.length
  })
}

const tagNotes = computed(() => {
  return notes.value.filter((el) => {
    return el.tag_id === tagId.value
  })
})

const setActiveTag = (id) => {
  tags.value.forEach((el) => {
    el.active = el.id === id
  })
  tagId.value = id
  setActiveNote(undefined)
}

const setActiveNote = (id) => {
  notes.value.forEach((el) => {
    el.active = el.id === id
  })
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
            :notes="tagNotes"
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
          <NoteContent />
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
