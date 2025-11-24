/**
 * User data transformation utilities
 */

import type { 
  User, 
  JobSeekerInfo, 
  CompanyInfo, 
  UpdateUserPayload 
} from '$lib/types';

/**
 * Transform frontend data to backend format
 */
export function transformToBackendFormat(
  userData: Record<string, unknown>, 
  userType: string, 
  changedFields?: string[]
): UpdateUserPayload {
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
      payload.userInfo = buildSeekerInfo(userData);
    } else if (userType === 'company') {
      payload.userInfo = buildCompanyInfo(userData);
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
        'fullName', 'location', 'phone', 'linkedIn', 'desiredRole', 
        'aboutMe', 'dateOfBirth', 'portfolio', 'github', 'skills'
      ];
      
      // Check if any userInfo field has changed
      const hasUserInfoChanges = userInfoFieldNames.some(field => 
        changedFields.includes(field)
      );
      
      // If any userInfo field changed, send ALL userInfo fields
      if (hasUserInfoChanges) {
        payload.userInfo = buildSeekerInfo(userData);
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
        payload.userInfo = buildCompanyInfo(userData);
      }
    }
  }
  
  return payload;
}

/**
 * Transform backend data to frontend format
 */
export function transformToFrontendFormat(user: User): Record<string, unknown> {
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
    frontendData.linkedIn = '';
    frontendData.desiredRole = '';
    frontendData.aboutMe = '';
    frontendData.dateOfBirth = '';
    frontendData.gender = '';
    frontendData.portfolio = '';
    frontendData.github = '';
    frontendData.skills = [];
    
    // Override with actual values if they exist
    if (user.userInfo) {
      const info = user.userInfo as JobSeekerInfo;
      frontendData.fullName = info.fullName || '';
      frontendData.location = info.location || '';
      frontendData.phone = info.phone || '';
      frontendData.linkedIn = info.linkedIn || '';
      frontendData.desiredRole = info.desiredRole || '';
      frontendData.aboutMe = info.aboutMe || '';
      frontendData.dateOfBirth = info.dateOfBirth || '';
      frontendData.gender = info.gender || '';
      frontendData.portfolio = info.portfolio || '';
      frontendData.github = info.github || '';
      // Parse skills - handle both string and array formats
      if (info.skills) {
        if (typeof info.skills === 'string') {
          frontendData.skills = info.skills.split(',').map(s => s.trim()).filter(s => s.length > 0);
        } else if (Array.isArray(info.skills)) {
          frontendData.skills = info.skills;
        }
      }
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

/**
 * Build seeker userInfo object
 */
function buildSeekerInfo(userData: Record<string, unknown>): JobSeekerInfo {
  return {
    fullName: (userData.fullName as string) || '',
    location: (userData.location as string) || '',
    phone: (userData.phone as string) || '',
    linkedIn: (userData.linkedIn as string) || '', // Fixed: use linkedIn not linkedin
    desiredRole: (userData.desiredRole as string) || '',
    aboutMe: (userData.aboutMe as string) || '',
    dateOfBirth: (userData.dateOfBirth as string) || '',
    portfolio: (userData.portfolio as string) || '',
    github: (userData.github as string) || '',
    skills: Array.isArray(userData.skills) ? userData.skills.join(', ') : (userData.skills as string) || ''
  };
}

/**
 * Build company userInfo object
 */
function buildCompanyInfo(userData: Record<string, unknown>): CompanyInfo {
  return {
    userID: (userData.userID as string) || '',
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

/**
 * Process backend user data (handle _id mapping and userInfo parsing)
 */
export function processBackendUser(user: any): User {
  // Handle null or undefined user
  if (!user) {
    throw new Error('User data is null or undefined');
  }
  
  // Map _id from backend to id for frontend
  if (user._id) {
    user.id = user._id;
  }
  
  // Parse userInfo if it's an array of key-value pairs (company format)
  if (Array.isArray(user.userInfo)) {
    const userInfoObj: Record<string, any> = {};
    user.userInfo.forEach((item: { Key: string; Value: any }) => {
      userInfoObj[item.Key] = item.Value;
    });
    user.userInfo = userInfoObj;
  }
  
  return user as User;
}