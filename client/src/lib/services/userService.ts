import { apiFetch } from '$lib/utils/api';

const API_BASE = import.meta.env.VITE_BACKEND || 'http://localhost:8080';

export interface User {
	id: string;
	userID: string;
	provider: string;
	email: string;
	name: string;
	avatarURL: string;
	role: string;
	verified: boolean;
	updatedAt: string;
	createdAt: string;
	userInfo?: JobSeekerInfo | CompanyInfo;
}

export interface JobSeekerInfo {
	fullName?: string;
	location: string;
	phone: string;
	linkedIn: string;
	desiredRole?: string;
	aboutMe?: string;
	dateOfBirth?: string;
	gender?: string;
	portfolio?: string;
	github?: string;
}

export interface CompanyInfo {
	id?: string;
	userID: string;
	name: string;
	aboutUs: string;
	industry: string;
	size: string;
	website: string;
	logo?: string;
	foundedYear?: string;
	headquarters?: string;
	linkedIn?: string;
	benefits?: string[];
}

export interface FileDocument {
	id?: string;
	parentID: string;
	parentColl: string;
	content?: string;
	fileExtension: string;
	name?: string;
	size?: string;
	uploadDate?: string;
}

export interface UpdateUserPayload {
	name?: string;
	email?: string;
	avatarURL?: string;
	provider?: string;
	userID?: string;
	role?: string;
	verified?: boolean;
	userInfo?: Partial<JobSeekerInfo | CompanyInfo>;
	files?: string[];
}

class UserService {
	// Get current user profile
	async getCurrentUser(): Promise<User> {
		const userStr = localStorage.getItem('user');
		if (!userStr) {
			throw new Error('No user data found');
		}
		
		const userData = JSON.parse(userStr);
		const userId = userData.id || userData._id || userData.userID;
		
		if (!userId) {
			throw new Error('User ID not found');
		}
		
		const response = await apiFetch(`${API_BASE}/users/${userId}`);
		if (!response.ok) {
			throw new Error('Failed to fetch user profile');
		}
		
		const user = await response.json();
		
		// Map _id from backend to id for frontend
		if (user._id) {
			user.id = user._id;
		}
		
		return user;
	}
	
	// Update user profile
	async updateUser(userId: string, data: UpdateUserPayload): Promise<User> {
		const response = await apiFetch(`${API_BASE}/users/${userId}`, {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data)
		});
		
		if (!response.ok) {
			const error = await response.text();
			throw new Error(`Failed to update profile: ${error}`);
		}
		
		const updatedUser = await response.json();
		
		// Map _id from backend to id for frontend
		if (updatedUser._id) {
			updatedUser.id = updatedUser._id;
		}
		
		// Update localStorage with new user data, but only update specific fields
		const currentUser = JSON.parse(localStorage.getItem('user') || '{}');
		localStorage.setItem('user', JSON.stringify({
			...currentUser,
			name: updatedUser.name || currentUser.name,
			email: updatedUser.email || currentUser.email,
			role: updatedUser.role || currentUser.role,
			avatarURL: updatedUser.avatarURL || currentUser.avatarURL,
			userID: currentUser.userID // Keep the original userID
		}));
		
