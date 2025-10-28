<script>
  import DOMPurify from "dompurify";
  export let html = "";

  $: safeHTML = DOMPurify.sanitize(html, { USE_PROFILES: { html: true } });
</script>

<!-- eslint-disable-next-line svelte/no-at-html-tags -->
<div class="job-description">{@html safeHTML}</div>

<style>
  .job-description :global(h2),
  .job-description :global(h3) {
    font-weight: 500;
    margin-top: 1.5rem;
    margin-bottom: 0.75rem;
    color: #1f2937;
  }

  .job-description :global(h2) {
    font-size: 1.125rem;
  }

  .job-description :global(strong),
  .job-description :global(b) {
    font-weight: 600;
    color: #374151;
  }

</style>
