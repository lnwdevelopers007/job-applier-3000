import { goto } from '$app/navigation';
import { authStore, type UserInfo } from '$lib/stores/auth.svelte';

export function isAuthenticated(): boolean {
	return authStore.isAuthenticated;
}

export function getUserInfo(): UserInfo | null {
	// Add a small retry mechanism for timing issues
	const user = authStore.user;
	if (user && user.userID) { // Make sure we have essential data
		return {
			name: user.name || 'Guest',
			email: user.email || '',
			role: user.role || 'jobSeeker',
			avatarURL: user.avatarURL,
			userID: user.userID,
			verified: user.verified
		};
	}
	return null;
}

export function logout() {
	authStore.logout();
}

export function requireAuth(callback: () => void, pendingAction?: { key: string; value: string }) {
	if (isAuthenticated()) {
		callback();
	} else {
		// Store pending action if provided
		if (pendingAction) {
			sessionStorage.setItem(pendingAction.key, pendingAction.value);
		}
		// Return false to indicate auth required - component can show modal
		return false;
	}
	return true;
}

export function navigateWithAuth(path: string, pendingAction?: { key: string; value: string }) {
	const authSuccess = requireAuth(() => {
		// Check verification before navigation
		const user = authStore.user;
		if (user && !user.verified && path !== '/unverified' && path !== '/') {
			// Check if trying to access protected routes
			if (path.startsWith('/app/') || path.startsWith('/company/') || path.startsWith('/faculty/') || path.startsWith('/admin/')) {
				goto('/unverified');
				return;
			}
		}
		goto(path);
	}, pendingAction);
	return authSuccess;
}
