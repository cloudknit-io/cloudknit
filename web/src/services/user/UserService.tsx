import { ForgotPasswordValues, RegisterFormValues, ResetPasswordValues } from 'models/auth.models';
import { Created, Response } from 'models/response.models';
import { User } from 'models/user.models';
import ApiClient from 'utils/apiClient';

export class UserService {
	static createClient(data: RegisterFormValues): Promise<Response<Created>> {
		return ApiClient.post<Created>('/users/client', data);
	}

	static createExpert(data: RegisterFormValues): Promise<Response<Created>> {
		return ApiClient.post<Created>('/users/expert', data);
	}

	static async createAndLoginIfExpert(data: RegisterFormValues): Promise<Response<User>> {
		const action = data.is_organisation ? UserService.createExpert : UserService.createClient;
		const register = await action(data);
		if (!data.is_organisation) {
			return {
				...register,
				data: {} as User,
			};
		}
		return { data: {} } as Response<User>;
	}

	static forgotPassword(data: ForgotPasswordValues): Promise<Response<Created>> {
		return ApiClient.post<Created>('/users/forgot/password', data);
	}

	static updatePassword(data: ResetPasswordValues): Promise<Response> {
		return ApiClient.post('/users/forgot/password/update', data);
	}
}
