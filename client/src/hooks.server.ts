import type { Handle } from '@sveltejs/kit';
import { redirect } from '@sveltejs/kit';
import { jwtDecode } from 'jwt-decode';

// Routes that don't require authentication
const publicRoutes = [
	'/',
	'/login',
	'/signup',
	'/signup/student', 
	'/signup/company',
	'/callback',
	'/unverified'
];

type JWTPayload = {
	email?: string;
	name?: string;
	avatarURL?: string;
	userID?: string;
	role?: string;
	verified?: boolean;
	exp?: number;
};

export const handle: Handle = async ({ event, resolve }) => {
	const path = event.url.pathname;
	
	const isPublicRoute = publicRoutes.some(route => 
		path === route || path.startsWith(route + '/')
	);

	const accessToken = event.cookies.get('access_token');
	const refreshToken = event.cookies.get('refresh_token');

	if (accessToken) {
		let decoded: JWTPayload | null = null;
		try {
			decoded = jwtDecode<JWTPayload>(accessToken);
		} catch (error) {
			console.error('Failed to decode JWT:', error);
		}
		
		if (decoded) {
			// Check if token is expired
			const isExpired = decoded.exp && decoded.exp * 1000 < Date.now();
			
			if (isExpired && refreshToken) {
				// Try to refresh the token
				try {
					const backendUrl = import.meta.env.VITE_BACKEND || 'http://localhost:8080';
					const refreshResponse = await fetch(`${backendUrl}/auth/refresh`, {
						method: 'POST',
						headers: {
							'Cookie': `refresh_token=${refreshToken}`
						},
						credentials: 'include'
					});

					if (refreshResponse.ok) {
						const data = await refreshResponse.json();
						const newDecoded = jwtDecode<JWTPayload>(data.access_token);
						const jwtExpiresIn = newDecoded.exp ? Math.max(0, newDecoded.exp - Math.floor(Date.now() / 1000)) : 60 * 60 * 24;
						
						event.cookies.set('access_token', data.access_token, {
							path: '/',
							httpOnly: true,
							secure: import.meta.env.MODE === 'production',
							sameSite: 'lax',
							maxAge: jwtExpiresIn
						});
						
						if (newDecoded) {
							const normalizedRole = normalizeRole(newDecoded.role);
							
							// Set user in locals
							event.locals.user = {
								email: newDecoded.email || '',
								name: newDecoded.name || '',
								avatarURL: newDecoded.avatarURL,
								userID: newDecoded.userID || '',
								role: normalizedRole,
								verified: newDecoded.verified || false,
								isAuthenticated: true
							};
						}
					} else {
						// Refresh failed, clear cookies
						event.cookies.delete('access_token', { path: '/' });
						event.cookies.delete('refresh_token', { path: '/' });
						
						if (!isPublicRoute) {
							throw redirect(303, `/login?returnUrl=${encodeURIComponent(path)}`);
						}
					}
				} catch (error) {
					console.error('Token refresh failed:', error);
					event.cookies.delete('access_token', { path: '/' });
					event.cookies.delete('refresh_token', { path: '/' });
					
					if (!isPublicRoute) {
						throw redirect(303, `/login?returnUrl=${encodeURIComponent(path)}`);
					}
				}
			} else if (!isExpired) {
				// Token is valid, set user in locals
				const normalizedRole = normalizeRole(decoded.role);
				
				event.locals.user = {
					email: decoded.email || '',
					name: decoded.name || '',
					avatarURL: decoded.avatarURL,
					userID: decoded.userID || '',
					role: normalizedRole,
					verified: decoded.verified || false,
					isAuthenticated: true
				};
			} else {
				// Token expired and no refresh token
				event.cookies.delete('access_token', { path: '/' });
				
				if (!isPublicRoute) {
					throw redirect(303, `/login?returnUrl=${encodeURIComponent(path)}`);
				}
			}
		}
	} else if (!isPublicRoute) {
		// No access token and trying to access protected route
		if (refreshToken) {
			// Try to get new access token with refresh token
			try {
				const backendUrl = import.meta.env.VITE_BACKEND || 'http://localhost:8080';
				const refreshResponse = await fetch(`${backendUrl}/auth/refresh`, {
					method: 'POST',
					headers: {
						'Cookie': `refresh_token=${refreshToken}`
					},
					credentials: 'include'
				});

				if (refreshResponse.ok) {
					const data = await refreshResponse.json();
					const newDecoded = jwtDecode<JWTPayload>(data.access_token);
					const refreshExpiresIn = newDecoded.exp ? Math.max(0, newDecoded.exp - Math.floor(Date.now() / 1000)) : 60 * 60 * 24;
					
					event.cookies.set('access_token', data.access_token, {
						path: '/',
						httpOnly: true,
						secure: import.meta.env.MODE === 'production',
						sameSite: 'lax',
						maxAge: refreshExpiresIn
					});
					
					if (newDecoded) {
						const normalizedRole = normalizeRole(newDecoded.role);
						
						event.locals.user = {
							email: newDecoded.email || '',
							name: newDecoded.name || '',
							avatarURL: newDecoded.avatarURL,
							userID: newDecoded.userID || '',
							role: normalizedRole,
							verified: newDecoded.verified || false,
							isAuthenticated: true
						};
					}
				} else {
					event.cookies.delete('refresh_token', { path: '/' });
					throw redirect(303, `/login?returnUrl=${encodeURIComponent(path)}`);
				}
			} catch (error) {
				console.error('Token refresh failed:', error);
				event.cookies.delete('refresh_token', { path: '/' });
				throw redirect(303, `/login?returnUrl=${encodeURIComponent(path)}`);
			}
		} else {
			// No tokens at all
			throw redirect(303, `/login?returnUrl=${encodeURIComponent(path)}`);
		}
	}

	// Role-based access control
	if (event.locals.user?.isAuthenticated) {
		const userRole = event.locals.user.role;
		
		// Redirect authenticated users away from login/signup
		if (path === '/login' || path === '/signup') {
			const redirectPath = getDefaultPathForRole(userRole);
			throw redirect(303, redirectPath);
		}
		
		// Check role-specific access
		if (path.startsWith('/app/') && userRole !== 'jobSeeker') {
			throw redirect(303, getDefaultPathForRole(userRole));
		}
		
		if (path.startsWith('/company/') && userRole !== 'company') {
			throw redirect(303, getDefaultPathForRole(userRole));
		}
		
		if (path.startsWith('/admin/') && userRole !== 'admin') {
			throw redirect(303, getDefaultPathForRole(userRole));
		}
	}

	// Continue with request
	const response = await resolve(event);
	return response;
};

// Helper function to normalize role names
function normalizeRole(role?: string): 'jobSeeker' | 'company' | 'admin' {
	if (!role) return 'jobSeeker';
	
	const normalized = role.toLowerCase();
	if (normalized === 'jobseeker' || normalized === 'job_seeker') {
		return 'jobSeeker';
	}
	if (normalized === 'company') {
		return 'company';
	}
	if (normalized === 'admin') {
		return 'admin';
	}
	return 'jobSeeker';
}

// Helper function to get default path for role
function getDefaultPathForRole(role: string): string {
	switch (role) {
		case 'company':
			return '/company/dashboard';
		case 'admin':
			return '/admin/dashboard';
		case 'jobSeeker':
		default:
			return '/app/jobs';
	}
}