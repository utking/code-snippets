<script setup>
import { ref } from "@vue/reactivity"


const props = defineProps({
  id: Number,
  title: String,
  content: String,
  indent: Number,
  tagAlias: String,
  shareHash: String
})

const copied = ref(false)
const DAY = 86400000
const getDate = (days) => {
  const d = new Date(new Date().getTime() + days * DAY)
  return d.toISOString()
}

const emits = defineEmits(['note:edit', 'note:delete', 'note:share'])

const emitEdit = () => {
  emits("note:edit", props.id)
}

const emitDelete = () => {
  emits("note:delete", props.id)
}

const shareNote = () => {
  emits("note:share", {
    NoteID: props.id,
    ValidUntil: getDate(30)
  })
}

const copyToClipboard = () => {
  if (!navigator.clipboard) {
    console.log('Clipboard API is not supported')
    return
  }
  let copyText = document.querySelector("#snippet-content");
  navigator.clipboard.writeText(copyText.textContent)
  copied.value =true
}

</script>

<template>
  <div class="card">
    <div class="card-body">
      <div>
        <span class="title h5">{{tagAlias}} | {{title}}</span> <button class="btn btn-light" title="Copy" @click="copyToClipboard">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-clipboard" viewBox="0 0 16 16">
            <path d="M4 1.5H3a2 2 0 0 0-2 2V14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V3.5a2 2 0 0 0-2-2h-1v1h1a1 1 0 0 1 1 1V14a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V3.5a1 1 0 0 1 1-1h1v-1z"/>
            <path d="M9.5 1a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5h-3a.5.5 0 0 1-.5-.5v-1a.5.5 0 0 1 .5-.5h3zm-3-1A1.5 1.5 0 0 0 5 1.5v1A1.5 1.5 0 0 0 6.5 4h3A1.5 1.5 0 0 0 11 2.5v-1A1.5 1.5 0 0 0 9.5 0h-3z"/>
          </svg></button>
        <button class="btn btn-light" title="Share" @click="shareNote">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-share" viewBox="0 0 16 16">
            <path d="M13.5 1a1.5 1.5 0 1 0 0 3 1.5 1.5 0 0 0 0-3zM11 2.5a2.5 2.5 0 1 1 .603 1.628l-6.718 3.12a2.499 2.499 0 0 1 0 1.504l6.718 3.12a2.5 2.5 0 1 1-.488.876l-6.718-3.12a2.5 2.5 0 1 1 0-3.256l6.718-3.12A2.5 2.5 0 0 1 11 2.5zm-8.5 4a1.5 1.5 0 1 0 0 3 1.5 1.5 0 0 0 0-3zm11 5.5a1.5 1.5 0 1 0 0 3 1.5 1.5 0 0 0 0-3z"/>
          </svg>
        </button>
        <span class="show mx-2" role="alert" v-if="copied">
          Copied!
            <button type="button" class="btn-close" @click="copied = false" aria-label="Close"></button>
        </span>
        <span class="shared-url btn btn-link" v-if="shareHash">
          <router-link :to="'/note/fe/'+shareHash">Shared URL</router-link>
        </span>
      </div>
      
      <!-- body -->
      <div class="card">
        <div class="card-body">
          <pre v-highlightjs><code id="snippet-content" class="automatic">{{ content }}</code></pre>
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
