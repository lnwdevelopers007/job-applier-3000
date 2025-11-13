import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { jwtDecode } from 'jwt-decode';

type JWTPayload = {
  banned?: boolean;
  role?: string;
  exp?: number;
};

export const load: PageServerLoad = async ({ cookies, locals }) => {
  console.log('=== BANNED PAGE LOAD ===');
  console.log('Has access_token:', !!cookies.get('access_token'));
  console.log('Has refresh_token:', !!cookies.get('refresh_token'));
  
  const accessToken = cookies.get('access_token');
  const refreshToken = cookies.get('refresh_token');
  
  // Case 1: User has access token
  if (accessToken) {
    console.log('Case 1: Has access token, checking...');
    try {
      const decoded = jwtDecode<JWTPayload>(accessToken);
      console.log('JWT banned status:', decoded.banned);
      
      // Trust the JWT for banned status
      if (decoded.banned === true) {
        // User is banned - clear cookies (logout) and show banned page
        console.log('User is banned, clearing cookies and showing banned page');
        cookies.delete('access_token', { path: '/' });
        cookies.delete('refresh_token', { path: '/' });
        
        return { 
          user: null,
          wasBanned: true 
        };
      } else {
        // User is NOT banned - redirect to dashboard
        console.log('User is NOT banned, redirecting to dashboard');
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
      console.error('Error in Case 1:', error);
      if (error && typeof error === 'object' && 'status' in error) {
        throw error;
      }
    }
  }
  
  // Case 2: No access token but has refresh token - try to refresh
  console.log('Case 2: No access token, trying refresh token...');
  if (refreshToken) {
    try {
      const backendUrl = import.meta.env.VITE_BACKEND || 'http://localhost:8080';
      console.log('Calling refresh endpoint...');
      const refreshResponse = await fetch(`${backendUrl}/auth/refresh`, {
        method: 'POST',
        headers: {
          'Cookie': `refresh_token=${refreshToken}`
        },
        credentials: 'include'
      });

      console.log('Refresh response status:', refreshResponse.status);

      if (refreshResponse.ok) {
        const data = await refreshResponse.json();
        const newDecoded = jwtDecode<JWTPayload>(data.access_token);
        
        console.log('New token decoded, banned status:', newDecoded.banned);
        
        // Check ban status from NEW token
        if (newDecoded.banned === true) {
          console.log('Newly refreshed token shows user is banned - clearing all tokens');
          // Clear both tokens - don't set the new one
          cookies.delete('access_token', { path: '/' });
          cookies.delete('refresh_token', { path: '/' });
          return { 
            user: null,
            wasBanned: true 
          };
        } else {
          // Not banned, set token and redirect to dashboard
          console.log('User is NOT banned, setting token and redirecting');
          const jwtExpiresIn = newDecoded.exp 
            ? Math.max(0, newDecoded.exp - Math.floor(Date.now() / 1000)) 
            : 60 * 60 * 24;
          
          cookies.set('access_token', data.access_token, {
            path: '/',
            httpOnly: true,
            secure: import.meta.env.MODE === 'production',
            sameSite: 'lax',
            maxAge: jwtExpiresIn
          });
          
          const role = newDecoded.role?.toLowerCase();
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
      } else {
        // Refresh failed - check if it's due to ban
        console.log('Refresh failed, checking error...');
        const errorData = await refreshResponse.json().catch(() => ({}));
        console.log('Refresh error data:', errorData);
        
        if (errorData.error === 'account_banned') {
          // User is banned - clear tokens and show banned page
          console.log('Refresh failed due to ban, clearing tokens');
          cookies.delete('access_token', { path: '/' });
          cookies.delete('refresh_token', { path: '/' });
          return { 
            user: null,
            wasBanned: true 
          };
        }
        
        // Other error - redirect to login
        console.log('Refresh failed for other reason, redirecting to login');
        cookies.delete('access_token', { path: '/' });
        cookies.delete('refresh_token', { path: '/' });
        throw redirect(303, '/login');
      }
    } catch (error) {
      console.error('Error in Case 2:', error);
      if (error && typeof error === 'object' && 'status' in error) {
        throw error;
      }
      cookies.delete('access_token', { path: '/' });
      cookies.delete('refresh_token', { path: '/' });
      throw redirect(303, '/login');
    }
  }
  
  // Case 3: No tokens at all - someone manually visited /banned or was logged out
  // Allow them to see the page with a generic message
  console.log('Case 3: No tokens at all, showing banned page');
  return { 
    user: null,
    wasBanned: false // Might be someone just visiting the URL
  };
};
