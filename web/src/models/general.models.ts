import { ReactNode } from 'react';

export interface OptionItem {
	id: string;
	name: string;
}

export interface ListItem {
	label: string;
	value: string | ReactNode;
}
