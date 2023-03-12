import { configure, shallow } from 'enzyme';
import Adapter from '@cfaester/enzyme-adapter-react-18';
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
