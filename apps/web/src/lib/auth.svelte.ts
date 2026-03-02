import { goto } from '$app/navigation';

export interface User {
	email: string;
	firstName: string;
	lastName: string;
	isAdmin: boolean;
}

class AuthState {
	user = $state<User | null>(null);
	initialized = $state(false);
	isLoggedIn = $derived(this.user !== null);

	async check() {
		try {
			// Your actual API call to backend
			const response = await fetch('/api/me');
			if (response.ok) {
				this.user = await response.json();
			} else {
				this.user = null;
			}
		} catch {
			this.user = null;
		} finally {
			this.initialized = true;
		}
	}

	async logout() {
		try {
			const response = await fetch('/api/logout', { method: 'post' });
		} finally {
			this.user = null;
		}

		goto('/login');
	}
}

export const auth = new AuthState();
// $inspect(auth.user)
