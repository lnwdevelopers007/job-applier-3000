import type { JobDisplay } from '$lib/types';
import JobCell from '../cells/JobCell.svelte';
import ActionsCell from '../cells/ActionsCell.svelte';
import StatusBadge from '$lib/components/ui/StatusBadge.svelte';

export type JobColumnOptions = {
	showCompanyColumn: boolean;
	onView: (job: JobDisplay) => void;
	onEdit?: (job: JobDisplay) => void;
	onDelete: (job: JobDisplay) => void;
};

export function createJobColumns(options: JobColumnOptions) {
	const { showCompanyColumn, onView, onEdit, onDelete } = options;

	return [
		...(showCompanyColumn
			? [
					{
						id: 'job',
						accessorKey: 'title' as keyof JobDisplay,
						header: 'Job',
						cell: (job: JobDisplay) => ({
							component: JobCell,
							props: {
								title: job.title,
								company: job.company,
								companyLogo: job.companyLogo
							}
						}),
						width: '400px',
						cellClass: 'px-6 py-2.5 whitespace-nowrap'
					}
				]
			: [
					{
						id: 'title',
						accessorKey: 'title' as keyof JobDisplay,
						header: 'Job Title',
						width: '400px',
						cellClass: 'px-6 py-4 whitespace-nowrap'
					}
				]),
		{
			id: 'location',
			accessorKey: 'location' as keyof JobDisplay,
			header: 'Location',
			width: '200px',
			cellClass: 'px-6 py-4 whitespace-nowrap'
		},
		{
			id: 'posted',
			accessorKey: 'posted' as keyof JobDisplay,
			header: 'Posted',
			width: '150px',
			cellClass: 'px-6 py-4 whitespace-nowrap'
		},
		{
			id: 'expires',
			accessorKey: 'expires' as keyof JobDisplay,
			header: 'Deadline',
			width: '150px',
			cellClass: 'px-6 py-4 whitespace-nowrap'
		},
		{
			id: 'applicants',
			accessorKey: 'applicants' as keyof JobDisplay,
			header: 'Applied',
			width: '100px',
			cellClass: 'px-6 py-4 whitespace-nowrap'
		},
		{
			id: 'status',
			accessorKey: 'status' as keyof JobDisplay,
			header: 'Status',
			cell: (job: JobDisplay) => ({
				component: StatusBadge,
				props: { status: job.status, type: 'status' }
			}),
			width: '120px',
			cellClass: 'px-6 py-4 whitespace-nowrap'
		},
		{
			id: 'actions',
			header: '',
			cell: (job: JobDisplay) => ({
				component: ActionsCell,
				props: {
					onView: () => onView(job),
					onEdit: onEdit ? () => onEdit(job) : undefined,
					onDelete: () => onDelete(job)
				}
			}),
			width: '80px',
			cellClass: 'px-6 py-4 whitespace-nowrap text-center relative'
		}
	];
}