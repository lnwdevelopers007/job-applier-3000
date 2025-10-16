import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ cookies }) => {
	// Clear the auth cookies
	cookies.delete('access_token', { path: '/' });
	cookies.delete('refresh_token', { path: '/' });

	// Optional: Call your backend logout endpoint to invalidate tokens
	// You could also make a request to your Go backend's logout endpoint here
	
	return json({ success: true });
};