		return updatedUser;
	}
	
	// Upload document (Not implemented - console log only)
	async uploadDocument(file: File): Promise<Record<string, unknown>> {
		console.log('📤 Document upload not implemented yet for:', file.name);
		return {};
	}
	
	// Delete document (Not implemented - console log only)
	async deleteDocument(documentId: string): Promise<void> {
		console.log('🗑️ Document delete not implemented yet for ID:', documentId);
	}
	
	// Update password
	async updatePassword(currentPassword: string, newPassword: string): Promise<void> {
		const response = await apiFetch(`${API_BASE}/users/password`, {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				currentPassword,
				newPassword
			})
		});
		
		if (!response.ok) {
			const error = await response.text();
			throw new Error(`Failed to update password: ${error}`);
		}
	}
	
	// Helper to transform frontend data to backend format
	// Only send changed fields while always preserving userID and provider
	transformToBackendFormat(userData: Record<string, unknown>, userType: string, changedFields?: string[]): UpdateUserPayload {
		const payload: UpdateUserPayload = {};
		
		// Always preserve userID and provider fields (needed for OAuth disconnect)
		payload.provider = (userData.provider as string) || '';
		payload.userID = (userData.userID as string) || '';
		
		// If no changedFields specified, send all fields (backwards compatibility)
		if (!changedFields) {
			payload.name = (userData.name as string) || '';
			payload.email = (userData.email as string) || '';
			payload.avatarURL = (userData.avatar as string) || '';
			payload.role = (userData.role as string) || (userType === 'seeker' ? 'jobSeeker' : 'company');
			payload.verified = (userData.verified as boolean) || false;
			
			// Build userInfo object based on user type
			if (userType === 'seeker' || userType === 'jobSeeker') {
				payload.userInfo = {
					fullName: (userData.fullName as string) || '',
					location: (userData.location as string) || '',
					phone: (userData.phone as string) || '',
					linkedIn: (userData.linkedin as string) || '',
					desiredRole: (userData.desiredRole as string) || '',
					aboutMe: (userData.aboutMe as string) || '',
					dateOfBirth: (userData.dateOfBirth as string) || '',
					portfolio: (userData.portfolio as string) || '',
					github: (userData.github as string) || ''
				};
			} else if (userType === 'company') {
				payload.userInfo = {
					name: (userData.companyName as string) || '',
					aboutUs: (userData.aboutCompany as string) || '',
					industry: (userData.industry as string) || '',
					size: (userData.companySize as string) || '',
					website: (userData.companyWebsite as string) || '',
					logo: (userData.companyLogo as string) || '',
					foundedYear: (userData.foundedYear as string) || '',
					headquarters: (userData.headquarters as string) || '',
					linkedIn: (userData.companyLinkedin as string) || ''
				};
			}
		} else {
			// Only include fields that have changed
			if (changedFields.includes('name') && userData.name !== undefined) {
				payload.name = userData.name as string;
			}
			if (changedFields.includes('email') && userData.email !== undefined) {
				payload.email = userData.email as string;
			}
			if (changedFields.includes('avatar') && userData.avatar !== undefined) {
				payload.avatarURL = userData.avatar as string;
			}
			if (changedFields.includes('role') && userData.role !== undefined) {
				payload.role = userData.role as string;
			}
			if (changedFields.includes('verified') && userData.verified !== undefined) {
				payload.verified = userData.verified as boolean;
			}
			
			// Handle userInfo fields - Always send all fields since BSON can't use omitempty
			if (userType === 'seeker' || userType === 'jobSeeker') {
				const userInfoFieldNames = [
					'fullName', 'location', 'phone', 'linkedin', 'desiredRole', 
					'aboutMe', 'dateOfBirth', 'portfolio', 'github'
				];
				
				// Check if any userInfo field has changed
				const hasUserInfoChanges = userInfoFieldNames.some(field => 
					changedFields.includes(field)
				);
				
				// If any userInfo field changed, send ALL userInfo fields
				if (hasUserInfoChanges) {
					payload.userInfo = {
						fullName: (userData.fullName as string) || '',
						location: (userData.location as string) || '',
						phone: (userData.phone as string) || '',
						linkedIn: (userData.linkedin as string) || '',
						desiredRole: (userData.desiredRole as string) || '',
						aboutMe: (userData.aboutMe as string) || '',
						dateOfBirth: (userData.dateOfBirth as string) || '',
						portfolio: (userData.portfolio as string) || '',
						github: (userData.github as string) || ''
					};
				}
			} else if (userType === 'company') {
				const userInfoFieldNames = [
					'companyName', 'aboutCompany', 'industry', 'companySize', 
					'companyWebsite', 'companyLogo', 'foundedYear', 'headquarters', 'companyLinkedin'
				];
				
				// Check if any userInfo field has changed
				const hasUserInfoChanges = userInfoFieldNames.some(field => 
					changedFields.includes(field)
				);
				
				// If any userInfo field changed, send ALL userInfo fields
				if (hasUserInfoChanges) {
					payload.userInfo = {
						name: (userData.companyName as string) || '',
						aboutUs: (userData.aboutCompany as string) || '',
						industry: (userData.industry as string) || '',
						size: (userData.companySize as string) || '',
						website: (userData.companyWebsite as string) || '',
						logo: (userData.companyLogo as string) || '',
						foundedYear: (userData.foundedYear as string) || '',
						headquarters: (userData.headquarters as string) || '',
						linkedIn: (userData.companyLinkedin as string) || ''
					};
				}
			}
		}
		
		return payload;
	}
	
	// Helper to transform backend data to frontend format
	transformToFrontendFormat(user: User): Record<string, unknown> {
		const frontendData: Record<string, unknown> = {
			id: user.id,
			fullName: user.name || '',
			email: user.email || '',
			avatar: user.avatarURL || '',
			googleConnected: user.provider === 'google',
			// Initialize all fields with empty strings to avoid undefined binding
			currentPassword: '',
			newPassword: '',
			confirmPassword: ''
		};
		
		if (user.role === 'jobseeker' || user.role === 'jobSeeker') {
			// Initialize all seeker fields with defaults
			frontendData.fullName = '';
			frontendData.location = '';
			frontendData.phone = '';
			frontendData.linkedin = '';
			frontendData.desiredRole = '';
			frontendData.aboutMe = '';
			frontendData.dateOfBirth = '';
			frontendData.gender = '';
			frontendData.portfolio = '';
			frontendData.github = '';
			
			// Override with actual values if they exist
			if (user.userInfo) {
				const info = user.userInfo as JobSeekerInfo;
				frontendData.fullName = info.fullName || '';
				frontendData.location = info.location || '';
				frontendData.phone = info.phone || '';
				frontendData.linkedin = info.linkedIn || '';
				frontendData.desiredRole = info.desiredRole || '';
				frontendData.aboutMe = info.aboutMe || '';
				frontendData.dateOfBirth = info.dateOfBirth || '';
				frontendData.gender = info.gender || '';
				frontendData.portfolio = info.portfolio || '';
				frontendData.github = info.github || '';
			}
		} else if (user.role === 'company') {
			// Initialize all company fields with defaults
			frontendData.companyName = '';
			frontendData.aboutCompany = '';
			frontendData.industry = '';
			frontendData.companySize = '';
			frontendData.companyWebsite = '';
			frontendData.companyLogo = '';
			frontendData.foundedYear = '';
			frontendData.headquarters = '';
			frontendData.companyLinkedin = '';
			frontendData.benefits = [];
			
			// Override with actual values if they exist
			if (user.userInfo) {
				const info = user.userInfo as CompanyInfo;
				frontendData.companyName = info.name || '';
				frontendData.aboutCompany = info.aboutUs || '';
				frontendData.industry = info.industry || '';
				frontendData.companySize = info.size || '';
				frontendData.companyWebsite = info.website || '';
				frontendData.companyLogo = info.logo || '';
				frontendData.foundedYear = info.foundedYear || '';
				frontendData.headquarters = info.headquarters || '';
				frontendData.companyLinkedin = info.linkedIn || '';
				frontendData.benefits = info.benefits || [];
			}
		}
		
		// Initialize documents array
		frontendData.documents = [];
		
		return frontendData;
	}
}

export const userService = new UserService();