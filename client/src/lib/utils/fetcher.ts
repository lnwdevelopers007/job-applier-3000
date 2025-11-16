import { apiFetch } from "./api";

export const DEFAULT_COMPANY_LOGO =
  'https://images.unsplash.com/photo-1534237710431-e2fc698436d0?fm=jpg&q=60&w=3000&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8M3x8YnVpbGRpbmd8ZW58MHx8MHx8fDA%3D';

export const DEFAULT_COMPANY_NAME = 'Unknown Company';

export async function fetchData(endpoint: string, params: string = "") {
  const res = await apiFetch(`/${endpoint}/${params}`)
  if (!res.ok) throw new Error(`Failed to fetch ${endpoint} details (Status: ${res.status})`);
  const data = await res.json();
  return data;
}


export async function fetchJob(jobID: string) {
  return fetchData('jobs', jobID)
}

export async function fetchUser(userID: string) {
  return fetchData('users', userID)
}

export async function fetchUsers(params?: string) {
  return fetchData('users', params)
}

/**
 * Fetch company name and logo (with custom logo support)
 * @deprecated Use fetchCompanyPublicInfo instead - it now handles custom logos automatically
 */
export async function fetchCompanyNameLogo(companyID: any) {
  let companyName = DEFAULT_COMPANY_NAME
  let companyLogo = DEFAULT_COMPANY_LOGO

  if (!companyID) {
    return [companyName, companyLogo];
  }

  try {
    const companyData = await fetchCompanyPublicInfo(companyID);
    companyName = companyData.name || DEFAULT_COMPANY_NAME;
    companyLogo = companyData.profileImage || DEFAULT_COMPANY_LOGO;
  } catch (err) {
    // Only log warning for non-404 errors to reduce noise
    if (err instanceof Error && !err.message.includes('not found')) {
      console.warn(`Failed to load company info for ID ${companyID}:`, err);
    }
  }

  return [companyName, companyLogo];
}

/**
 * Fetch public company information (name, role, logo)
 * This endpoint now automatically returns custom logo from userInfo if available,
 * falling back to avatarURL if not set
 */
export async function fetchCompanyPublicInfo(userID: string): Promise<{
  name: string;
  role: string;
  profileImage: string;
  userInfo: {
    name: string;
    logo: string;
  };
}> {
  // Use relative URL to go through Vite proxy
  const res = await fetch(`/users/public/${userID}`);

  if (!res.ok) {
    throw new Error(`Failed to fetch public user info for ${userID}`);
  }

  const data = await res.json();

  // Combine top-level and nested userInfo for maximum compatibility
  const userInfo = {
    name: data.userInfo?.name || "",
    logo: data.userInfo?.logo || "",
  };

  return {
    name: data.name || userInfo.name || DEFAULT_COMPANY_NAME,
    role: data.role || "company",
    profileImage: data.profileImage || userInfo.logo || DEFAULT_COMPANY_LOGO,
    userInfo,
  };
}
