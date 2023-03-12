import { configure, shallow } from 'enzyme';
import Adapter from '@cfaester/enzyme-adapter-react-18';
import { Environments } from 'pages/authorized/environments/Environments';
import React from 'react';
import { MemoryRouter } from 'react-router-dom';

configure({ adapter: new Adapter() });
describe('Environments', () => {
	it('renders without crashing', () => {
		shallow(
			<MemoryRouter>
				<Environments />
			</MemoryRouter>
		);
	});
});
