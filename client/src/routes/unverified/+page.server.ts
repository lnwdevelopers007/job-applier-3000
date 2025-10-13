import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ url, request, cookies }) => {
	const userName = url.searchParams.get('name') || '';
	
	// Check if this is coming from callback (has referer) or direct access
	const referer = request.headers.get('referer');
	const isFromCallback = referer && referer.includes('/callback');
	
	// Clear any existing authentication cookies since user is unverified
	cookies.delete('access_token', { path: '/' });
	cookies.delete('refresh_token', { path: '/' });
	
	// If not from callback and no name, redirect to login
	if (!isFromCallback && !userName) {
		throw redirect(303, '/login');
	}
	
	return {
		userName: decodeURIComponent(userName)
	};
};