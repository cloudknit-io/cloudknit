import { ReactComponent as AWSIcon } from 'assets/images/icons/AWS.svg';
import { ReactComponent as ChevronRight } from 'assets/images/icons/chevron-right.svg';
import { ReactComponent as MoreOptionsIcon } from 'assets/images/icons/more-options.svg';
import { NotificationsApi, NotificationType } from 'components/argo-core';
import { Checkbox } from 'components/argo-core/checkbox';
import { TableColumn } from 'components/atoms/table/Table';
import {
	CostRenderer,
	renderHealthStatus,
	renderLabels,
	renderSyncedStatus,
} from 'components/molecules/cards/renderFunctions';
import { MenuItem, ZDropdownMenuJSX } from 'components/molecules/dropdown-menu/DropdownMenu';
import { ZStreamRenderer } from 'components/molecules/zasync-renderer/ZStreamRenderer';
import { ApplicationCondition, HealthStatuses, OperationPhase, OperationPhases } from 'models/argo.models';
import { Environment } from 'models/entity.store';
import { EnvironmentItem } from 'models/projects.models';
import React from 'react';
import { ArgoComponentsService } from 'services/argo/ArgoComponents.service';
import { ArgoEnvironmentsService } from 'services/argo/ArgoEnvironments.service';
import { ArgoTeamsService } from 'services/argo/ArgoProjects.service';
import { CostingService } from 'services/costing/costing.service';

export const mockOriginalYaml = `
apiVersion: stable.compuzest.com/v1
kind: Environment
metadata:
  name: dev-checkout-adprod
  namespace: zlifecycle
spec:
  teamName: checkout
  envName: adprod
  autoApprove: true
  teardown: true
  components:
    - name: static-assets
      type: terraform
      module:
        source: aws
        name: s3-bucket
      variables:
        - name: bucket
          value: "dev-checkout-adprod-static-assets"
      tags:
        - name: componentType
          value: data
        - name: cloudProvider
          value: aws
    - name: networking
      type: terraform
      module:
        source: aws
        name: vpc
      variablesFile:
        source: "git@github.com:zl-dev-tech/checkout-team-config.git"
        path: "adprod/tfvars/networking.tfvars"
      outputs:
        - name: private_subnets
          sensitive: true
      tags:
        - name: componentType
          value: app
    - name: platform-eks
      type: terraform
      dependsOn: [networking]
      module:
        source: aws
        name: s3-bucket
      variablesFile:
        source: "git@github.com:zl-dev-tech/checkout-team-config.git"
        path: "adprod/tfvars/platform-eks.tfvars"
      tags:
        - name: cloudProvider
          value: aws
    - name: eks-addons
      type: terraform
      dependsOn: [platform-eks]
      module:
        source: aws
        name: s3-bucket
      outputs:
        - name: s3_bucket_arn
    - name: postgres
      type: terraform
      dependsOn: [networking]
      module:
        source: aws
        name: s3-bucket
      variables:
        - name: bucket
          value: "dev-checkout-adprod-postgres"
      tags:
        - name: componentType
          value: data`;

export const mockModifiedYaml = `
apiVersion: stable.compuzest.com/v1
kind: Environment
metadata:
  name: dev-checkout-prod
  namespace: zlifecycle
spec:
  teamName: checkout
  envName: prod
  teardown: true
  autoApprove: true
  components:
    - name: static-assets
      type: terraform
      module:
        source: aws
        name: s3-bucket
      variables:
        - name: bucket
          value: "dev-checkout-prod-static-assets"
      tags:
        - name: componentType
          value: data
        - name: cloudProvider
          value: aws
    - name: networking
      type: terraform
      module:
        source: aws
        name: vpc
      variablesFile:
        source: "git@github.com:zl-dev-tech/checkout-team-config.git"
        path: "prod/tfvars/networking.tfvars"
      outputs:
        - name: private_subnets
          sensitive: true
      tags:
        - name: componentType
          value: app
    - name: platform-eks
      type: terraform
      dependsOn: [networking]
      module:
        source: aws
        name: s3-bucket
      variablesFile:
        source: "git@github.com:zl-dev-tech/checkout-team-config.git"
        path: "prod/tfvars/platform-eks.tfvars"
      tags:
        - name: cloudProvider
          value: aws
    - name: eks-addons
      type: terraform
      dependsOn: [platform-eks]
      module:
        source: aws
        name: s3-bucket
      outputs:
        - name: s3_bucket_arn
    - name: platform-ec2
      type: terraform
      dependsOn: [networking]
      module:
        source: aws
        name: ec2-instance
      variables:
        - name: subnet_id
          valueFrom: networking.private_subnets[0]
      variablesFile:
        source: "git@github.com:zl-dev-tech/checkout-team-config.git"
        path: "prod/tfvars/ec2.tfvars"
    - name: postgres
      type: terraform
      dependsOn: [platform-ec2]
      module:
        source: aws
        name: s3-bucket
      variables:
        - name: bucket
          value: "dev-checkout-prod-postgres"
      tags:
        - name: componentType
          value: data`;

export const renderSync = (data: any) => (
	<>
		{renderHealthStatus(data.healthStatus)}
		{renderSyncedStatus(data.syncStatus, data.operationPhase, data.runningStatus)}
	</>
);

export const renderSyncStatus = (data: any) => (
	<>{renderSyncedStatus(data.syncStatus, data.operationPhase, data.runningStatus)}</>
);

export const getEnvironmentErrorCondition = (conditions: ApplicationCondition[]) => {
	const errors = conditions.filter(e => e.type.toLowerCase().includes('error'));
	if (errors.length > 0)
		return 'There seems to be a problem with your environment. Please contact your administrator.'; //errors.map(e => `${e.type}: ${e.message}`).join('\n');
	return null;
};

