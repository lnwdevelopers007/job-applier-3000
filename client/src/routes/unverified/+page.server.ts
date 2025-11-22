import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ url, request, cookies }) => {
	const userName = url.searchParams.get('name') || '';
	
	// Check if this is coming from callback (has referer) or if user has valid tokens (signup flow)
	const referer = request.headers.get('referer');
	const isFromCallback = referer && referer.includes('/callback');
	const hasValidTokens = cookies.get('access_token') && cookies.get('refresh_token');
	
	// Only clear cookies if this is NOT from the callback AND user doesn't have valid tokens
	// This preserves tokens during signup flow
	if (!isFromCallback && !hasValidTokens) {
		cookies.delete('access_token', { path: '/' });
		cookies.delete('refresh_token', { path: '/' });
	}
	
	// If not from callback and no name, redirect to login
	if (!isFromCallback && !userName) {
		throw redirect(303, '/login');
	}
	
	return {
		userName: decodeURIComponent(userName)
	};
};