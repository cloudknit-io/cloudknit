import { configure, shallow } from 'enzyme';
import Adapter from '@cfaester/enzyme-adapter-react-18';
import { NotFound } from 'pages/anonymous/not-found/NotFound';
import React from 'react';
import { MemoryRouter } from 'react-router-dom';

configure({ adapter: new Adapter() });
describe('NotFound', () => {
	it('renders without crashing', () => {
		shallow(
			<MemoryRouter>
				<NotFound />
			</MemoryRouter>
		);
	});
});
