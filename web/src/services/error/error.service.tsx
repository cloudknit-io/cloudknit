import { EnvironmentStatus, ErrorEvent } from 'models/error.model';
import { EventClient } from 'utils/apiClient/EventClient';
import { Subject } from 'rxjs';
import ApiClient from 'utils/apiClient';

export class ErrorService {
	private static instance: ErrorService | null = null;
	private constructUri = (path: string) => `/events/${path}`;
	private constructApiUri = () => `/error-api`;
	private errorModelStream: Subject<EnvironmentStatus> | null = null;
  private streamMap = new Map<Subject<any>, EventClient<any>>();

	private constructor() {
	}

	static getInstance() {
		if (!ErrorService.instance) {
			ErrorService.instance = new ErrorService();
		}
		return ErrorService.instance;
	}

	subscribeToErrorStream() {
		if (!this.errorModelStream) {
			const url = this.constructUri(ErrorUriType.errorStream);
			const eventClient = new EventClient<EnvironmentStatus>(url);
			this.errorModelStream = eventClient.listen();
			this.streamMap.set(this.errorModelStream, eventClient);
		}
  	
		return this.errorModelStream;
	}

	disposeStreams(...streams: Subject<any>[]) {
		streams.forEach(stream => {
			const client = this.streamMap.get(stream);
			
			if (client) {
				client.close();
			}
		});
	}

	async getEventData() {
		const url = this.constructApiUri();
		try {
			const res = await ApiClient.get<any>(url);

			if (!res.data || !res.data.status) {
				return null;
			}

			// status.environmentStatus.environment.team
			const status = res.data.status
			const environmentStatus = status.environmentStatus;

			for (const teamKey of Object.keys(environmentStatus)) {
				const team = environmentStatus[teamKey];
				for (const envKey of Object.keys(environmentStatus[teamKey])) {
					const env = team[envKey];
					const envEvents : ErrorEvent[] = [];

					if (env.status.status === "ok") {
						continue;
					}
	
					for (let event of env.status.events) {
						envEvents.push({
							company: event.meta.company,
							createdAt: new Date(event.createdAt),
							environment: event.meta.environment,
							eventType: event.eventType,
							id: event.id,
							team: event.meta.team
						});
					}
					
					this.errorModelStream?.next({
						company: env.company,
						environment: env.environment,
						errors: env.status.status.validation?.errors,
						events: envEvents,
						status: env.status.status,
						team: env.team,
					});
				}
			}
			
			return true; // super useful
		} catch (err) {
			console.log("getEventData error:", err);
			return null;
		}
	}
}

class ErrorUriType {
	static errorStream = `stream`;
}
