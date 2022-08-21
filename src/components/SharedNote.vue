<script setup>
import { inject, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router';

const http = injectStrict('axios')
const copied = ref(false)
const note = ref({})
const error = ref(null)
const route = useRoute();

function injectStrict(key, fallback) {
  const resolved = inject(key, fallback)
  if (!resolved) {
    throw new Error(`Could not resolve ${key.description}`)
  }
  return resolved
}


onMounted(async () => {
  const hash = route.params.hash
  if (hash) {
    try {
      const resp = await http.get(`/note/s/${hash}`)
      note.value = resp.data
      error.value = ''
    } catch (err) {
      console.log(err)
      error.value = (err.response.data && err.response.data.Error) ? err.response.data.Error : err
    }
  }
})

const copyToClipboard = () => {
  if (!navigator.clipboard) {
    console.log('Clipboard API is not supported')
    return
  }
  navigator.clipboard.writeText(note.value.Content)
  copied.value =true
}

</script>

<template>
  <div v-if="error">
    <p class="alert alert-danger">{{ error }}</p>
  </div>
  <div class="card">
    <div class="card-body">
      <div>
        <span class="title h5">Shared | {{note.Title}}</span> <button class="btn btn-light" title="Copy" @click="copyToClipboard">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-clipboard" viewBox="0 0 16 16">
            <path d="M4 1.5H3a2 2 0 0 0-2 2V14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V3.5a2 2 0 0 0-2-2h-1v1h1a1 1 0 0 1 1 1V14a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V3.5a1 1 0 0 1 1-1h1v-1z"/>
            <path d="M9.5 1a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5h-3a.5.5 0 0 1-.5-.5v-1a.5.5 0 0 1 .5-.5h3zm-3-1A1.5 1.5 0 0 0 5 1.5v1A1.5 1.5 0 0 0 6.5 4h3A1.5 1.5 0 0 0 11 2.5v-1A1.5 1.5 0 0 0 9.5 0h-3z"/>
          </svg></button>
        <span class="show mx-2" role="alert" v-if="copied">
          Copied!
            <button type="button" class="btn-close" @click="copied = false" aria-label="Close"></button>
        </span>
      </div>
      
      <!-- body -->
      <div class="card">
        <div class="card-body">
          <pre v-highlightjs><code id="snippet-content" class="automatic">{{ note.Content }}</code></pre>
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
