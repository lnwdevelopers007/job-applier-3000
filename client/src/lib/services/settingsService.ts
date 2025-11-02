import { UserService } from './userService';
import type { 
  User, 
  JobSeekerInfo, 
  CompanyInfo, 
  BaseUserData,
  SeekerUserData,
  CompanyUserData,
  UserData,
  Tab
} from '$lib/types';

// Field mappings for each user type
const SEEKER_FIELD_MAPPING = {
	fullName: (info: JobSeekerInfo) => info.fullName || '',
	location: (info: JobSeekerInfo) => info.location || '',
	phone: (info: JobSeekerInfo) => info.phone || '',
	linkedin: (info: JobSeekerInfo) => info.linkedIn || '',
	desiredRole: (info: JobSeekerInfo) => info.desiredRole || '',
	aboutMe: (info: JobSeekerInfo) => info.aboutMe || '',
	dateOfBirth: (info: JobSeekerInfo) => info.dateOfBirth || '',
	portfolio: (info: JobSeekerInfo) => info.portfolio || '',
	github: (info: JobSeekerInfo) => info.github || '',
	skills: (info: JobSeekerInfo) => {
		if (info.skills) {
			// Handle both string and array formats
			if (typeof info.skills === 'string') {
				return info.skills.split(',').map(s => s.trim()).filter(s => s.length > 0);
			} else if (Array.isArray(info.skills)) {
				return info.skills;
			}
		}
		return [];
	}
} as const;

const COMPANY_FIELD_MAPPING = {
	companyName: (info: CompanyInfo) => info.name || '',
	aboutCompany: (info: CompanyInfo) => info.aboutUs || '',
	industry: (info: CompanyInfo) => info.industry || '',
	companySize: (info: CompanyInfo) => info.size || '',
	companyWebsite: (info: CompanyInfo) => info.website || '',
	companyLogo: (info: CompanyInfo) => info.logo || '',
	foundedYear: (info: CompanyInfo) => info.foundedYear || '',
	headquarters: (info: CompanyInfo) => info.headquarters || '',
	companyLinkedin: (info: CompanyInfo) => info.linkedIn || ''
} as const;

export const SEEKER_TABS: Tab[] = [
	{
		id: 'user',
		label: 'User',
		title: 'User info',
		description: 'Update your photo and user details here.'
	},
	{
		id: 'personal',
		label: 'Personal Info',
		title: 'Personal Info',
		description: 'Manage your personal information that will be shared with employers.'
	},
	{
		id: 'documents',
		label: 'Documents',
		title: 'Documents',
		description: 'Upload and manage your documents and files.'
	}
];

export const COMPANY_TABS: Tab[] = [
	{
		id: 'user',
		label: 'User',
		title: 'User info',
		description: 'Update your admin account details.'
	},
	{
		id: 'company',
		label: 'Company',
		title: 'Company Info',
		description: 'Manage your company information and branding.'
	},
	{
		id: 'documents',
		label: 'Documents',
		title: 'Company Documents',
		description: 'Upload and manage company documents and certifications.'
	}
];

// Initial data templates
export const INITIAL_SEEKER_DATA: SeekerUserData = {
	id: '',
	name: '',
	email: '',
	avatar: '',
	provider: '',
	userID: '',
	role: 'jobSeeker',
	verified: false,
	googleConnected: false,
	fullName: '',
	location: '',
	phone: '',
	linkedin: '',
	desiredRole: '',
	aboutMe: '',
	dateOfBirth: '',
	portfolio: '',
	github: '',
	skills: [],
	documents: []
};

export const INITIAL_COMPANY_DATA: CompanyUserData = {
	id: '',
	name: '',
	email: '',
	avatar: '',
	provider: '',
	userID: '',
	role: 'company',
	verified: false,
	googleConnected: false,
	companyName: '',
	aboutCompany: '',
	industry: '',
	companySize: '',
	companyWebsite: '',
	companyLogo: '',
	foundedYear: '',
	headquarters: '',
	companyLinkedin: '',
	documents: []
};

