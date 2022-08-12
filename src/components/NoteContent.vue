<script setup>

defineProps({
  tag: String,
  description: String,
  alias: String,
  syntax: String,
  content: String,
  indent: Number,
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
              placeholder="Snippet title" required>
          </div>
        </div>

        <div class="row mb-1">
          <label for="note-description" class="col-sm-3 col-form-label">Description</label>
          <div class="col-sm-9">
            <input type="text" class="form-control" id="note-description"
              placeholder="Snippet description">
          </div>
        </div>
        
        <div class="row">
          <label for="note-tag" class="col-sm-3 col-form-label" aria-required="true">Tag</label>
          <div class="col-sm-9">
            <input type="text" class="form-control" id="note-tag"
              placeholder="Snippet tag" required>
          </div>
        </div>
      </div>

      <!-- body -->
      <div class="card">
        <div class="card-header">
          <div class="row">
            <div class="col-6 col-sm-4 col-md-3 col-lg-4 col-xl-2">
              <!-- Syntax -->
              <select class="form-select" aria-label="Snippet Syntax">
                <option selected value="auto">Auto</option>
                <option value="mysql">MySQL</option>
                <option value="shell">Shell</option>
                <option value="go">GoLang</option>
              </select>
            </div>

            <div class="col-6 col-sm-4 col-md-3 col-lg-3 col-xl-2">
              <!-- Indentation -->
              <select class="form-select" aria-label="Code indentation">
                <option value="2">2</option>
                <option value="4" selected>4</option>
                <option value="8">8</option>
              </select>
            </div>
          </div>
        </div>
        <div class="card-body">
          <textarea name="note-content" id="note-content" class="form-control"
            rows="8" placeholder="Paste a snippet of code..."></textarea>
        </div>
        <div class="card-footer">
          <div class="d-grid gap-2 d-md-flex justify-content-md-end">
            <button type="button" class="btn btn-sm btn-danger">Cancel</button>
            <button type="button" class="btn btn-sm btn-success float-right">Create</button>
          </div>
        </div>
      </div>
      <!-- end of body -->

    </div>
  </div>
</template>

<style scoped>
</style>
