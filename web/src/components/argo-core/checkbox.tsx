import * as React from 'react';

export const Checkbox = (props: {
	disabled?: boolean;
	checked?: boolean;
	onChange?: (val: boolean) => any;
	onNativeChange?: (e: React.ChangeEvent<HTMLInputElement>) => any;
	id?: string;
	value?: string;
}) => (
	<span className="argo-checkbox">
		<input
			id={props.id}
			type="checkbox"
			disabled={props.disabled}
			checked={props.checked}
			data-value={props.value}
			onChange={e =>
				(props.onNativeChange && props.onNativeChange(e) && true) ||
				(props.onChange && props.onChange(!props.checked))
			}
		/>
		<span>
			<i className="fa fa-check" />
		</span>
	</span>
);
