import { writable, get } from 'svelte/store';

// Interface for future backend integration
// interface Bookmark {
//   jobId: string;
//   userId: string;
//   createdAt: Date;
// }

class BookmarkService {
  // Store for bookmarked job IDs
  private bookmarkedJobs = writable<Set<string>>(new Set());
  
  // Initialize bookmarks from localStorage or API
  async initializeBookmarks(userId: string) {
    try {
      // For now, use localStorage. Later this can be replaced with API call
      const stored = localStorage.getItem(`bookmarks_${userId}`);
      if (stored) {
        const jobIds = JSON.parse(stored) as string[];
        this.bookmarkedJobs.set(new Set(jobIds));
      } else {
        this.bookmarkedJobs.set(new Set());
      }
      
      // TODO: When backend is ready, replace with:
      // const response = await fetch(`/api/bookmarks?userId=${userId}`);
      // const bookmarks = await response.json();
      // const jobIds = bookmarks.map((b: Bookmark) => b.jobId);
      // this.bookmarkedJobs.set(new Set(jobIds));
    } catch (error) {
      console.error('Failed to initialize bookmarks:', error);
      this.bookmarkedJobs.set(new Set());
    }
  }
  
  // Toggle bookmark for a job
  async toggleBookmark(jobId: string, userId: string): Promise<boolean> {
    const currentBookmarks = get(this.bookmarkedJobs);
    const isCurrentlyBookmarked = currentBookmarks.has(jobId);
    
    if (isCurrentlyBookmarked) {
      await this.removeBookmark(jobId, userId);
      return false;
    } else {
      await this.addBookmark(jobId, userId);
      return true;
    }
  }
  
  // Add a bookmark
  async addBookmark(jobId: string, userId: string) {
    try {
      // Update store
      this.bookmarkedJobs.update(jobs => {
        jobs.add(jobId);
        return new Set(jobs);
      });
      
      // Persist to localStorage
      this.persistToLocalStorage(userId);
      
      // TODO: When backend is ready, add API call:
      // await fetch('/api/bookmarks', {
      //   method: 'POST',
      //   headers: { 'Content-Type': 'application/json' },
      //   body: JSON.stringify({ jobId, userId })
      // });
      
      return true;
    } catch (error) {
      console.error('Failed to add bookmark:', error);
      // Rollback on error
      this.bookmarkedJobs.update(jobs => {
        jobs.delete(jobId);
        return new Set(jobs);
      });
      return false;
    }
  }
  
  // Remove a bookmark
  async removeBookmark(jobId: string, userId: string) {
    try {
      // Update store
      this.bookmarkedJobs.update(jobs => {
        jobs.delete(jobId);
        return new Set(jobs);
      });
      
      // Persist to localStorage
      this.persistToLocalStorage(userId);
      
      // TODO: When backend is ready, add API call:
      // await fetch(`/api/bookmarks/${jobId}`, {
      //   method: 'DELETE',
      //   headers: { 'Content-Type': 'application/json' },
      //   body: JSON.stringify({ userId })
      // });
      
      return true;
    } catch (error) {
      console.error('Failed to remove bookmark:', error);
      // Rollback on error
      this.bookmarkedJobs.update(jobs => {
        jobs.add(jobId);
        return new Set(jobs);
      });
      return false;
    }
  }
  
  // Check if a job is bookmarked
  isBookmarked(jobId: string): boolean {
    const currentBookmarks = get(this.bookmarkedJobs);
    return currentBookmarks.has(jobId);
  }
  
  // Get all bookmarked job IDs
  getBookmarkedJobIds(): string[] {
    const currentBookmarks = get(this.bookmarkedJobs);
    return Array.from(currentBookmarks);
  }
  
  // Subscribe to bookmark changes
  subscribe(callback: (value: Set<string>) => void) {
    return this.bookmarkedJobs.subscribe(callback);
  }
  
  // Clear all bookmarks (useful for logout)
  clearBookmarks() {
    this.bookmarkedJobs.set(new Set());
  }
  
  // Persist to localStorage
  private persistToLocalStorage(userId: string) {
    const currentBookmarks = get(this.bookmarkedJobs);
    const jobIds = Array.from(currentBookmarks);
    localStorage.setItem(`bookmarks_${userId}`, JSON.stringify(jobIds));
  }
}

// Export singleton instance
export const bookmarkService = new BookmarkService();