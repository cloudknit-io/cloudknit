import { User } from 'models/user.models';

export interface LoginFormValues {
	email: string;
	password: string;
	verification_url?: string;
}

export interface RegisterFormValues extends LoginFormValues {
	terms: boolean;
	name: string;
	verifyPassword: string;
	token?: string;
	is_organisation?: boolean;
}

export type ForgotPasswordValues = {
	email: string;
	url: string;
};

export interface ResetPasswordValues {
	password: string;
	verifyPassword: string;
	token: string;
}

export interface AuthState {
	user: User;
	setAuthState: (user?: User) => void;
}
