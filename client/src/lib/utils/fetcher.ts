export const DEFAULT_COMPANY_LOGO =
  'https://images.unsplash.com/photo-1534237710431-e2fc698436d0?fm=jpg&q=60&w=3000&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8M3x8YnVpbGRpbmd8ZW58MHx8MHx8fDA%3D';

export const DEFAULT_COMPANY_NAME = 'Unknown Company';



export async function fetchJob(jobID: string) {
  const jobRes = await fetch(`/jobs/${jobID}`, {
    credentials: 'include'
  });
  if (!jobRes.ok) throw new Error('Failed to fetch job details');
  const jobData = await jobRes.json();
  return jobData;
}

export async function fetchCompany(companyID: string) {
  const companyRes = await fetch(`/users/${companyID}`, {
    credentials: 'include'
  });
  if (!companyRes.ok) throw new Error('Failed to fetch company details');
  const companyData = await companyRes.json()
  return companyData
}

export async function fetchCompanyNameLogo(companyID: any) {

  let companyName = DEFAULT_COMPANY_NAME
  let companyLogo = DEFAULT_COMPANY_LOGO
  try {
    const companyData = await fetchCompany(companyID);
    [companyName, companyLogo] = getCompanyInfo(companyData);
  } catch (err) {
    console.warn(`Failed to load company info for ID ${companyID}:`, err);
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