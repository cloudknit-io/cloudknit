import { ZCardLabel } from 'components/atoms/card-label/CardLabel';
import { configure, mount, shallow } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import React from 'react';

configure({ adapter: new Adapter() });
describe('CardLabel', () => {
	it('renders without crashing', () => {
		shallow(<ZCardLabel text="DEMO" color="orange" />);
	});

	it('contains text', () => {
		const wrapper = shallow(<ZCardLabel text="DEMO" color="orange" />);
		expect(wrapper.find('div').text()).toEqual('DEMO');
	});

	it('contains all passed properties', () => {
		const wrapper = mount(<ZCardLabel text="DEMO" color="orange" />);
		expect(wrapper.props()).toEqual({
			color: 'orange',
			text: 'DEMO',
		});
	});
});
