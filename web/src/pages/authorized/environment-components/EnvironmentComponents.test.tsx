import { configure, shallow } from 'enzyme';
import Adapter from '@cfaester/enzyme-adapter-react-18';
import { EnvironmentComponents } from 'pages/authorized/environment-components/EnvironmentComponents';
import React from 'react';
import { MemoryRouter } from 'react-router-dom';

configure({ adapter: new Adapter() });
describe('EnvironmentComponents', () => {
	it('renders without crashing', () => {
		shallow(
			<MemoryRouter>
				<EnvironmentComponents />
			</MemoryRouter>
		);
	});
});