export const syncMe = async (
	environment: Environment,
	syncStarted: any,
	setSyncStarted: any,
	notificationManager: NotificationsApi,
	watcherStatus: OperationPhase | undefined
) => {
	if (syncStarted) {
		return;
	}
	setSyncStarted(true);
	try {
		notificationManager.show({
			content: `Reconciling ${environment.name}`,
			type: NotificationType.Success,
		});
		// if (watcherStatus === OperationPhases.Failed) {
		// 	await hardSync(environment.labels?.project_id || '', environment.displayValue || '', notificationManager);
		// 	return;
		// }
		await ArgoEnvironmentsService.deleteEnvironment(environment.argoId as string);
		setTimeout(async () => {
			await ArgoEnvironmentsService.syncEnvironment(environment.argoId as string);
		}, 1500);
	} catch (err) {
		// if ((err as any)?.message?.includes('not found') && environment.syncStatus === 'OutOfSync') {
		// 	await ArgoEnvironmentsService.syncEnvironment(environment.id as string);
		// }
		// console.log(err);
	}
};

export const hardSync = async (projectId: string, envName: string, notificationManager: NotificationsApi) => {
	try {
		await ArgoTeamsService.hardSyncTeam(projectId, envName);
	} catch (err) {}
};

const renderServices = () => <AWSIcon />;

const renderActions = () => <MoreOptionsIcon />;

export const renderCost = (teamId?: string, environmentName?: string) => {
	if (!teamId || !environmentName) {
		return <></>;
	}
	return (
		<ZStreamRenderer
			subject={CostingService.getInstance().getEnvironmentCostStream(teamId, environmentName)}
			defaultValue={CostingService.getInstance().getCachedValue(`${teamId}-${environmentName}`)}
			Component={CostRenderer}
		/>
	);
};

export const environmentTableColumns: TableColumn[] = [
	{
		id: 'name',
		name: 'Name',
		width: 250,
	},
	{
		id: 'services',
		name: 'Services',
		width: 100,
		render: renderServices,
	},
	{
		id: 'labels',
		name: 'Labels',
		width: 250,
		render: renderLabels,
	},
	{
		id: 'labels',
		name: 'Cost',
		width: 100,
		render: data => renderCost(data.project_id, data.env_name),
	},
	{
		id: 'healthStatus',
		name: 'Status',
		width: 150,
		combine: true,
		render: renderSyncStatus,
	},
	{
		id: 'repository',
		name: 'Repository',
	},
	{
		id: 'path',
		name: 'Path',
	},
	{
		id: 'actions',
		name: '',
		width: 30,
		render: renderActions,
	},
];

export const renderSyncStatusItems = (
	syncStatuses: any,
	statusFilter: Set<any>,
	setStatusFilter: any,
	title: string,
	getCount?: (status: string) => number,
) => {
	const menuItems: MenuItem[] = [];
	for (const status in syncStatuses) {
		const id = `sync-${status}`;
		const value = syncStatuses[status];
		menuItems.push({
			text: '',
			jsx: (
				<>
					<Checkbox
						id={id}
						value={value}
						onNativeChange={e => setStatusValue(e, statusFilter, setStatusFilter)}
					/>
					<label htmlFor={id}>{status}</label>
					<label>&nbsp;({getCount && getCount(value)})</label>
				</>
			),
			action: () => {},
		});
	}
	return (
		<ZDropdownMenuJSX
			className="checkbox-filter__inline-flex"
			label={title}
			isOpened={true}
			items={menuItems}
		/>
	);
};

export const renderHealthStatusItems = (
	statusFilter: Set<any>,
	setStatusFilter: any,
	getCount?: (status: string) => number
) => {
	const menuItems: MenuItem[] = [];
	for (const status in HealthStatuses) {
		const id = `health-${status}`;
		menuItems.push({
			text: '',
			jsx: (
				<>
					<Checkbox
						id={id}
						value={status}
						onNativeChange={e => setStatusValue(e, statusFilter, setStatusFilter)}
					/>
					<label htmlFor={id}>{status}</label>
					<label>&nbsp; ({getCount && getCount(status)})</label>
				</>
			),
			action: () => {},
		});
	}
	return (
		<ZDropdownMenuJSX
			className="checkbox-filter__inline-flex"
			label="Health Status"
			isOpened={true}
			items={menuItems}
		/>
	);
};

export const getCheckBoxFilters = (
	filterDropDownOpen: boolean,
	setToggleFilterDropDownValue: any,
	filterItems: Array<() => JSX.Element>
) => {
	return (
		<div
			className={`checkbox-filter ${filterDropDownOpen ? 'visible' : ''}`}
			onMouseLeave={() => setToggleFilterDropDownValue(false)}>
			<button
				onClick={() => {
					setToggleFilterDropDownValue(!filterDropDownOpen);
				}}>
				Filters
				<ChevronRight className="dropdown-icon" />
			</button>
			<div
				className={`checkbox-filter__container ${
					filterDropDownOpen ? 'checkbox-filter__container-visible' : ''
				}`}>
				{filterItems.map((e, _i) => (
					<React.Fragment key={_i}>{e()}</React.Fragment>
				))}
			</div>
		</div>
	);
};

export const setStatusValue = (
	e: React.ChangeEvent<HTMLInputElement>,
	statusFilter: Set<any>,
	setStatusFilter: any
): void => {
	const checked = e.currentTarget.checked;
	const value = e.currentTarget.dataset.value;
	if (checked) {
		statusFilter.add(value);
	} else {
		statusFilter.delete(value);
	}
	setStatusFilter(new Set(statusFilter));
};
