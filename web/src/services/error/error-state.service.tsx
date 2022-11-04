import { EnvironmentStatus, ErrorEvent, ErrorStatus, EventMessage } from 'models/error.model';
import { ErrorService } from './error.service';
import { Subject } from 'rxjs';

export class ErrorStateService {
	private static instance: ErrorStateService | null = null;
	private errorStateEnvironment: Map<string, EnvironmentStatus> = new Map<string, EnvironmentStatus>();
	public updates: Subject<void> = new Subject<void>();

	private constructor() {
		const errorInstance = ErrorService.getInstance();
		errorInstance.subscribeToErrorStream().subscribe((data: EnvironmentStatus) => this.errorStateData(data));
		errorInstance.getEventData();
	}

	static getInstance() {
		if (!ErrorStateService.instance) {
			ErrorStateService.instance = new ErrorStateService();
		}
		return ErrorStateService.instance;
	}

	private getTimestamp(e: EnvironmentStatus, err: string) {
		const errorEvent = e.events.find(ev => (ev.payload || []).includes(err));
		return errorEvent?.createdAt ? new Date(errorEvent?.createdAt).toUTCString() : null;
	}

	private errorMapper(e: EnvironmentStatus, err: string): EventMessage {
		// @ts-ignore
		return {
			company: e.company,
			environment: e.environment,
			team: e.team,
			message: err,
			timestamp: e.events.length > 0 ? e.events[0].createdAt.toLocaleString() : 'N/A'
		};
	}

	errorStateData(e: EnvironmentStatus) {
		this.errorStateEnvironment.set(e.environment, e);
		this.updates.next();
	}

	errorsInEnvironment(environmentId: string) {
		const field = this.errorStateEnvironment.get(environmentId);
		if (!field) {
			return [];
		}
		return (field.errors || []).map(e => this.errorMapper(field, e));
	}

	public get Errors(): EventMessage[] {
		const msgs: EventMessage[] = [];

		const errors = [...this.errorStateEnvironment.values()];
		const filtered = errors.filter(e => e.errors?.length > 0);
		
		filtered.forEach(error => {
			msgs.push(...error.errors.map(er => this.errorMapper(error, er)));
		});
			
		return msgs;
	}

	public get ErrorsEnvs(): EnvironmentStatus[] {
		return [...this.errorStateEnvironment.values()].filter(e => e.errors?.length > 0);
	}
}
