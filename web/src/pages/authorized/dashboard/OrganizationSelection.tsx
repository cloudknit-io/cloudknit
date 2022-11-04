import AuthStore from 'auth/AuthStore';
import { getDropDownList } from 'helpers/environment-builder.helper';
import { Organization} from 'models/user.models';
import React, { useState } from 'react';
import { FC } from 'react';

export const OrganizationSelection: FC = () => {
	const user = AuthStore.getUser();
	const [selectedOrg, setSelectedOrg] = useState<Organization>();

	if (!user) {
		AuthStore.login();
	}

	const orgs = user?.organizations || [];

	return (
		<section className="d-flex align-center justify-center flex-dir-column" style={{ height: '100%' }}>
			<em style={{ textAlign: 'center' }}>Please choose an organization you want to log in</em>
			<section>
				{getDropDownList(
					new Set<string>(orgs.map(e => e.name)),
					(val: string) => {
						setSelectedOrg(orgs.find(e => e.name === val));
					},
					(val: string) => (selectedOrg?.name === val ? 'selected' : ''),
					true,
					selectedOrg?.name || '--None--',
					'',
					['h-30px', 'm-t-10']
				)}
			</section>
			<section>
				<button
					className="base-btn m-t-10"
					disabled={!Boolean(selectedOrg)}
					onClick={async () => {
						await AuthStore.selectOrganization(selectedOrg?.name);
					}}>
					Proceed
				</button>
			</section>
		</section>
	);
};
