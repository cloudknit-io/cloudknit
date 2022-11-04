import { LocalStorageKey } from 'models/localStorage';

export class LocalStorage {
	static getItem<T>(key: LocalStorageKey): T | null {
		const value = window.localStorage.getItem(key);
		if (value) {
			return JSON.parse(value) as T;
		}
		return null;
	}

	static setItem<T>(key: LocalStorageKey, value: T | null): void {
		try {
			window.localStorage.setItem(key, JSON.stringify(value));
		} catch(err) {
			window.localStorage.clear();
		}
		
	}

	static removeItem(key: LocalStorageKey): void {
		window.localStorage.removeItem(key);
	}
}
