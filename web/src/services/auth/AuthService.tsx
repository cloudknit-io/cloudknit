import { LoginFormValues } from 'models/auth.models';
import { Response } from 'models/response.models';
import { Organization, User } from 'models/user.models';
import { AuthMapper } from 'services/auth/AuthMapper';
import ApiClient from 'utils/apiClient';

export class AuthService {
	static login(credentials: LoginFormValues): Promise<Response<User>> {
		return ApiClient.post<User>('/auth/token', credentials);
	}

	static me(): Promise<Response<User>> {
		return ApiClient.get<User>('/session').then(AuthMapper.me);
	}

	static logout() {
		const logoutUrl = `${process.env.LOGOUT_URL}`;
		window.location.href = logoutUrl;
	}

	static async loginAndProfile(): Promise<Response<User>> {
		return { data: {} } as Response<User>;
	}

	static verify(token: string): Promise<Response> {
		return ApiClient.post(`/token/${token}`);
	}

	static addOrganization(orgName: string) {
		return ApiClient.post<Organization>(`/registration/orgs`, {
			name: orgName,
			githubRepo: '',
		});
	}

	static selectOrganization(orgName: string) {
		return ApiClient.post<Organization>(`/auth/select-org`, {
			selectOrg: orgName,
		});
	}

	static fetchOrganization(orgId: number) {
		return ApiClient.get<Organization>(`/orgs/${orgId}`);
	}

	static fetchOrganizationStatus() {
		return ApiClient.get<Organization>('/ops/is-provisioned');
	}
}
