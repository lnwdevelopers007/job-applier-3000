// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			user?: {
				email: string;
				name: string;
				avatarURL?: string;
				userID: string;
				role: 'jobSeeker' | 'company' | 'faculty' | 'admin';
				verified: boolean;
				isAuthenticated: boolean;
			};
		}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
