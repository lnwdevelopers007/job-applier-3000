import { userApi } from '$lib/api';
import { getUserInfo, isAuthenticated } from '$lib/utils/auth';
import { authStore } from '$lib/stores/auth.svelte';
import { processBackendUser, transformToBackendFormat, transformToFrontendFormat } from '$lib/utils/userTransforms';
import type { 
  User, 
  UserFilters, 
  UpdateUserPayload 
} from '$lib/types';

export class UserService {
	/**
	 * Query users with filters - GET /users/query
	 */
	static async queryUsers(filters?: UserFilters): Promise<User[]> {
		try {
			return await userApi.query(filters);
		} catch (error) {
			console.error('Error querying users:', error);
			throw error;
		}
	}

	/**
	 * Get all users - GET /users/
	 */
	static async getAllUsers(): Promise<User[]> {
		try {
			return await userApi.getAll();
		} catch (error) {
			console.error('Error fetching all users:', error);
			throw error;
		}
	}

	/**
	 * Get a specific user by ID - GET /users/:id
	 */
	static async getUserById(id: string): Promise<User> {
		try {
			const user = await userApi.getById(id);
			if (!user) {
				throw new Error('User not found or unauthorized');
			}
			return processBackendUser(user);
		} catch (error) {
			console.error('Error fetching user:', error);
			throw error;
		}
	}

	/**
	 * Create a new user - POST /users/
	 */
	static async createUser(userData: Partial<User>): Promise<User> {
		try {
			return await userApi.create(userData);
		} catch (error) {
			console.error('Error creating user:', error);
			throw error;
		}
	}

	/**
	 * Update a user - PUT /users/:id
	 */
	static async updateUser(id: string, userData: UpdateUserPayload): Promise<User> {
		try {
			const user = await userApi.update(id, userData);
			return processBackendUser(user);
		} catch (error) {
			console.error('Error updating user:', error);
			throw error;
		}
	}

	/**
	 * Delete a user - DELETE /users/:id
	 */
	static async deleteUser(id: string): Promise<void> {
		try {
			await userApi.delete(id);
		} catch (error) {
			console.error('Error deleting user:', error);
			throw error;
		}
	}

	/**
	 * Approve a user (set verified to true) - PATCH /users/:id/verify
	 */
	static async approveUser(id: string): Promise<void> {
		try {
			await userApi.updateVerification(id, true);
		} catch (error) {
			console.error('Error approving user:', error);
			throw error;
		}
	}

	/**
	 * Update user role - PATCH /users/:id/role
	 */
	static async updateUserRole(id: string, role: string): Promise<void> {
		try {
			await userApi.updateRole(id, role);
		} catch (error) {
			console.error('Error updating user role:', error);
			throw error;
		}
	}

	// === BUSINESS LOGIC METHODS ===

	/**
	 * Get current user profile (business logic wrapper)
	 */
	static async getCurrentUser(): Promise<User> {
		try {
			if (!isAuthenticated()) {
				throw new Error('No user data found');
			}
			
			const userInfo = getUserInfo();
			if (!userInfo?.userID) {
				throw new Error('User ID not found');
			}
			
			return this.getUserById(userInfo.userID);
		} catch (error) {
			console.error('Error getting current user:', error);
			throw error;
		}
	}

	/**
	 * Update current user profile with auth store sync
	 */
	static async updateCurrentUser(data: UpdateUserPayload): Promise<User> {
		try {
			const userInfo = getUserInfo();
			if (!userInfo?.userID) {
				throw new Error('User ID not found');
			}
			
			const updatedUser = await this.updateUser(userInfo.userID, data);
			
			// Update auth store with new user data
			if (authStore.user) {
				authStore.updateUser({
					name: updatedUser.name || authStore.user.name,
					email: updatedUser.email || authStore.user.email,
					role: (updatedUser.role || authStore.user.role) as 'jobSeeker' | 'company' | 'admin',
					avatarURL: updatedUser.avatarURL || authStore.user.avatarURL,
					userID: authStore.user.userID,
					verified: authStore.user.verified
				});
			}
			
			return updatedUser;
		} catch (error) {
			console.error('Error updating current user:', error);
			throw error;
		}
	}
	
	// === HELPER METHODS ===

	/**
	 * Transform frontend data to backend format
	 */
	static transformToBackendFormat = transformToBackendFormat;

	/**
	 * Transform backend data to frontend format
	 */
	static transformToFrontendFormat = transformToFrontendFormat;
}

// Export for backward compatibility (can be removed when all imports are updated)
export const userService = UserService;