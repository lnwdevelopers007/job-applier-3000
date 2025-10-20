import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals }) => {
	// Server-side authentication check
	if (!locals.user?.isAuthenticated) {
		throw redirect(303, '/login?returnUrl=/app/applications');
	}

	return {
		user: locals.user
	};
};