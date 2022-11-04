import { configure, shallow } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import { Login } from 'pages/anonymous/login/Login';
import React from 'react';
import { MemoryRouter } from 'react-router-dom';

configure({ adapter: new Adapter() });
describe('Login', () => {
	it('renders without crashing', () => {
		shallow(
			<MemoryRouter>
				<Login />
			</MemoryRouter>
		);
	});
});
