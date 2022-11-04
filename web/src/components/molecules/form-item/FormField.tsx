import classNames from 'classnames';
import { FormItemProps } from 'components/molecules/form-item/FormItem';
import { FieldInputProps, FieldMetaProps } from 'formik/dist/types';
import React from 'react';

interface FormFieldProps extends FormItemProps {
	field: FieldInputProps<any>;
	meta: FieldMetaProps<any>;
}

export const FormField: React.FC<FormFieldProps> = props => {
	const { id, disabled, size, placeholder, type, field, meta } = props;
	const { touched, error } = meta;

	if (type === 'textarea') {
		return (
			<textarea
				id={id}
				className={classNames('form-control', {
					[`form-control-${size}`]: Boolean(size),
					'is-invalid': touched && error,
				})}
				placeholder={placeholder}
				disabled={disabled}
				{...field}
			/>
		);
	}

	return (
		<input
			id={id}
			className={classNames('argo-field', {
				[`form-control-${size}`]: Boolean(size),
				'is-invalid': touched && error,
			})}
			placeholder={placeholder}
			type={type}
			{...field}
		/>
	);
};
