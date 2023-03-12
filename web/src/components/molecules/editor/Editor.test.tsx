import Editor from '@monaco-editor/react';
import { configure, shallow } from 'enzyme';
import Adapter from '@cfaester/enzyme-adapter-react-18';
import React from 'react';

configure({ adapter: new Adapter() });
describe('Editor', () => {
	it('renders without crashing', () => {
		shallow(<Editor />);
	});
});
