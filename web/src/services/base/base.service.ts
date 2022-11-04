import { debounce } from 'lodash';
import { Subject } from 'rxjs';
import ApiClient from 'utils/apiClient';

export class BaseService {
	protected requestMap = new Map<string, any>();
	protected streamMap = new Map<string, Subject<any>>();
	protected localStorageCache = window.localStorage;
	private debounceTime: number = Number.MAX_SAFE_INTEGER;
	protected cache = new Map<string, any>();
	private cacheKey = '';

	constructor(debouceTime: number, cacheKey: string) {
		this.debounceTime = debouceTime;
		this.cacheKey = cacheKey;
	}

	private createRequest = async (requestor: () => Promise<any>) => {
		try {
			return await requestor();
		} catch (err) {
			return err;
		}
	};

	protected setStreamHandler<T>(url: string, key: string) {
		const subject = new Subject<any>();
		const requestHandler = debounce(
			this.createRequest.bind(null, async () => {
				const { data } = await ApiClient.get<T>(url);
				this.notifySubscribers(key, data, subject);
			}),
			this.debounceTime,
			{
				leading: true,
			}
		);
		this.streamMap.set(key, subject);
		this.requestMap.set(key, requestHandler);
	}

	protected notifySubscribers(key: string, data: any, subject: Subject<any>) {
		this.cache.set(key, data);
		try {
			this.localStorageCache.setItem(this.cacheKey, JSON.stringify(Array.from(this.cache.entries())));
		} catch(err) {
			this.localStorageCache.clear();
		}
		
		subject.next(data);
	}

	protected getStream<T>(key: string, url: string) {
		if (!key || !url) {
			throw 'Key or url cannot be empty';
		}
		if (!this.streamMap.has(key)) {
			this.setStreamHandler<T>(url, key);
		}
		this.requestMap.get(key)();
		return this.streamMap.get(key) as Subject<any>;
	}

	getCachedValue(key?: string) {
		if (!key) {
			throw 'Key cannot be empty';
		}

		if (this.cache.has(key)) {
			return this.cache.get(key);
		}

		if (this.cache.size === 0 && this.localStorageCache.getItem(this.cacheKey)) {
			this.cache = new Map<string, any>(JSON.parse(this.localStorageCache.getItem(this.cacheKey) || ''));
		}
		if (this.cache.has(key)) {
			return this.cache.get(key);
		}
		return null;
	}
}
