import ActionsCell from '../cells/ActionsCell.svelte';
import ClickableNumberCell from '../cells/ClickableNumberCell.svelte';
import StatusBadge from '$lib/components/ui/StatusBadge.svelte';

export interface CompanyJobDisplay {
	id: string;
	title: string;
	status: string;
	applicants: number;
	posted: string;
	expires: string;
}

export type CompanyJobColumnOptions = {
	onEdit: (job: CompanyJobDisplay) => void;
	onView: (job: CompanyJobDisplay) => void;
	onDelete: (job: CompanyJobDisplay) => void;
	onViewApplicants: (job: CompanyJobDisplay) => void;
};

export function createCompanyJobColumns(options: CompanyJobColumnOptions) {
	const { onEdit, onView, onDelete, onViewApplicants } = options;

	return [
		{
			id: 'title',
			accessorKey: 'title' as keyof CompanyJobDisplay,
			header: 'Job Title',
			cell: (job: CompanyJobDisplay) => job.title,
			width: '250px',
			cellClass: 'px-6 py-4 whitespace-nowrap font-medium text-gray-900'
		},
		{
			id: 'applicants',
			accessorKey: 'applicants' as keyof CompanyJobDisplay,
			header: 'Applicants',
			cell: (job: CompanyJobDisplay) => ({
				component: ClickableNumberCell,
				props: {
					value: job.applicants,
					onClick: () => onViewApplicants(job)
				}
			}),
			width: '100px',
			cellClass: 'px-6 py-4 whitespace-nowrap text-center'
		},
		{
			id: 'posted',
			accessorKey: 'posted' as keyof CompanyJobDisplay,
			header: 'Posted',
			cell: (job: CompanyJobDisplay) => job.posted,
			width: '120px',
			cellClass: 'px-6 py-4 whitespace-nowrap'
		},
		{
			id: 'expires',
			accessorKey: 'expires' as keyof CompanyJobDisplay,
			header: 'Expires',
			cell: (job: CompanyJobDisplay) => job.expires,
			width: '120px',
			cellClass: 'px-6 py-4 whitespace-nowrap'
		},
		{
			id: 'status',
			accessorKey: 'status' as keyof CompanyJobDisplay,
			header: 'Status',
			cell: (job: CompanyJobDisplay) => ({
				component: StatusBadge,
				props: { 
					status: job.status, 
					type: 'job_status',
					rawStatus: job.status.toLowerCase()
				}
			}),
			width: '120px',
			cellClass: 'px-6 py-4 whitespace-nowrap'
		},
		{
			id: 'actions',
			header: '',
			cell: (job: CompanyJobDisplay) => ({
				component: ActionsCell,
				props: {
					onView: () => onView(job),
					onEdit: () => onEdit(job),
					onDelete: () => onDelete(job),
					customActions: []
				}
			}),
			width: '80px',
			cellClass: 'px-6 py-4 whitespace-nowrap text-center relative'
		}
	];
}