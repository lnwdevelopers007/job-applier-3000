/**
 * User API layer - handles all user-related HTTP requests
 */

import { ApiClient } from './client';
import type { User, UserFilters, UpdateUserPayload } from '$lib/types';

export class UserApi {
  constructor(private client: ApiClient) {}

  /**
   * Get all users - GET /users/
   */
  async getAll(): Promise<User[]> {
    const response = await this.client.get<User[]>('/users/');
    return response.data;
  }

  /**
   * Query users with filters - GET /users/query
   */
  async query(filters?: UserFilters): Promise<User[]> {
    const response = await this.client.get<User[]>('/users/query', filters);
    return response.data;
  }

  /**
   * Get a specific user by ID - GET /users/:id
   */
  async getById(id: string): Promise<User> {
    const response = await this.client.get<User>(`/users/${id}`);
    return response.data;
  }

  /**
   * Create a new user - POST /users/
   */
  async create(userData: Partial<User>): Promise<User> {
    const response = await this.client.post<User>('/users/', userData);
    return response.data;
  }

  /**
   * Update a user - PUT /users/:id
   */
  async update(id: string, userData: UpdateUserPayload): Promise<User> {
    const response = await this.client.put<User>(`/users/${id}`, userData);
    return response.data;
  }

  /**
   * Delete a user - DELETE /users/:id
   */
  async delete(id: string): Promise<void> {
    await this.client.delete(`/users/${id}`);
  }

  /**
   * Get current user profile - GET /users/profile
   */
  async getProfile(): Promise<User> {
    const response = await this.client.get<User>('/users/profile');
    return response.data;
  }

  /**
   * Update current user profile - PUT /users/profile
   */
  async updateProfile(userData: UpdateUserPayload): Promise<User> {
    const response = await this.client.put<User>('/users/profile', userData);
    return response.data;
  }

  /**
   * Get users by role - GET /users/role/:role
   */
  async getByRole(role: string): Promise<User[]> {
    const response = await this.client.get<User[]>(`/users/role/${role}`);
    return response.data;
  }
}