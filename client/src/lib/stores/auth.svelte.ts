import { goto } from '$app/navigation';
import { browser } from '$app/environment';

export type UserRole = 'jobSeeker' | 'company' | 'faculty' | 'admin';

export type UserInfo = {
	name: string;
	email: string;
	role: UserRole;
	avatarURL?: string;
	userID: string;
	verified: boolean;
};

export interface AuthState {
	isAuthenticated: boolean;
	user: UserInfo | null;
	isLoading: boolean;
}

class AuthStore {
	private state = $state<AuthState>({
		isAuthenticated: false,
		user: null,
		isLoading: false
	});

	// Initialize from server data (called from +layout.svelte)
	initFromPageData(userData: App.Locals['user']) {
		if (userData && userData.isAuthenticated) {
			this.state.isAuthenticated = true;
			this.state.user = {
				name: userData.name,
				email: userData.email,
				role: userData.role,
				avatarURL: userData.avatarURL,
				userID: userData.userID,
				verified: userData.verified
			};
		} else {
			this.state.isAuthenticated = false;
			this.state.user = null;
		}
		this.state.isLoading = false;
	}

	get isAuthenticated() {
		return this.state.isAuthenticated;
	}

	get user() {
		return this.state.user;
	}

	get role() {
		return this.state.user?.role || null;
	}

	get isLoading() {
		return this.state.isLoading;
	}

	hasRole(requiredRole: UserRole | UserRole[]): boolean {
		if (!this.state.isAuthenticated || !this.state.user) {
			return false;
		}

		const roles = Array.isArray(requiredRole) ? requiredRole : [requiredRole];
		return roles.includes(this.state.user.role);
	}

	updateUser(updates: Partial<UserInfo>) {
		if (this.state.user) {
			this.state.user = { ...this.state.user, ...updates };
		}
	}

	clearUser() {
		// Clear local state only (keep cookies for signup flow)
		this.state.isAuthenticated = false;
		this.state.user = null;
	}

	async logoutAndClearTokens() {
		// Clear cookies by calling logout endpoint but don't redirect
		if (browser) {
			try {
				await fetch('/api/logout', {
					method: 'POST',
					credentials: 'include'
				});
			} catch (error) {
				console.error('Logout failed:', error);
			}
			
			// Clear local state
			this.clearUser();
		}
	}

	async logout() {
		// Clear cookies and redirect to home
		await this.logoutAndClearTokens();
		if (browser) {
			goto('/');
		}
	}
}

export const authStore = new AuthStore();