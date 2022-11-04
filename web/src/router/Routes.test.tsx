import { configure, shallow } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import React from 'react';
import { MemoryRouter } from 'react-router-dom';
import Routes from 'router/Routes';

configure({ adapter: new Adapter() });
describe('Routes', () => {
	it('renders without crashing', () => {
		shallow(
			<MemoryRouter>
				<Routes />
			</MemoryRouter>
		);
	});
});
