import { Button } from 'components/atoms/button/Button';
import { Loader } from 'components/atoms/loader/Loader';
import { configure, mount, shallow } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import React from 'react';

configure({ adapter: new Adapter() });
describe('Button', () => {
	it('renders without crashing', () => {
		shallow(<Button>DEMO</Button>);
	});

	it('contains text', () => {
		const wrapper = shallow(<Button>DEMO</Button>);
		expect(wrapper.find('button').text()).toEqual('DEMO');
	});

	it('contains all passed properties', () => {
		const wrapper = mount(
			<Button disabled block isLoading>
				DEMO
			</Button>
		);
		expect(wrapper.props()).toEqual({
			disabled: true,
			block: true,
			isLoading: true,
			children: 'DEMO',
		});
	});

	it('triggers function on click', () => {
		const mockCallBack = jest.fn();

		const wrapper = shallow(<Button onClick={mockCallBack}>DEMO</Button>);
		wrapper.find('button').simulate('click');
		expect(mockCallBack.mock.calls.length).toEqual(1);
	});

	it('have loading', () => {
		const wrapper = shallow(<Button isLoading>DEMO</Button>);
		expect(wrapper.containsMatchingElement(<Loader width={20} height={20} color="#ffffff" />)).toEqual(true);
	});
});
