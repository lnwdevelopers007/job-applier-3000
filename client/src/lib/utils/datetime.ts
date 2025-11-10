/**
 * Utility functions for formatting dates and times
 */

export interface DateTimeOptions {
	includeTime?: boolean;
	use24Hour?: boolean;
	showSeconds?: boolean;
}

/**
 * Formats a date to DD/MM/YYYY format
 */
export function formatDateDMY(date: string | Date, options: DateTimeOptions = {}): string {
	const d = new Date(date);
	
	if (isNaN(d.getTime())) {
		return 'Invalid Date';
	}

	const day = d.getDate().toString().padStart(2, '0');
	const month = (d.getMonth() + 1).toString().padStart(2, '0');
	const year = d.getFullYear();
	
	let result = `${day}/${month}/${year}`;
	
	if (options.includeTime) {
		const time = formatTime(d, options.use24Hour, options.showSeconds);
		result += ` ${time}`;
	}
	
	return result;
}

/**
 * Formats time in 12-hour or 24-hour format
 */
function formatTime(date: Date, use24Hour = false, showSeconds = false): string {
	if (use24Hour) {
		const hours = date.getHours().toString().padStart(2, '0');
		const minutes = date.getMinutes().toString().padStart(2, '0');
		const seconds = date.getSeconds().toString().padStart(2, '0');
		
		return showSeconds ? `${hours}:${minutes}:${seconds}` : `${hours}:${minutes}`;
	} else {
		let hours = date.getHours();
		const minutes = date.getMinutes().toString().padStart(2, '0');
		const seconds = date.getSeconds().toString().padStart(2, '0');
		const ampm = hours >= 12 ? 'PM' : 'AM';
		
		hours = hours % 12;
		hours = hours ? hours : 12; // the hour '0' should be '12'
		const hoursStr = hours.toString();
		
		const timeStr = showSeconds ? `${hoursStr}:${minutes}:${seconds}` : `${hoursStr}:${minutes}`;
		return `${timeStr} ${ampm}`;
	}
}

/**
 * Formats a date into human-readable relative time
 */
export function formatRelativeTime(date: string | Date): string {
	const d = new Date(date);
	const now = new Date();
	
	if (isNaN(d.getTime())) {
		return 'Unknown';
	}
	
	const diffTime = now.getTime() - d.getTime();
	const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));
	const diffHours = Math.floor(diffTime / (1000 * 60 * 60));
	const diffMinutes = Math.floor(diffTime / (1000 * 60));
	
	// Future dates
	if (diffTime < 0) {
		const futureDiffDays = Math.abs(diffDays);
		const futureDiffHours = Math.abs(diffHours);
		const futureDiffMinutes = Math.abs(diffMinutes);
		
		if (futureDiffDays > 0) {
			if (futureDiffDays === 1) return 'Tomorrow';
			if (futureDiffDays < 7) return `In ${futureDiffDays} days`;
			if (futureDiffDays < 30) return `In ${Math.floor(futureDiffDays / 7)} weeks`;
			if (futureDiffDays < 365) return `In ${Math.floor(futureDiffDays / 30)} months`;
			return `In ${Math.floor(futureDiffDays / 365)} years`;
		} else if (futureDiffHours > 0) {
			return futureDiffHours === 1 ? 'In 1 hour' : `In ${futureDiffHours} hours`;
		} else if (futureDiffMinutes > 0) {
			return futureDiffMinutes === 1 ? 'In 1 minute' : `In ${futureDiffMinutes} minutes`;
		} else {
			return 'In a moment';
		}
	}
	
	// Past dates
	if (diffDays > 0) {
		if (diffDays === 1) return '1 day ago';
		if (diffDays < 7) return `${diffDays} days ago`;
		if (diffDays < 30) {
			const weeks = Math.floor(diffDays / 7);
			return weeks === 1 ? '1 week ago' : `${weeks} weeks ago`;
		}
		if (diffDays < 365) {
			const months = Math.floor(diffDays / 30);
			if (months === 1) return '1 month ago';
			return diffDays >= 30 ? '30+ days ago' : `${months} months ago`;
		}
		const years = Math.floor(diffDays / 365);
		return years === 1 ? '1 year ago' : `${years} years ago`;
	} else if (diffHours > 0) {
		return diffHours === 1 ? '1 hour ago' : `${diffHours} hours ago`;
	} else if (diffMinutes > 0) {
		return diffMinutes === 1 ? '1 minute ago' : `${diffMinutes} minutes ago`;
	} else {
		return 'Today';
	}
}


/**
 * Formats a date into a short readable format (e.g., "Mar 15, 2024")
 */
export function formatDateShort(date: string | Date): string {
	const d = new Date(date);
	
	if (isNaN(d.getTime())) {
		return 'Invalid Date';
	}
	
	return d.toLocaleDateString('en-US', {
		year: 'numeric',
		month: 'short',
		day: 'numeric'
	});
}

/**
 * Formats a date into format "30 Dec 2024"
 */
export function formatDateCompact(date: string | Date): string {
	const d = new Date(date);
	
	if (isNaN(d.getTime())) {
		return 'Invalid Date';
	}
	
	return d.toLocaleDateString('en-GB', {
		year: 'numeric',
		month: 'short',
		day: 'numeric'
	});
}

/**
 * Get the number of days between two dates
 */
export function getDaysBetween(date1: string | Date, date2: string | Date): number {
	const d1 = new Date(date1);
	const d2 = new Date(date2);
	
	if (isNaN(d1.getTime()) || isNaN(d2.getTime())) {
		return 0;
	}
	
	const diffTime = Math.abs(d2.getTime() - d1.getTime());
	return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
}

/**
 * Formats a date string to include deadline time (16:59:00 UTC)
 */
export function withDeadlineTime(dateStr: string): string {
	const d = new Date(dateStr);
	d.setUTCHours(16, 59, 0, 0);
	return d.toISOString();
}