export class SettingsService {
	// Map backend user data to frontend format
	static mapUserToUserData<T extends UserData>(user: User, userType: 'seeker' | 'company'): T {
		const baseData: BaseUserData = {
			id: user.id || (user as any)._id as string,
			name: user.name || '',
			email: user.email || '',
			avatar: user.avatarURL || '',
			provider: user.provider || '',
			userID: user.userID || '',
			role: user.role || userType === 'seeker' ? 'jobSeeker' : 'company',
			verified: user.verified || false,
			googleConnected: user.provider === 'google',
			documents: []
		};

		if (userType === 'seeker') {
			const userInfo = user.userInfo as JobSeekerInfo | undefined;
			const seekerData: SeekerUserData = {
				...baseData,
				...Object.fromEntries(
					Object.entries(SEEKER_FIELD_MAPPING).map(([key, mapper]) => [
						key,
						userInfo ? mapper(userInfo) : ''
					])
				)
			};
			return seekerData as T;
		} else {
			const userInfo = user.userInfo as CompanyInfo | undefined;
			const companyData: CompanyUserData = {
				...baseData,
				...Object.fromEntries(
					Object.entries(COMPANY_FIELD_MAPPING).map(([key, mapper]) => [
						key,
						userInfo ? mapper(userInfo) : ''
					])
				)
			};
			return companyData as T;
		}
	}

	// Load user data from backend
	static async loadUserData<T extends UserData>(userType: 'seeker' | 'company'): Promise<T> {
		try {
			const user = await UserService.getCurrentUser();
			return this.mapUserToUserData<T>(user, userType);
		} catch (error) {
			console.error('Failed to load user data:', error);
			throw new Error('Failed to load profile data. Please try again.');
		}
	}

	// Save user profile data
	static async saveUserData(
		userData: UserData,
		userType: 'seeker' | 'company',
		changedFields?: string[]
	): Promise<User> {
		const userId = userData.id;
		if (!userId) {
			throw new Error('User ID not found');
		}

		const payload = UserService.transformToBackendFormat(userData as Record<string, unknown>, userType, changedFields);
		
		// Handle file uploads for company logo
		if (userType === 'company' && (userData as CompanyUserData & { companyLogoFile?: File }).companyLogoFile) {
			return new Promise((resolve, reject) => {
				const reader = new FileReader();
				reader.onload = async (e) => {
					try {
						if (payload.userInfo && e.target?.result) {
							(payload.userInfo as Record<string, unknown>).logo = e.target.result as string;
						}
						const result = await UserService.updateUser(userId, payload);
						resolve(result);
					} catch (error) {
						reject(error);
					}
				};
				reader.onerror = () => reject(new Error('Failed to read file'));
				reader.readAsDataURL((userData as CompanyUserData & { companyLogoFile: File }).companyLogoFile);
			});
		}

		return UserService.updateUser(userId, payload);
	}

	// Update local userData with backend response
	static updateLocalData<T extends UserData>(
		userData: T,
		updatedUser: User,
		userType: 'seeker' | 'company'
	): void {
		// Update common fields
		if (updatedUser.name) userData.name = updatedUser.name;
		if (updatedUser.avatarURL) userData.avatar = updatedUser.avatarURL;

		// Update type-specific fields
		if (updatedUser.userInfo) {
			if (userType === 'seeker') {
				const info = updatedUser.userInfo as JobSeekerInfo;
				const seekerData = userData as SeekerUserData;
				Object.entries(SEEKER_FIELD_MAPPING).forEach(([key]) => {
					const backendKey = key === 'linkedin' ? 'linkedIn' : key;
					const infoRecord = info as unknown as Record<string, unknown>;
					const seekerRecord = seekerData as unknown as Record<string, unknown>;
					if (infoRecord[backendKey] !== undefined) {
						seekerRecord[key] = infoRecord[backendKey];
					}
				});
			} else if (userType === 'company') {
				const info = updatedUser.userInfo as CompanyInfo;
				const companyData = userData as CompanyUserData;
				Object.entries(COMPANY_FIELD_MAPPING).forEach(([key]) => {
					let backendKey = key;
					// Handle field name mappings
					if (key === 'companyName') backendKey = 'name';
					else if (key === 'aboutCompany') backendKey = 'aboutUs';
					else if (key === 'companySize') backendKey = 'size';
					else if (key === 'companyWebsite') backendKey = 'website';
					else if (key === 'companyLogo') backendKey = 'logo';
					else if (key === 'companyLinkedin') backendKey = 'linkedIn';

					const infoRecord = info as unknown as Record<string, unknown>;
					const companyRecord = companyData as unknown as Record<string, unknown>;
					if (infoRecord[backendKey] !== undefined) {
						companyRecord[key] = infoRecord[backendKey];
					}
				});
			}
		}
	}

	// Handle document uploads and saving (Not implemented)
	static async saveDocuments(userData: UserData): Promise<void> {
		console.log('📄 Document upload not implemented yet');
		console.log('User ID:', userData.id);
		console.log('Documents that would be uploaded:', userData.documents);
	}
}