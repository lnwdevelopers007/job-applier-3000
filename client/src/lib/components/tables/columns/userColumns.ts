import type { User } from '$lib/types';
import UserCell from '../cells/UserCell.svelte';
import ActionsCell from '../cells/ActionsCell.svelte';
import DocumentsCell from '../cells/DocumentsCell.svelte';
import StatusBadge from '$lib/components/ui/StatusBadge.svelte';
import { formatDateCompact } from '$lib/utils/datetime';

export interface UserDisplay extends User {
	banned?: boolean;
	files?: string[];
}

export type UserColumnOptions = {
	onEdit: (user: UserDisplay) => void;
	onBan: (user: UserDisplay) => void;
	onDelete: (user: UserDisplay) => void;
	getCurrentUserId: () => string | undefined;
};

export function createUserColumns(options: UserColumnOptions) {
	const { onEdit, onBan, onDelete, getCurrentUserId } = options;

	return [
		{
			id: 'user',
			accessorKey: 'name' as keyof UserDisplay,
			header: 'User',
			cell: (user: UserDisplay) => ({
				component: UserCell,
				props: {
					name: user.name,
					email: user.email,
					avatarURL: user.avatarURL
				}
			}),
			width: '300px',
			cellClass: 'px-6 py-4 whitespace-nowrap'
		},
		{
			id: 'documents',
			header: 'Documents',
			cell: (user: UserDisplay) => ({
				component: DocumentsCell,
				props: {
					files: user.files || [],
					userName: user.name
				}
			}),
			width: '100px',
			cellClass: 'px-6 py-4 whitespace-nowrap'
		},
		{
			id: 'role',
			accessorKey: 'role' as keyof UserDisplay,
			header: 'Role',
			cell: (user: UserDisplay) => {
				const roleLabels: Record<string, string> = {
					jobSeeker: 'Job Seeker',
					company: 'Company',
					faculty: 'Faculty',
					admin: 'Admin'
				};
				return {
					component: StatusBadge,
					props: { 
						status: roleLabels[user.role] || user.role, 
						type: 'role',
						rawStatus: user.role 
					}
				};
			},
			width: '120px',
			cellClass: 'px-6 py-4 whitespace-nowrap'
		},
		{
			id: 'verified',
			accessorKey: 'verified' as keyof UserDisplay,
			header: 'Verified',
			cell: (user: UserDisplay) => ({
				component: StatusBadge,
				props: { 
					status: user.verified ? 'Verified' : 'Unverified', 
					type: 'verification',
					rawStatus: user.verified ? 'verified' : 'unverified'
				}
			}),
			width: '120px',
			cellClass: 'px-6 py-4 whitespace-nowrap'
		},
		{
			id: 'banned',
			accessorKey: 'banned' as keyof UserDisplay,
			header: 'Status',
			cell: (user: UserDisplay) => ({
				component: StatusBadge,
				props: { 
					status: user.banned ? 'Banned' : 'Active',
					type: 'user_status',
					rawStatus: user.banned ? 'banned' : 'active'
				}
			}),
			width: '120px',
			cellClass: 'px-6 py-4 whitespace-nowrap'
		},
		{
			id: 'createdAt',
			accessorKey: 'createdAt' as keyof UserDisplay,
			header: 'Joined',
			cell: (user: UserDisplay) => formatDateCompact(user.createdAt),
			width: '120px',
			cellClass: 'px-6 py-4 whitespace-nowrap'
		},
		{
			id: 'actions',
			header: '',
			cell: (user: UserDisplay) => {
				const currentUserId = getCurrentUserId();
				const isCurrentUser = currentUserId === user.id;
				
				// Don't allow users to ban themselves
				const customActions = !isCurrentUser ? [
					{
						label: user.banned ? 'Unban' : 'Ban',
						action: () => onBan(user),
						variant: 'warning' as const
					}
				] : [];
				
				return {
					component: ActionsCell,
					props: {
						onEdit: () => onEdit(user),
						onDelete: () => onDelete(user),
						customActions
					}
				};
			},
			width: '80px',
			cellClass: 'px-6 py-4 whitespace-nowrap text-center relative'
		}
	];
}