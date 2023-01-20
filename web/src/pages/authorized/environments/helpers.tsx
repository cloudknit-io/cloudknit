import { ReactComponent as AWSIcon } from 'assets/images/icons/AWS.svg';
import { ReactComponent as ChevronRight } from 'assets/images/icons/chevron-right.svg';
import { ReactComponent as MoreOptionsIcon } from 'assets/images/icons/more-options.svg';
import { ReactComponent as SyncIcon } from 'assets/images/icons/sync-icon.svg';
import { Checkbox } from 'components/argo-core/checkbox';
import { TableColumn } from 'components/atoms/table/Table';
import { renderHealthStatus, renderLabels, renderSyncedStatus } from 'components/molecules/cards/renderFunctions';
import { MenuItem, ZDropdownMenuJSX } from 'components/molecules/dropdown-menu/DropdownMenu';
import { ApplicationCondition, HealthStatuses, ZEnvSyncStatus } from 'models/argo.models';
import { Environment } from 'models/entity.type';
import React, { ReactElement } from 'react';
import { ArgoTeamsService } from 'services/argo/ArgoProjects.service';

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
		return 'There seems to be a problem with your environment. Please contact your administrator.';
	//errors.map(e => `${e.type}: ${e.message}`).join('\n');
	return null;
};

export const hardSync = async (projectId: string, envName: string) => {
	try {
		await ArgoTeamsService.hardSyncTeam(projectId, envName);
	} catch (err) {}
};

export const TreeReconcile = (
	env: Environment,
	reconciling: boolean,
	triggerSync: () => Promise<any>
): ReactElement => {
	return (
		<button
			className="dag-controls-reconcile"
			onClick={async (e: any) => {
				e.stopPropagation();
				if (!reconciling) await triggerSync();
			}}>
			<span
				className={`tooltip ${
					// environmentItem?.healthStatus !== 'Progressing' &&
					// !syncStarted &&
					// environmentCondition &&
					// 'error'
					''
				}`}>{`${reconciling ? 'Reconciling...' : 'Reconcile'}`}</span>
			<SyncIcon
				className={`large-health-icon-container__sync-button large-health-icon-container__sync-button${getSyncIconClass(
					env
				)} large-health-icon-container__sync-button${reconciling ? '--in-progress' : ''}`}
				title="Reconcile Environment"
			/>
			Reconcile
		</button>
	);
};

export const EnvCardReconcile = (env: Environment, reconciling: boolean, triggerSync: () => Promise<any>) => {
	return (
		<SyncIcon
			className={`large-health-icon-container__sync-button large-health-icon-container__sync-button${getSyncIconClass(
				env
			)} large-health-icon-container__sync-button${reconciling ? '--in-progress' : ''}`}
			title={'Reconcile Environment'}
			onClick={async e => {
				e.stopPropagation();
				//TODO: Syncing env
				if (!reconciling) await triggerSync();
			}}
		/>
	);
};

const getSyncIconClass = (environment: Environment) => {
	if ([ZEnvSyncStatus.DestroyFailed, ZEnvSyncStatus.ProvisionFailed].includes(environment.status as ZEnvSyncStatus)) {
		return '--out-of-sync';
	} else if ([ZEnvSyncStatus.Provisioned, ZEnvSyncStatus.Destroyed].includes(environment.status as ZEnvSyncStatus)) {
		return '--in-sync';
	} else {
		return '--in-sync';
	}
};

const renderServices = () => <AWSIcon />;

const renderActions = () => <MoreOptionsIcon />;

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
		render: data => -1,
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
	getCount?: (status: string) => number
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
		<ZDropdownMenuJSX className="checkbox-filter__inline-flex" label={title} isOpened={true} items={menuItems} />
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
