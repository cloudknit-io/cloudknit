import { configure } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import { LocalStorageKey } from 'models/localStorage';
import { LocalStorage } from 'utils/localStorage/localStorage';

configure({ adapter: new Adapter() });
describe('LocalStorage', () => {
	it('should fail to retrieve from localStorage', () => {
		expect(LocalStorage.getItem(LocalStorageKey.USER)).not.toEqual('DEMO');
	});

	it('stores to localStorage', () => {
		LocalStorage.setItem(LocalStorageKey.USER, 'DEMO');
		expect(LocalStorage.getItem(LocalStorageKey.USER)).toEqual('DEMO');
	});

	it('should fail to retrieve from localStorage after item ise removed', () => {
		LocalStorage.setItem(LocalStorageKey.USER, 'DEMO');
		LocalStorage.removeItem(LocalStorageKey.USER);
		expect(LocalStorage.getItem(LocalStorageKey.USER)).not.toEqual('DEMO');
	});
});
