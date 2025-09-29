<script>
  import DOMPurify from "dompurify";
  export let html = "";

  $: safeHTML = DOMPurify.sanitize(html, { USE_PROFILES: { html: true } });
</script>
<!-- eslint-disable-next-line svelte/no-at-html-tags -->
<div>{@html safeHTML}</div>
