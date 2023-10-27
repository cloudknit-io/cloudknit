import { LocalStorageKey } from 'models/localStorage';
import { Organization, User } from 'models/user.models';
import { AuthService } from 'services/auth/AuthService';
import { LocalStorage } from 'utils/localStorage/localStorage';

class AuthStore {
	user: User | null;

	constructor() {
		this.user = null;
	}

	login(): void {
		window.location.href = `${process.env.REACT_APP_AUTHORIZE_URL}`;
	}

	logoutUrl(): string {
		return process.env.REACT_APP_LOGOUT_URL || "/auth/logout";
	}

	logout(): void {
		window.location.href = this.logoutUrl();
	}

	redirectToHome() {
		window.location.href = `${process.env.REACT_APP_BASE_URL}`;
	}

	async refresh(): Promise<User> {
		const { data } = await AuthService.me();

		this.user = data;
		const selectedOrg = this.getOrganization();
		this.user.selectedOrg = selectedOrg as Organization;
		LocalStorage.setItem(LocalStorageKey.USER, this.user);

		return this.user;
	}

	async fetchOrganization(id?: number) {
		if (!id) throw 'Not a valid org id';
		const user = await this.refresh();
		const org = user.organizations.find(org => org.id === id);
		this.patchOrgToUser(org?.name);
		return org;
	}

	getUser(): User | null {
		return this.user || LocalStorage.getItem<User>(LocalStorageKey.USER);
	}

	async selectOrganization(orgName?: string, refresh: boolean = true) {
		if (!orgName) return;
		await AuthService.selectOrganization(orgName);
		this.patchOrgToUser(orgName);
		if (refresh) {
			this.redirectToHome();
		}
	}

	getOrganization() {
		const selectedOrg = this.getUser()?.selectedOrg || LocalStorage.getItem(LocalStorageKey.SELECTED_ORG);
		if (selectedOrg && this.getUser()?.organizations.some(e => e.name === selectedOrg.name)) {
			return selectedOrg;
		}
		return null;
	}

	async fetchOrganizationStatus() {
		const { data } = await AuthService.fetchOrganizationStatus();
		return data;
	}

	async addOrganization(orgName: string) {
		// Add organization to the list and set it in the local storage
		const { data } = await AuthService.addOrganization(orgName);
		this.user?.organizations.push(data);
		this.patchOrgToUser(data.name);
		return this.getOrganization();
	}

	private patchOrgToUser(orgName?: string) {
		if (this.user && orgName) {
			const org = this.getUser()?.organizations.find(e => e.name === orgName);
			if (!org) return;
			this.user = {
				...this.user,
				selectedOrg: org,
			};
		}
		// Updating user and selected org.
		if (this.user?.selectedOrg) {
			LocalStorage.setItem(LocalStorageKey.SELECTED_ORG, this.user?.selectedOrg);
		}
		LocalStorage.setItem(LocalStorageKey.USER, this.user);
	}
}

export default new AuthStore();
