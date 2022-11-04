export interface BreadcrumbInfo {
	title: string;
	path: string;
	filters?: BreadcrumbFilter[];
}

export interface BreadcrumbFilter {
	title: string;
	routePath: string;
}
