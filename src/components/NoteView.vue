<script setup>
import { ref } from "@vue/reactivity"


const props = defineProps({
  id: Number,
  title: String,
  content: String,
  indent: Number,
  tagAlias: String
})

const copied = ref(false)

const emits = defineEmits(['note:edit', 'note:delete'])

const emitEdit = () => {
  emits("note:edit", props.id)
}

const emitDelete = () => {
  emits("note:delete", props.id)
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
          &#x2752;</button>
        <span class="fade show mx-2" role="alert" v-if="copied">
          Copied!
          <button type="button" class="btn-close btn-sm" @click="copied = false"
            data-bs-dismiss="alert" aria-label="Close"></button>
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
