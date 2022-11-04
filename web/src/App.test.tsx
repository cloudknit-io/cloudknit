import App from 'App';
import { configure, shallow } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import React from 'react';
import { MemoryRouter } from 'react-router-dom';

configure({ adapter: new Adapter() });
describe('App', () => {
	it('renders without crashing', () => {
		shallow(
			<MemoryRouter>
				<App />
			</MemoryRouter>
		);
	});
});
