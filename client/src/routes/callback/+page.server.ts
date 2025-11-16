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

export const load: PageServerLoad = async ({ url, cookies }) => {
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

	// Check if user is verified before setting cookies
	if (!decoded.verified) {
		const userName = decoded.name ? encodeURIComponent(decoded.name) : '';
		throw redirect(303, `/unverified?name=${userName}`);
	}

	const jwtExpiresIn = decoded.exp
		? Math.max(0, decoded.exp - Math.floor(Date.now() / 1000))
		: 60 * 60 * 24;

	cookies.set('access_token', token, {
		path: '/',
		httpOnly: true,
		secure: import.meta.env.MODE === 'production',
		sameSite: 'lax',
		maxAge: jwtExpiresIn
	});

	// Determine redirect path based on role
	const role = decoded.role?.toLowerCase();
	const next = url.searchParams.get('step');
	let redirectPath = '/';
	if (role === 'company') {
		if (next === 'signup') {
			redirectPath = '/signup/company?currentStep=2';
		} else {
			redirectPath = '/company/dashboard';
		}
	}
	if (role === 'faculty') {
		redirectPath = '/app/jobs';
	} else if (role === 'admin') {
		redirectPath = '/admin/jobs';
	}
	if (role === 'jobseeker') {
		if (next === 'signup') {
			redirectPath = '/signup/student?currentStep=2';
		} else {
			redirectPath = '/app/jobs';
		}
	}
	if (role === 'admin') {
		redirectPath = '/admin/jobs';
	}

	throw redirect(303, redirectPath);
};
