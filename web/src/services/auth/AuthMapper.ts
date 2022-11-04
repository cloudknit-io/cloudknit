import { Response } from 'models/response.models';
import { User } from 'models/user.models';

export class AuthMapper {
	static async me(response: Response<User>): Promise<Response<User>> {
		return {
			...response,
			data: {
				...response.data,
				name: response.data.name,
			},
		};
	}
}
