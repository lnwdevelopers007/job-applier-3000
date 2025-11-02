/**
 * User-related type definitions
 */

// Backend User schema
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
  skills?: string | string[];
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

export interface UserFilters {
  id?: string;
  role?: string;
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

// Frontend data interfaces for settings
export interface BaseUserData {
  id?: string;
  name?: string;
  email?: string;
  avatar?: string;
  provider?: string;
  userID?: string;
  role?: string;
  verified?: boolean;
  googleConnected?: boolean;
  documents?: unknown[];
}

export interface SeekerUserData extends BaseUserData {
  fullName?: string;
  location?: string;
  phone?: string;
  linkedin?: string;
  desiredRole?: string;
  aboutMe?: string;
  dateOfBirth?: string;
  portfolio?: string;
  github?: string;
  skills?: string[];
}

export interface CompanyUserData extends BaseUserData {
  companyName?: string;
  aboutCompany?: string;
  industry?: string;
  companySize?: string;
  companyWebsite?: string;
  companyLogo?: string;
  foundedYear?: string;
  headquarters?: string;
  companyLinkedin?: string;
}

export type UserData = SeekerUserData | CompanyUserData;

export interface Tab {
  id: string;
  label: string;
  title: string;
  description: string;
}