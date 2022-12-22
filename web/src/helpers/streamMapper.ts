import { ApplicationWatchEvent } from 'models/argo.models';

const validator = (
	checkType: string,
	data: ApplicationWatchEvent,
	queryParams?: { projectId?: string; environmentId?: string; workflowId?: string }
): boolean => {
	const pId = data.application?.metadata?.labels?.project_id;
	const eId = data.application?.metadata?.labels?.environment_id;

	if (checkType === 'environment') {
		return pId === queryParams?.projectId;
	}

	if (checkType === 'config') {
		return pId === queryParams?.projectId && eId === queryParams?.environmentId;
	}

	return true;
};

export const streamMapperWF = (data: ApplicationWatchEvent | null) => {
	return data;
};

export const streamMapper = <T>(
	data: ApplicationWatchEvent | null,
	items: any[],
	mapper: Function,
	checkType: string,
	queryParams?: { projectId?: string; environmentId?: string; workflowId?: string }
): any[] => {
	if (data) {
		const type = data.application?.metadata?.labels?.type;

		if (type === checkType && validator(checkType, data, queryParams)) {
			const item = mapper(data.application);

			let newItems: T[] = items;

			const index = items.findIndex(x => x.id === item.id);
			if (index === -1 && data.type === 'ADDED') {
				newItems = [...items, item];
			}

			if (data.type === 'DELETED') {
				items.splice(index, 1);
				newItems = items;
			}

			if (index !== -1 && data.type === 'MODIFIED') {
				newItems = items.map(x => {
					if (x.id === item.id) {
						//NOTE: This is done so that argo doesn't override cost and status since we are fetching them from API
						item.componentCost = x.componentCost;
						item.componentStatus = x.componentStatus;
						item.costResources = x.costResources;
						item.syncFinishedAt = x.syncFinishedAt
						return item;
					}
					return x;
				});
			}

			return newItems;
		}
	}
	return items;
};
