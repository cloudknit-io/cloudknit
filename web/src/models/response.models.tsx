import { AxiosResponse } from 'axios';

export type ServerError = {
	message: string;
	code: number;
};

export interface Response<T = unknown> extends AxiosResponse<T> {
	error?: ServerError;
}

export type ResponseError = AxiosResponse<ServerError>;

export interface Created {
	id: string;
}
