import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { jwtDecode } from 'jwt-decode';

type JWTPayload = {
  banned?: boolean;
  role?: string;
};

export const load: PageServerLoad = async ({ cookies }) => {
  const accessToken = cookies.get('access_token');
  
  if (accessToken) {
    try {
      const decoded = jwtDecode<JWTPayload>(accessToken);
      
      // If user is not banned, redirect them away
      if (!decoded.banned) {
        // Redirect to their appropriate dashboard
        const role = decoded.role?.toLowerCase();
        if (role === 'company') {
          throw redirect(303, '/company/dashboard');
        } else if (role === 'faculty') {
          throw redirect(303, '/app/jobs');
        } else if (role === 'admin') {
          throw redirect(303, '/admin/dashboard');
        } else {
          throw redirect(303, '/app/jobs');
        }
      }
    } catch (error) {
      // If JWT decode fails, allow access to banned page
      // (user might be manually navigating here, or token is invalid)
      console.error('JWT decode error on banned page:', error);
    }
  }
  // If no token or user is actually banned, allow access to the page
  return {};
};
