import { WorkflowPayload } from 'services/argo/ArgoWorkflows.service';
import { EventClientCD, EventClientCDResourceTree, EventClientCDWatcher, EventClientCost, EventClientParallelWF, EventClientWF } from 'utils/apiClient/EventClient';

export class ArgoStreamService {
	static stream(resourceVersion: string): void {
		new EventClientCD(`/cd/api/v1/stream/projects/${resourceVersion}`).listen();
	}

	static streamWF({ projectId, environmentId, configId, workflowId }: WorkflowPayload): void {
		new EventClientWF(
			`/wf/api/v1/stream/projects/${projectId}/environments/${environmentId}/config/${configId}/${workflowId}`
		).listen();
	}

	static streamEnvironment(environmentId: string): void {
		new EventClientCD(`/cd/api/v1/stream/environments/${environmentId}`).listen();
	}

	static streamWatcher(teamName: string): void {
		new EventClientCDWatcher(`/cd/api/v1/stream/watcher/projects/${teamName}`).listen();
	}

	static streamResourceTree(name: string): void {
		new EventClientCDResourceTree(`/cd/api/v1/stream/applications/${name}/resource-tree`).listen();
	}
}
