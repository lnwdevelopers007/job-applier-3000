import { goto } from '$app/navigation';

export type UserInfo = {
	name: string;
	email: string;
	role: string;
	avatarURL?: string;
	userID: string;
};

export function isAuthenticated(): boolean {
	return !!localStorage.getItem('access_token');
}

export function getUserInfo(): UserInfo | null {
	const storedUser = localStorage.getItem('user');
	if (storedUser) {
		try {
			const parsed = JSON.parse(storedUser);
			return {
				name: parsed.name || 'Guest',
				email: parsed.email || '',
				role: parsed.role || 'User',
				avatarURL: parsed.avatarURL,
				userID: parsed.userID
			};
		} catch (err) {
			console.error('Failed to parse stored user info:', err);
		}
	}
	return null;
}

export function logout() {
	localStorage.removeItem('access_token');
	localStorage.removeItem('refresh_token');
	localStorage.removeItem('user');
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
	const authSuccess = requireAuth(() => goto(path), pendingAction);
	return authSuccess;
}
