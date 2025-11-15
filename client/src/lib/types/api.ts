/**
 * Generic API-related type definitions
 */

export interface ApiResponse<T = any> {
  data: T;
  status: number;
  ok: boolean;
}

export interface ApiError {
  message: string;
  status: number;
  code?: string;
}

