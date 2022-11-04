export interface ErrorModel {
	status: Map<string, Map<string, EnvironmentStatus[]>>;
}

export interface EnvironmentStatus {
	events: ErrorEvent[];
	status: ErrorStatus;
	team: string;
	environment: string;
	company: string;
	errors: string[];
}

export interface ErrorEvent {
	id: string;
	company: string;
	team: string;
	environment: string;
	createdAt: Date;
	eventType: EventType;
	payload?: string[];
	debug?: string;
}

export interface EventMessage {
	company: string;
	team: string;
	environment: string;
	message: string;
	timestamp: string;
}

export type EventType = 'validation_success' | 'validation_error';
export type ErrorStatus = 'ok' | 'error';

export const eventErrorColumns = [
	{
		id: 'team',
		name: 'Team',
		width: 100,
	},
	{
		id: 'environment',
		name: 'Environment',
		width: 100,
	},
	{
		id: 'message',
		name: 'Message',
	},
	{
		id: 'timestamp',
		name: 'Timestamp',
	},
];
