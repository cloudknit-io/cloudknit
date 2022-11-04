import { Loader } from 'components/atoms/loader/Loader';
import { configure, mount, shallow } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import React from 'react';

configure({ adapter: new Adapter() });
describe('Loader', () => {
	it('renders without crashing', () => {
		shallow(<Loader />);
	});

	it('is white', () => {
		const wrapper = shallow(<Loader color="#ffffff" />);
		expect(wrapper.find('circle').props().stroke).toEqual('#ffffff');
	});

	it('contains all passed properties', () => {
		const wrapper = mount(<Loader width={23} height={23} color="#ffffff" />);
		expect(wrapper.props()).toEqual({
			color: '#ffffff',
			width: 23,
			height: 23,
		});
	});
});
