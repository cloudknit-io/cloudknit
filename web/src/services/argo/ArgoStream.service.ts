import { ApplicationWatchEvent } from 'models/argo.models';
import { WorkflowPayload } from 'services/argo/ArgoWorkflows.service';
import { EventClient, EventClientWF } from 'utils/apiClient/EventClient';

export class ArgoStreamService {

	static streamWF({ projectId, environmentId, configId, workflowId }: WorkflowPayload): void {
		new EventClientWF(
			`/wf/api/v1/stream/projects/${projectId}/environments/${environmentId}/config/${configId}/${workflowId}`
		).listen();
	}

	static streamEnvironment(environmentId: string): EventClient<ApplicationWatchEvent> {
		return new EventClient<ApplicationWatchEvent>(`/cd/api/v1/stream/environments/${environmentId}`);
	}
}
