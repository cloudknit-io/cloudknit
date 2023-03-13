import { configure, shallow } from 'enzyme';
import Adapter from '@cfaester/enzyme-adapter-react-18';
import { Profile } from 'pages/authorized/profile/Profile';
import React from 'react';
import { MemoryRouter } from 'react-router-dom';

configure({ adapter: new Adapter() });
describe('Profile', () => {
	it('renders without crashing', () => {
		shallow(
			<MemoryRouter>
				<Profile />
			</MemoryRouter>
		);
	});
});
