<script lang="ts">
  import Badge from '../ui/Badge.svelte';
  
  type BadgeVariant = 'primary' | 'secondary' | 'warning' | 'danger' | 'success' | 'info' | 'purple' | 'indigo' | 'pink' | 'teal' | 'orange' | 'cyan' | 'lime' | 'rose';

  let {
    type,
    value,
    size = 'sm'
  }: {
    type: 'workType' | 'workArrangement';
    value: string;
    size?: 'xs' | 'sm' | 'md';
  } = $props();

  // Get variant based on work type (based on actual schema options)
  function getWorkTypeVariant(workType: string): BadgeVariant {
    switch (workType.toLowerCase()) {
      case 'full-time':
      case 'fulltime':
        return 'success';   // Green
      case 'part-time':
      case 'parttime':
        return 'info';      // Sky blue
      case 'contract':
        return 'orange';    // Orange
      case 'casual':
        return 'purple';    // Purple
      default:
        return 'secondary'; // Gray
    }
  }

  // Get variant based on work arrangement (based on actual schema options)
  function getWorkArrangementVariant(arrangement: string): BadgeVariant {
    switch (arrangement.toLowerCase()) {
      case 'remote':
        return 'teal';      // Teal (different from work types)
      case 'on-site':
      case 'onsite':
        return 'indigo';    // Indigo (different from work types)
      case 'hybrid':
        return 'pink';      // Pink (different from work types)
      default:
        return 'secondary'; // Gray
    }
  }

  const variant = $derived(
    type === 'workType' 
      ? getWorkTypeVariant(value)
      : getWorkArrangementVariant(value)
  );
</script>

<Badge variant={variant as any} {size} text={value} />