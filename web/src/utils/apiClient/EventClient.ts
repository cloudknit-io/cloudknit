import { ApplicationWatchEvent } from 'models/argo.models';
import { processNodeLogs } from 'pages/authorized/environment-components/helpers';
import { BehaviorSubject, Subject } from 'rxjs';

const baseUrl = `${process.env.REACT_APP_STREAM_URL}`;

let eventSourceCD: EventSource;
let eventSourceWF: EventSource;
let eventSourceCost: EventSource;
export const subscriber = new BehaviorSubject<ApplicationWatchEvent | null>(null);
export const subscriberWF = new BehaviorSubject<ApplicationWatchEvent | null>(null);
export const subscriberCost = new BehaviorSubject<any>({});
export const subscriberWatcher = new BehaviorSubject<ApplicationWatchEvent | null>(null);
export const subscriberResourceTree = new BehaviorSubject<ApplicationWatchEvent | null>(null);

export class EventClientCD {
	constructor(url: string) {
		if (eventSourceCD) eventSourceCD.close();

		eventSourceCD = new EventSource(baseUrl + url, { withCredentials: true });
	}

	listen(): void {
		eventSourceCD.onopen = (): void => {
			return;
		};
		eventSourceCD.onmessage = (event: MessageEvent): void => {
			subscriber.next(JSON.parse(event.data)?.result);
		};
		eventSourceCD.onerror = (err): void => {
			console.log(err);
			return;
		};
	}

	close(): void {
		eventSourceCD.close();
	}
}

export class EventClientCDResourceTree {
	constructor(url: string) {
		if (eventSourceCD) eventSourceCD.close();

		eventSourceCD = new EventSource(baseUrl + url, { withCredentials: true });
	}

	listen(): void {
		eventSourceCD.onopen = (): void => {
			return;
		};
		eventSourceCD.onmessage = (event: MessageEvent): void => {
			subscriberResourceTree.next(JSON.parse(event.data)?.result);
		};
		eventSourceCD.onerror = (err): void => {
			console.log(err);
			return;
		};
	}

	close(): void {
		eventSourceCD.close();
	}
}

export class EventClientCDWatcher {
	private eventSourceWatcher: EventSource;
	constructor(url: string) {
		this.eventSourceWatcher = new EventSource(baseUrl + url, { withCredentials: true });
	}

	listen(): void {
		this.eventSourceWatcher.onopen = (): void => {
			return;
		};
		this.eventSourceWatcher.onmessage = (event: MessageEvent): void => {
			subscriberWatcher.next(JSON.parse(event.data)?.result);
		};
		this.eventSourceWatcher.onerror = (err): void => {
			console.log(err);
			return;
		};
	}

	close(): void {
		this.eventSourceWatcher.close();
	}
}

export class EventClientWF {
	constructor(url: string) {
		if (eventSourceWF) eventSourceWF.close();

		eventSourceWF = new EventSource(baseUrl + url, { withCredentials: true });
	}

	listen(): void {
		eventSourceWF.onopen = (): void => {
			return;
		};
		eventSourceWF.onmessage = (event: MessageEvent): void => {
			const resp = JSON.parse(event.data)?.result;
			subscriberWF.next(JSON.parse(event.data)?.result);
			// console.log('--------> workflow data', resp, '-------->', new Date(Date.now()).toDateString());

		};
		eventSourceWF.onerror = (err): void => {
			console.log(err);
			return;
		};
	}

	close(): void {
		eventSourceWF.close();
	}
}

export class EventClientParallelWF {
	subject = new BehaviorSubject<any>({});
	constructor(url: string) {
		if (eventSourceWF) eventSourceWF.close();

		eventSourceWF = new EventSource(baseUrl + url, { withCredentials: true });
	}

	listen(): any {
		eventSourceWF.onopen = (): void => {
			return;
		};
		eventSourceWF.onmessage = (event: MessageEvent): void => {
			this.subject.next(JSON.parse(event.data)?.result);
		};
		eventSourceWF.onerror = (err): void => {
			console.log(err);
			return;
		};
		return this.subject;
	}

	close(): void {
		eventSourceWF.close();
	}
}

export class EventClientCost {
	constructor(url: string) {
		if (eventSourceCost) eventSourceCost.close();

		eventSourceCost = new EventSource(baseUrl + url, { withCredentials: true });
	}

	listen(): void {
		eventSourceCost.onopen = (): void => {
			return;
		};
		eventSourceCost.onmessage = (event: MessageEvent): void => {
			subscriberCost.next(JSON.parse(event.data));
		};
		eventSourceCost.onerror = (err): void => {
			console.log(err);
			return;
		};
	}

	close(): void {
		eventSourceCost.close();
	}
}

export class EventClientAudit {
	private publisher: Subject<any> = new Subject();
	private eventSource: EventSource;
	constructor(url: string) {
		this.eventSource = new EventSource(baseUrl + url, { withCredentials: true });
	}

	listen(): Subject<any> {
		this.eventSource.onopen = (): void => {
			return;
		};
		this.eventSource.onmessage = (event: MessageEvent): void => {
			this.publisher.next(JSON.parse(event.data));
		};
		this.eventSource.onerror = (err): void => {
			console.log(err);
			return;
		};

		return this.publisher;
	}

	close(): void {
		this.eventSource.close();
	}
}

export class EventClientLogs {
	private publisher: Subject<any> = new Subject();
	private eventSource: EventSource;
	private content = '';
	constructor(url: string) {
		this.eventSource = new EventSource(baseUrl + url, { withCredentials: true });
	}

	listen(): Subject<any> {
		this.eventSource.onopen = (): void => {
			this.content = '';
			return;
		};
		this.eventSource.onmessage = (event: MessageEvent): void => {
			this.content += event.data + '\n';
			this.publisher.next(processNodeLogs(this.content));
		};
		this.eventSource.onerror = (err): void => {
			setTimeout(() => {
				this.content = '';
				this.close();
			}, 1000);
			console.log(err);
			return;
		};

		return this.publisher;
	}

	close(): void {
		this.eventSource.close();
	}
}

export class EventClient<T> {
	private publisher: Subject<T> = new Subject<T>();
	private eventSource: any;
	private listenerType = '';
	private handler = (event: MessageEvent): any => {
		this.publisher.next(JSON.parse(event.data) as T);
	};
	
	constructor(url: string, listenerType: string = '') {
		this.eventSource = new EventSource(baseUrl + url, { withCredentials: true });
		this.listenerType = listenerType;
	}
	

	listen(): Subject<T> {
		this.eventSource.onopen = (): void => {
			return;
		};
		
		if (this.listenerType) {
			this.eventSource.addEventListener(this.listenerType, this.handler); 
		} else {
			this.eventSource.onmessage = this.handler;
		}
		
		this.eventSource.onerror = (ev: any): void => {
			console.log('EventSource error', this.eventSource.url, ev);
			return;
		};

		return this.publisher;
	}

	close(): void {
		this.eventSource.close();
		this.listenerType && this.eventSource.removeEventListener(this.listenerType, this.handler);
		this.publisher.unsubscribe();
	}
}
