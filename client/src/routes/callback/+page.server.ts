import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { jwtDecode } from 'jwt-decode';

type JWTPayload = {
	email?: string;
	name?: string;
	avatarURL?: string;
	userID?: string;
	role?: string;
	verified?: boolean;
	exp?: number;
};

export const load: PageServerLoad = async ({ url }) => {
	const token = url.searchParams.get('token');
	const error = url.searchParams.get('error');

	// Handle error from OAuth
	if (error) {
		throw redirect(303, `/login?error=${encodeURIComponent(error)}`);
	}

	// No token provided
	if (!token) {
		throw redirect(303, '/login?error=No+authentication+token+received');
	}

	// Decode and validate token
	let decoded: JWTPayload | null = null;
	try {
		decoded = jwtDecode<JWTPayload>(token);
	} catch (error) {
		console.error('Failed to decode JWT in callback:', error);
		throw redirect(303, '/login?error=Invalid+authentication+token');
	}

	if (!decoded) {
		throw redirect(303, '/login?error=Invalid+authentication+token');
	}

	// Check if token is expired
	if (decoded.exp && decoded.exp * 1000 < Date.now()) {
		throw redirect(303, '/login?error=Authentication+token+expired');
	}

	// Check if this is part of signup flow
	const step = url.searchParams.get('step');
	const role = decoded.role?.toLowerCase();
	
	// Skip verification check if user is in signup flow
	if (!decoded.verified && step !== 'signup') {
		const userName = decoded.name ? encodeURIComponent(decoded.name) : '';
		throw redirect(303, `/unverified?name=${userName}`);
	}

	// Backend already sets the cookies properly during OAuth callback
	// Frontend should NOT overwrite them with the URL token parameter
	// The token in URL is just for validation, not for setting cookies

	// Determine redirect path based on role
	let redirectPath = '/';
	if (role === 'company') {
		if (step === 'signup') {
			redirectPath = '/signup/company?currentStep=2';
		} else {
			redirectPath = '/company/dashboard';
		}
	} else if (role === 'faculty') {
		redirectPath = '/app/jobs';
	} else if (role === 'admin') {
		redirectPath = '/admin/jobs';
	} else if (role === 'jobseeker') {
		if (step === 'signup') {
			redirectPath = '/signup/student?currentStep=2';
		} else {
			redirectPath = '/app/jobs';
		}
	}

	// Return data to the page component to handle redirect - server redirects lose cookies
	return {
		redirectPath,
		token,
		userInfo: {
			email: decoded.email,
			name: decoded.name,
			avatarURL: decoded.avatarURL,
			userID: decoded.userID,
			role: decoded.role,
			verified: decoded.verified
		}
	};
};
