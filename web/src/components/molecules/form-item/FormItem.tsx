import 'components/molecules/form-item/FormItem.scss';

import classNames from 'classnames';
import { FormField } from 'components/molecules/form-item/FormField';
import { Field, FieldProps } from 'formik';
import React, { ReactNode } from 'react';

export type InputSizes = 'sm' | 'lg';

export interface FormItemProps {
	className?: string;
	id: string;
	label?: string | ReactNode;
	name: string;
	placeholder?: string;
	type?: string;
	disabled?: boolean;
	size?: InputSizes;
}

export const FormItem: React.FC<FormItemProps> = props => {
	const { className, id, label, name, type } = props;
	return (
		<Field type={type} name={name}>
			{(fieldProps: FieldProps): ReactNode => {
				const { field, meta } = fieldProps;
				const { touched, error } = meta;
				const hasError = touched && error;
				return (
					<div className={classNames('form-group', className)}>
						{label && <label htmlFor={id}>{label}</label>}
						<FormField {...props} field={field} meta={meta} />
						<div className="invalid-feedback">{hasError ? error : ''}</div>
					</div>
				);
			}}
		</Field>
	);
};
