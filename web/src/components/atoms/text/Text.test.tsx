import { ZText } from 'components/atoms/text/Text';
import { configure, mount, shallow } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import React from 'react';

configure({ adapter: new Adapter() });
describe('Text', () => {
	it('renders without crashing', () => {
		shallow(<ZText>DEMO</ZText>);
	});

	it('contains text with element p', () => {
		const wrapper = shallow(<ZText.Body>DEMO</ZText.Body>);
		expect(wrapper.find('p').text()).toEqual('DEMO');
	});

	it('contains text with element h2', () => {
		const wrapper = shallow(<ZText.Headline>DEMO</ZText.Headline>);
		expect(wrapper.find('h2').text()).toEqual('DEMO');
	});

	it('contains all passed properties', () => {
		const wrapper = mount(
			<ZText.Body size="24" lineHeight="16" upperCase weight="bold">
				DEMO
			</ZText.Body>
		);
		expect(wrapper.props()).toEqual({
			children: 'DEMO',
			size: '24',
			lineHeight: '16',
			upperCase: true,
			weight: 'bold',
		});
	});
});
