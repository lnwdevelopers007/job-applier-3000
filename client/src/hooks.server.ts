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
	'/unverified',
	'/banned'
];

type JWTPayload = {
	email?: string;
	name?: string;
	avatarURL?: string;
	userID?: string;
	role?: string;
	verified?: boolean;
	banned?: boolean;
	exp?: number;
};

// Helper function to handle user verification and setup
function handleUserAuth(decoded: JWTPayload, event: any, path: string) {
	// Set user in locals
	event.locals.user = {
		email: decoded.email || '',
		name: decoded.name || '',
		avatarURL: decoded.avatarURL,
		userID: decoded.userID || '',
		role: (decoded.role as any) || 'jobSeeker',
		verified: decoded.verified || false,
		isAuthenticated: true // User has valid token, even if not verified
	};
	
	// Check if user is verified (allow access to /unverified route, signup flow, and callback)
	const isSignupFlow = path.startsWith('/signup/student') || path.startsWith('/signup/company');
	const isCallbackPath = path === '/callback';
	const isUnverifiedPath = path === '/unverified' || path.startsWith('/unverified/');
	const isAllowedUnverifiedPath = isCallbackPath || isUnverifiedPath || isSignupFlow;
	
	if (!decoded.verified && !isAllowedUnverifiedPath) {
		const userName = decoded.name ? encodeURIComponent(decoded.name) : '';
		throw redirect(303, `/unverified?name=${userName}`);
	}
}

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
			// Check if user is banned FIRST before anything else
			if (decoded.banned && path !== '/banned') {
				throw redirect(303, '/banned');
			}

			// Check if token is expired
			const isExpired = decoded.exp && decoded.exp * 1000 < Date.now();
			
			// Don't try to refresh tokens for unverified users in signup flow
			const isSignupFlow = path.startsWith('/signup/student') || path.startsWith('/signup/company');
			const skipRefresh = !decoded.verified && isSignupFlow;
			
			if (isExpired && refreshToken && !skipRefresh) {
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
						
						// Check if newly refreshed token shows user is banned
						if (newDecoded.banned) {
							throw redirect(303, '/banned');
						}

						const jwtExpiresIn = newDecoded.exp ? Math.max(0, newDecoded.exp - Math.floor(Date.now() / 1000)) : 60 * 60 * 24;
						
						event.cookies.set('access_token', data.access_token, {
							path: '/',
							httpOnly: true,
							secure: import.meta.env.MODE === 'production',
							sameSite: 'lax',
							maxAge: jwtExpiresIn
						});
						
						if (newDecoded) {
							handleUserAuth(newDecoded, event, path);
						}
					} else {
						// Check if refresh failed due to ban
						const errorData = await refreshResponse.json().catch(() => ({}));
						if (errorData.error === 'account_banned') {
							throw redirect(303, '/banned');
						}

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
				// Token is valid, handle user authentication
				handleUserAuth(decoded, event, path);
			} else if (skipRefresh) {
				// Token expired but skip refresh for unverified signup users - still allow access
				handleUserAuth(decoded, event, path);
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
					
					// Check if user is banned
					if (newDecoded.banned) {
						throw redirect(303, '/banned');
					}

					const refreshExpiresIn = newDecoded.exp ? Math.max(0, newDecoded.exp - Math.floor(Date.now() / 1000)) : 60 * 60 * 24;
					
					event.cookies.set('access_token', data.access_token, {
						path: '/',
						httpOnly: true,
						secure: import.meta.env.MODE === 'production',
						sameSite: 'lax',
						maxAge: refreshExpiresIn
					});
					
					if (newDecoded) {
						handleUserAuth(newDecoded, event, path);
					}
				} else {
					// Check if refresh failed due to ban
					const errorData = await refreshResponse.json().catch(() => ({}));
					if (errorData.error === 'account_banned') {
						throw redirect(303, '/banned');
					}

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
	if (event.locals.user?.isAuthenticated && event.locals.user?.verified) {
		const userRole = event.locals.user.role;
		
		// Redirect authenticated users away from login/signup to their last intended page or default
		if (path === '/login' || path === '/signup') {
			const returnUrl = event.url.searchParams.get('returnUrl');
			const redirectPath = returnUrl || getDefaultPathForRole(userRole);
			throw redirect(303, redirectPath);
		}
		
		// Check role-specific access
		if (path.startsWith('/app/jobs')) {
			// No restrictions - all authenticated users can access job routes
		} else if (path.startsWith('/app/') && userRole !== 'jobSeeker') {
			// Other /app/ routes are restricted to jobSeekers only
			throw redirect(303, getDefaultPathForRole(userRole));
		} else {
			// Check protected routes for specific roles
			const protectedRoutes = [
				{ path: '/company/', role: 'company' },
				{ path: '/faculty/', role: 'faculty' },
				{ path: '/admin/', role: 'admin' }
			] as const;
			
			for (const route of protectedRoutes) {
				if (path.startsWith(route.path) && userRole !== route.role) {
					throw redirect(303, getDefaultPathForRole(userRole));
				}
			}
		}
	}

	// Continue with request
	const response = await resolve(event);
	return response;
};

// Helper function to get default path for role
function getDefaultPathForRole(role: string): string {
	switch (role) {
		case 'company':
			return '/company/dashboard';
		case 'faculty':
			return '/app/jobs';
		case 'admin':
			return '/admin/jobs';
		case 'jobSeeker':
		default:
			return '/app/jobs';
	}
}