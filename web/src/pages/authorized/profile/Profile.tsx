import React, { useEffect, useMemo } from 'react';
import './styles.scss';

import { usePageHeader } from '../contexts/EnvironmentHeaderContext';
import { useState } from 'react';
import { AWSSecret } from '../secrets/aws-secrets';
import { ZLoaderCover } from 'components/atoms/loader/LoaderCover';
import { ENVIRONMENT_VARIABLES } from 'utils/environmentVariables';
import { SecretList } from '../secrets/secrets-list';
import { ZTablControl } from 'components/molecules/tab-control/TabControl';
import { HierachicalLeftView, Hierarchy } from '../secrets/hierarchical-left-view';
import { AWSSSMSecret } from '../secrets/aws-ssm-secrets';
import { Subject } from 'rxjs';
import { Secrets } from '../secrets/secrets';
import { AddUser } from '../addUser/AddUser';

const refresher = new Subject<string>();

const secretTabs = [
	{
		id: 'AWS',
		name: 'AWS Credentials',
		show: (secret: string) => true,
	},
	{
		id: 'SSM',
		name: 'Secrets',
		show: (secret: string) => true,
	},
	{
		id: 'tfState',
		name: 'TF State Credentials',
		show: (secret: string) => true,
	},
];

export const Profile: React.FC = () => {
	const { pageHeaderObservable, breadcrumbObservable } = usePageHeader();
	const [selectedSecret, setSelectedSecret] = useState<any>();
	const [selectedHierarchy, setSelectedHierarchy] = useState<Hierarchy>();
	const [loading, setLoading] = useState<boolean>(false);

	useEffect(() => {
		pageHeaderObservable.next({
			breadcrumbs: [],
			headerTabs: [],
			pageName: 'Settings',
			filterTitle: '',
			onSearch: () => {},
			buttonText: '',
			onViewChange: () => {},
		});
	});

	useEffect(() => {
		breadcrumbObservable.next(false);
	}, [breadcrumbObservable]);

	return (
		<ZLoaderCover loading={loading}>
			<div className="secrets-container">
				<div className="secrets-container__list">
					<HierachicalLeftView
						hierarchyChanged={(hierarchy: Hierarchy) => {
							setSelectedHierarchy(hierarchy);
							setSelectedSecret(hierarchy.id.replaceAll(':', '/'));
							return true;
						}}
						refreshView={refresher}
					/>
				</div>
				{selectedHierarchy?.type === 'SECRET' ? (
					<div className="secrets-container">
						<div className="secrets-container__tabs">
							<span className="d-flex">{selectedSecret.replaceAll('/', ' > ')}</span>
							<ZTablControl
								className={`secrets-container__tabs-control`}
								tabs={secretTabs.filter(st => st.show(selectedSecret))}
								selected={'AWS'}>
								<div id="AWS">
									<div className="aws-secrets">
										<div
											className={`secret-info secret-info-ssm`}
											onClick={e => {
												if (e.currentTarget === e.target) {
													// awsSecretCallback();
												}
											}}>
											<div className="secret-container">
												{
													<AWSSecret
														key={selectedSecret.id}
														newEnvironmentField={selectedSecret?.endsWith('/')}
														secretScope={selectedSecret}
														closeCallback={(secret?: string) => {
															selectedSecret?.endsWith('/') && refresher.next(secret);
														}}
													/>
												}
											</div>
										</div>
									</div>
								</div>
								<div id="SSM">
									<div className="aws-secrets">
										{selectedSecret.endsWith('/') ? (
											<div className="secret-info secret-info-ssm">
												<div className="secret-container">
													<AWSSSMSecret
														secretScope={selectedSecret}
														secretKey={null}
														saveCallback={(id: string) => {
															refresher.next(id);
														}}
														closeCallback={() => {}}
														scopeEditable={true}
													/>
												</div>
											</div>
										) : (
											<SecretList heading={''} secretKey={selectedSecret?.replaceAll('/', ':')} />
										)}
									</div>
								</div>
								<div id="tfState">
									<div className="aws-secrets">
										<div
											className={`secret-info secret-info-ssm`}
											onClick={e => {
												if (e.currentTarget === e.target) {
													// awsSecretCallback();
												}
											}}>
											<div className="secret-container">
												{
													<Secrets
														newEnvironmentField={selectedSecret?.endsWith('/')}
														secretScope={selectedSecret}
														closeCallback={(secret?: string) => {
															selectedSecret?.endsWith('/') && refresher.next(secret);
														}}
														secretModels={[
															{
																name: 'Bucket',
																key: 'state_bucket',
																immutable: true,
															},
															{
																name: 'Lock Table',
																key: 'state_lock_table',
																immutable: true,
															},
															{
																name: 'AWS Access Key Id',
																key: 'state_aws_access_key_id',
															},
															{
																name: 'AWS Secret Access Key',
																key: 'state_aws_secret_access_key',
																multiline: true,
															},
														]}
													/>
												}
											</div>
										</div>
									</div>
								</div>
							</ZTablControl>
						</div>
					</div>
				) : (
					<AddUser />
				)}
			</div>
		</ZLoaderCover>
	);
};
