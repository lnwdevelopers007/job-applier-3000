import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ url, request }) => {
	const userName = url.searchParams.get('name') || '';
	
	// Simple logic: if no name and not from callback, redirect to login
	const referer = request.headers.get('referer');
	const isFromCallback = referer && referer.includes('/callback');
	
	if (!isFromCallback && !userName) {
		throw redirect(303, '/login');
	}
	
	return {
		userName: decodeURIComponent(userName)
	};
};