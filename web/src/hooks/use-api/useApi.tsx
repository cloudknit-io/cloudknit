import { Response, ResponseError } from 'models/response.models';
import { useCallback, useEffect, useState } from 'react';

export type Params = any;
export type Body = any;

export type ApiCallType<T = any, Params = any, Body = any> = (params?: Params, body?: Body) => Promise<Response<T>>;

type ReturnType<T, Params = null, Body = null> = {
	data: T;
	error?: string;
	isLoading: boolean;
	fetch: ApiCallType<T, Params, Body>;
	setData: (data: T) => void;
	setLoading: (isLoading: boolean) => void;
};

type State<T> = {
	isLoading: boolean;
	error?: string;
	response?: Response<T>;
};

export const useApi = <T extends unknown>(
	api: ApiCallType<T, Params, Body>,
	params?: Params,
	body?: Body,
	skipOnLoad = true
): ReturnType<T, Params, Body> => {
	const [state, setState] = useState<State<T>>({
		isLoading: !skipOnLoad,
		error: undefined,
		response: undefined,
	});

	const fetch = useCallback(
		(paramsLoc?: Params, bodyLoc?: Body): Promise<any> => {
			setState(prevState => ({ ...prevState, isLoading: true, error: undefined }));
			return api(paramsLoc, bodyLoc)
				.then((responseLoc: Response<T>) => {
					responseLoc.data = responseLoc.data || ({} as T);
					setState({ isLoading: false, error: undefined, response: responseLoc });
					return new Promise<Response<T>>(resolve => resolve(responseLoc));
				})
				.catch((errorLoc: ResponseError) => {
					const errorMessage = errorLoc ? errorLoc?.data?.message : errorLoc || '';
					setState({ isLoading: false, error: errorMessage, response: undefined });
					return new Promise<T>(resolve => resolve({ error: errorMessage } as T));
				});
		},
		[api]
	);

	const setLoading = (isLoadingLoc: boolean): void =>
		setState(prevState => ({ ...prevState, isLoading: isLoadingLoc }));

	const setData = (dataLoc: T): void =>
		setState(prevState => ({ ...prevState, response: { ...prevState.response, data: dataLoc } as Response<T> }));

	useEffect(() => {
		if (!skipOnLoad) {
			fetch(params, body);
		}
	}, [fetch, params, body, skipOnLoad]);

	const { response, error, isLoading } = state;
	const data = response?.data || undefined;
	return {
		data: data as T,
		error,
		isLoading,
		fetch,
		setLoading,
		setData,
	};
};
