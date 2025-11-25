interface User {
	email: string;
	firstName: string;
	lastName: string;
}

let user = $state<User | null>(null);

export function getUser() {
	return user;
}

export function setUser(newUser: User | null) {
	user = newUser;
}

export function isLoggedIn() {
	return isUserLoggedIn;
}

const isUserLoggedIn = $derived(!!user);

export async function checkSession(
	fetch: (info: RequestInfo, init?: RequestInit) => Promise<Response>
) {
	try {
		const response = await fetch('/api/me');

		if (response.ok) {
			const userData = await response.json();
			user = userData;
			return userData;
		} else {
			// If the server returns a 401/403, session is expired
			console.log('session expired i guess?');
			user = null;
		}
	} catch (error) {
		// Catches network errors
		console.error('Session check failed:', error);
		user = null;
	}
	return null;
}
