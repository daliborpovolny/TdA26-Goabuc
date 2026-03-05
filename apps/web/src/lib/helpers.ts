
	export function formatTime(iso: string) {
		const date = new Date(iso);

		const formatted = date.toLocaleString();
		return formatted;
	}