<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { Editor } from '@tiptap/core';
  import StarterKit from '@tiptap/starter-kit';
  import Placeholder from '@tiptap/extension-placeholder';
  import { Bold, Italic, List, ListOrdered, Heading2 } from 'lucide-svelte';
  
  let { value = $bindable(''), placeholder = '' } = $props();
  
  let element: HTMLDivElement;
  let editor = $state<Editor | null>(null);
  
  onMount(() => {
    editor = new Editor({
      element: element,
      extensions: [
        StarterKit,
        Placeholder.configure({
          placeholder: placeholder
        })
      ],
      content: value,
      onTransaction: () => {
        editor = editor;
      },
      onUpdate: ({ editor }) => {
        value = editor.getHTML();
      },
      editorProps: {
        attributes: {
          class: 'tiptap-editor min-h-[200px] px-3 py-2 focus:outline-none text-sm'
        }
      }
    });
  });
  
  onDestroy(() => {
    if (editor) {
      editor.destroy();
    }
  });
  
  $effect(() => {
    if (editor && value !== editor.getHTML()) {
      editor.commands.setContent(value);
    }
  });
</script>

<div class="border border-gray-300 rounded-md focus-within:ring-1 focus-within:ring-gray-500 focus-within:border-gray-400">
  <!-- Toolbar -->
  <div class="border-b border-gray-200 bg-gray-50 px-3 py-2 flex items-center gap-1 rounded-t-md">
    <button
      type="button"
      onclick={() => editor?.chain().focus().toggleBold().run()}
      class:active={editor?.isActive('bold')}
      class="p-1.5 rounded hover:bg-gray-100 transition-colors"
      disabled={!editor}
    >
      <Bold class="w-4 h-4" />
    </button>
    
    <button
      type="button"
      onclick={() => editor?.chain().focus().toggleItalic().run()}
      class:active={editor?.isActive('italic')}
      class="p-1.5 rounded hover:bg-gray-100 transition-colors"
      disabled={!editor}
    >
      <Italic class="w-4 h-4" />
    </button>
    
    <div class="w-px h-6 bg-gray-300 mx-1"></div>
    
    <button
      type="button"
      onclick={() => editor?.chain().focus().toggleHeading({ level: 2 }).run()}
      class:active={editor?.isActive('heading', { level: 2 })}
      class="p-1.5 rounded hover:bg-gray-100 transition-colors"
      disabled={!editor}
    >
      <Heading2 class="w-4 h-4" />
    </button>
    
    <div class="w-px h-6 bg-gray-300 mx-1"></div>
    
    <button
      type="button"
      onclick={() => editor?.chain().focus().toggleBulletList().run()}
      class:active={editor?.isActive('bulletList')}
      class="p-1.5 rounded hover:bg-gray-100 transition-colors"
      disabled={!editor}
    >
      <List class="w-4 h-4" />
    </button>
    
    <button
      type="button"
      onclick={() => editor?.chain().focus().toggleOrderedList().run()}
      class:active={editor?.isActive('orderedList')}
      class="p-1.5 rounded hover:bg-gray-100 transition-colors"
      disabled={!editor}
    >
      <ListOrdered class="w-4 h-4" />
    </button>
  </div>
  
  <!-- Editor -->
  <div bind:this={element}></div>
</div>

<style>
  button.active {
    background-color: rgb(229 231 235);
  }
  
  :global(.ProseMirror) {
    min-height: 200px;
  }
  
  :global(.ProseMirror p) {
    margin-top: 0.375rem;
    margin-bottom: 0.375rem;
    line-height: 1.375;
  }
  
  :global(.ProseMirror h2) {
    font-size: 1.125rem;
    font-weight: 600;
    margin-top: 0.75rem;
    margin-bottom: 0.75rem;
    line-height: 1.25;
  }
  
  :global(.ProseMirror ul) {
    margin-top: 0.5rem;
    margin-bottom: 0.5rem;
    padding-left: 1.5rem;
    line-height: 1.5;
    list-style: disc !important;
  }
  
  :global(.ProseMirror ol) {
    margin-top: 0.5rem;
    margin-bottom: 0.5rem;
    padding-left: 1.5rem;
    line-height: 1.5;
    list-style: decimal !important;
  }
  
  :global(.ProseMirror li) {
    margin-top: 0.25rem;
    margin-bottom: 0.25rem;
    display: list-item !important;
  }
  
  :global(.ProseMirror ul ul) {
    list-style: circle !important;
  }
  
  :global(.ProseMirror ul ul ul) {
    list-style: square !important;
  }
  
  :global(.ProseMirror p.is-editor-empty:first-child::before) {
    color: rgb(156 163 175);
    content: attr(data-placeholder);
    float: left;
    height: 0;
    pointer-events: none;
  }
  
  :global(.ProseMirror:focus) {
    outline: none;
  }
</style>