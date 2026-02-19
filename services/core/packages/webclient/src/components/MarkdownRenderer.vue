<script setup>
import {computed} from 'vue'
import MarkdownIt from 'markdown-it'
import DOMPurify from 'dompurify'

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true
})

const props = defineProps({
  markdown: {
    type: String,
    default: ''
  }
})

const renderedHtml = computed(() => {
  const unsafeHtml = md.render(props.markdown)
  return DOMPurify.sanitize(unsafeHtml)
})
</script>

<template>
  <div class="markdown-content" v-html="renderedHtml"></div>
</template>

<style scoped>

::v-deep * {
  padding: revert;
  margin: revert;
}

::v-deep p {
  text-align: justify;
}

::v-deep a {
  color: #ffddb7 !important;
  text-decoration: underline;
}

::v-deep code, ::v-deep pre {
  background-color: rgba(104, 104, 104, 0.2);
  border-radius: 8px;
  font-family: 'Courier New', Courier, monospace;
  overflow-x: auto;
}

::v-deep code {
  padding: 2px 4px;
  border: 1px solid rgba(201, 189, 166, 0.2);
  font-size: 0.9em;
}

::v-deep pre code {
  padding: 0;
  border: none;
}

::v-deep pre {
  padding: 8px;
}

::v-deep pre code {
  background-color: transparent;
}

::v-deep ul, ::v-deep ol{
  padding-left: 1.5em;
}

::v-deep li {
  margin: 0.5em 0;
}

::v-deep img {
  max-width: 100%;
  height: auto;
  display: block;
  margin: 1em auto;
}

::v-deep table {
  width: 100%;
  border-collapse: collapse;
  margin: 1em 0;
}

::v-deep th, ::v-deep td {
  border: 1px solid rgba(201, 189, 166, 0.5);
  padding: 8px;
  text-align: left;
}

::v-deep hr {
  border: none;
  border-top: 1px solid rgba(201, 189, 166, 0.5);
  margin: 3em 0;
}

::v-deep blockquote {
  border-left: 4px solid #e69469;
  border-radius: 8px;
  background: rgba(255, 220, 167, 0.1);
  padding: 1em;
  margin: 1em 0;
}

::v-deep blockquote p {
  margin: 0;
  text-align: start;
}

</style>
