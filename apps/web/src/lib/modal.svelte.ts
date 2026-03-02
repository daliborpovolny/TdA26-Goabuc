class ModalService {
	show = $state(false);
	message = $state('');
	resolve = $state<(value: boolean) => void>(() => {});

	async confirm(msg: string): Promise<boolean> {
		this.message = msg;
		this.show = true;

		return new Promise((res) => {
			this.resolve = res;
		});
	}

	close(result: boolean) {
		this.show = false;
		this.resolve(result);
	}
}

export const modal = new ModalService();
