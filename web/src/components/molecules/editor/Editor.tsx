import './styles.scss';
import Editor, { EditorProps } from '@monaco-editor/react';
import React, { FC } from 'react';

interface Props {
	data: any;
	language?: string;
	height?: string;
	readOnly?: boolean;
	options?: any;
	onChange?: (value: EditorProps['value']) => void;
}

export const ZEditor: FC<Props> = (props: Props) => {
	const { data, language, height, readOnly, options, onChange } = props;

	return (
		<>
			<Editor
				wrapperClassName="zlifecycle-editor"
				height={height || '64vh'}
				defaultLanguage={language || 'yaml'}
				defaultValue={''}
				theme="light"
				value={data}
				options={{
					readOnly: readOnly,
					scrollbar: {
						alwaysConsumeMouseWheel: false
					},
					...options,
				}}
				onChange={onChange}
			/>
		</>
	);
};
