export enum Roles {
	CLIENT = 1,
	EXPERT = 2,
}

export interface Billing {
	phone: string;
	address: string;
	zip: string;
	city: string;
	country: number;
}

export interface Expert {
	title: string;
	description: string;
}

export interface Organization {
	id: number;
	name: string;
	updated: Date;
	created: Date;
	provisioned: boolean;
	githubRepo: string;
}

export interface Profile {
	id: number;
	username: string;
	email: string;
	termAgreementStatus: boolean;
	archived: boolean;
	role: string;
	created: Date;
	updated: Date;
	organizations: Organization[];
	picture: string;
	terms: boolean;
	name: string;
	selectedOrg: Organization;
}

export type User = Profile;
