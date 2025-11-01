export const DEFAULT_COMPANY_LOGO =
  'https://images.unsplash.com/photo-1534237710431-e2fc698436d0?fm=jpg&q=60&w=3000&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8M3x8YnVpbGRpbmd8ZW58MHx8MHx8fDA%3D';

export const DEFAULT_COMPANY_NAME = 'Unknown Company';

async function fetchThing(endpoint: string, params: string = "") {
  const res = await fetch(`/${endpoint}/${params}`, {
    credentials: 'include'
  });
  if (!res.ok) throw new Error(`Failed to fetch ${endpoint} details (Status: ${res.status})`);
  const data = await res.json();
  return data;
}


export async function fetchJob(jobID: string) {
  return fetchThing('jobs', jobID)
}

export async function fetchUser(userID: string) {
  return fetchThing('users', userID)
}

export async function fetchUsers(params?: string) {
  return fetchThing('users')
}

export async function fetchCompanyNameLogo(companyID: any) {
  let companyName = DEFAULT_COMPANY_NAME
  let companyLogo = DEFAULT_COMPANY_LOGO

  if (!companyID) {
    return [companyName, companyLogo];
  }

  try {
    const companyData = await fetchUser(companyID);
    [companyName, companyLogo] = getCompanyInfo(companyData);
  } catch (err) {
    // Only log warning for non-404 errors to reduce noise
    if (err instanceof Error && !err.message.includes('Company not found')) {
      console.warn(`Failed to load company info for ID ${companyID}:`, err);
    }
  }

  return [companyName, companyLogo];
}

function getCompanyInfo(company: any) {

  const infoArray = company.userInfo || [];
  const info = Object.fromEntries(infoArray.map((item: any) => [item.Key, item.Value]));

  const companyName: string = info.name || company.name || DEFAULT_COMPANY_NAME;
  const companyLogo: string = info.logo || company.avatarURL || DEFAULT_COMPANY_LOGO;
  return [companyName, companyLogo]
}
