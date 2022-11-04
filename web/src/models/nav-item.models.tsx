export interface NavItem {
	title: string;
	path: string;
	children?: NavItem[];
	optionalData?: any;
	visible?: () => boolean;
}
