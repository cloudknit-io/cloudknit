import { configure, shallow } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import { Dashboard } from 'pages/authorized/dashboard/Dashboard';
import React from 'react';
import { MemoryRouter } from 'react-router-dom';

configure({ adapter: new Adapter() });
describe('Dashboard', () => {
	it('renders without crashing', () => {
		shallow(
			<MemoryRouter>
				<Dashboard />
			</MemoryRouter>
		);
	});
});
