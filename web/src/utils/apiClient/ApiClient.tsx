import axios, { AxiosError, AxiosInstance, AxiosRequestConfig } from 'axios';
import { Response } from 'models/response.models';
import { User } from 'models/user.models';

const baseURL = process.env.REACT_APP_BASE_URL;
export interface ErrorResponse {
	statusCode: number;
    message: string;
    error: string;
}

class ApiClient {
	static client: AxiosInstance;
	static setAuthState: (user?: User) => void;

	static init(): ApiClient {
		const client = axios.create({
			baseURL,
			headers: {
				'Content-Type': 'application/json',
			},
			withCredentials: true,
		});

		const handleError = (error: AxiosError): Promise<string> => {
			if (error.response?.status === 401 || error.response?.status === 403) {
				if (error.config.url !== '/session') {
					window.location.href = '/';
					return Promise.reject('Unauthorized');
				}
			}
			return Promise.reject(error.response?.data);
		};

		client.interceptors.response.use(config => config, handleError);

		ApiClient.client = client;

		return ApiClient;
	}

	static get<T = undefined>(url: string, params?: unknown, config?: AxiosRequestConfig): Promise<Response<T>> {
		return ApiClient.client.get<unknown, Response<T>>(url, { ...config, params });
	}

	static post<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<Response<T>> {
		return ApiClient.client.post<unknown, Response<T>>(url, data, config);
	}

	static put<T = unknown>(url: string, data: unknown, config?: AxiosRequestConfig): Promise<Response<T>> {
		return ApiClient.client.put<unknown, Response<T>>(url, data, config);
	}

	static delete<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<Response<T>> {
		return ApiClient.client.delete<unknown, Response<T>>(url, config);
	}

	static patch<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<Response<T>> {
		return ApiClient.client.patch<unknown, Response<T>>(url, data, config);
	}

	static download(url: string): Promise<void> {
		return new Promise(resolve => {
			window.open(ApiClient.client.defaults.baseURL + url);
			return resolve();
		});
	}

	static upload<T = unknown>(
		formData: FormData,
		onUploadProgress?: (e: ProgressEvent) => void
	): Promise<Response<T>> {
		const config = {
			onUploadProgress,
		};
		return ApiClient.post<T>('/files/upload', formData, config);
	}
}

export default ApiClient;
