<script lang="ts">
  import { FileText, Eye, Download, Trash2 } from 'lucide-svelte';
  import { fileService } from '$lib/services/fileService';
  
  interface Props {
    file: {
      id: string;
      filename: string;
      category: string;
      size: number;
      uploadDate: string;
    };
    onPreview?: (event: MouseEvent) => void;
    onDelete?: (event: MouseEvent) => void;
  }

  let { file, onPreview = () => {}, onDelete = () => {} }: Props = $props();
  
  async function handleDownload() {
    try {
      const blob = await fileService.downloadFile(file.id);
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = file.filename;
      a.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      console.error('Download error:', err);
    }
  }
  
  const categoryColor = $derived(fileService.getCategoryColor(file.category));
  const categoryLabel = $derived(fileService.getCategoryLabel(file.category));
  const formattedSize = $derived(fileService.formatFileSize(file.size));
  const formattedDate = $derived(new Date(file.uploadDate).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  }));
</script>

<div class="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors group">
  <div class="flex items-center gap-4 flex-1 min-w-0">
    <!-- Icon -->
    <div class="w-12 h-12 bg-{categoryColor}-100 rounded-lg flex items-center justify-center flex-shrink-0">
      <FileText class="w-6 h-6 text-{categoryColor}-600" />
    </div>
    
    <!-- Info -->
    <div class="flex-1 min-w-0">
      <div class="flex items-center gap-2 mb-1">
        <p class="text-sm font-medium text-gray-900 truncate">{file.filename}</p>
        <span class="px-2 py-0.5 text-xs font-medium bg-{categoryColor}-100 text-{categoryColor}-700 rounded">
          {categoryLabel}
        </span>
      </div>
      <p class="text-xs text-gray-500">{formattedSize} • Uploaded on {formattedDate}</p>
    </div>
  </div>
  
  <!-- Actions -->
  <div class="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
    <button
      onclick={onPreview}
      class="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
      title="Preview"
      aria-label="Preview document"
    >
      <Eye class="w-5 h-5" />
    </button>
    <button
      onclick={handleDownload}
      class="p-2 text-gray-400 hover:text-green-600 hover:bg-green-50 rounded-lg transition-colors"
      title="Download"
      aria-label="Download document"
    >
      <Download class="w-5 h-5" />
    </button>
    <button
      onclick={onDelete}
      class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
      title="Delete"
      aria-label="Delete document"
    >
      <Trash2 class="w-5 h-5" />
    </button>
  </div>
</div>