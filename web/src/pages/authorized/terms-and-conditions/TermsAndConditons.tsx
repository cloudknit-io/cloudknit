import AuthStore from 'auth/AuthStore';
import { Notifications, NotificationsManager, NotificationType } from 'components/argo-core';
import React, { useEffect, useState } from 'react';
import { ErrorResponse } from 'utils/apiClient/ApiClient';
import './styles.scss';

export const TermsAndConditions: React.FC = () => {
	const formRef = React.useRef<HTMLFormElement>(null);
	const nm = new NotificationsManager();
	const user = AuthStore.getUser();
	const [retryCount, setRetryCount] = useState(60);

	useEffect(() => {
		if (user?.selectedOrg?.provisioned !== true) {
			const interval = setInterval(() => {
				setRetryCount(prev => {
					if (prev === 1) {
						checkOrganizationStatus();
					}
					return prev - 1;
				});
			}, 1000);

			return () => {
				clearInterval(interval);
			};
		}
	}, [user]);

	if (!user) {
		return <></>;
	}

	const checkOrganizationStatus = async () => {
		try {
			const resp = await AuthStore.fetchOrganization(AuthStore.getOrganization()?.id);
			if (resp?.provisioned) {
				await AuthStore.selectOrganization(resp.name);
			} else {
				setRetryCount(60);
				nm.show({
					content: resp?.name + ' is still being provisioned.',
					type: NotificationType.Warning,
				});
			}
		} catch (err) {
			setRetryCount(60);
			nm.show({
				content: 'There was some error fetching orgs data',
				type: NotificationType.Error,
			});
		}
	};

	const submitOrganizationForm = async () => {
		try {
			if (!formRef.current?.checkValidity()) {
				formRef.current?.classList.add('invalid');
				return;
			}
			const fd = new FormData(formRef.current);
			const data = {
				agreedByUsername: fd.get('name'),
				agreedByEmail: fd.get('email'),
				organizationName: fd.get('organizationName')?.toString(),
			};
			if (data.organizationName) {
				nm.show({
					content: 'Creating Organization...',
					type: NotificationType.Warning,
				});
				const org = await AuthStore.addOrganization(data.organizationName);
				nm.show({
					content: `Your organization ${org?.name} was successfully created.`,
					type: NotificationType.Success,
				});
				return;
			}
		} catch (err) {
			const { message } = err as ErrorResponse;
			let content = 'There was some error saving data.';
			if (message) {
				content = message;
			}
			nm.show({
				content,
				type: NotificationType.Error,
			});
		}
	};

	return (
		<>
			{!AuthStore.getOrganization() ? (
				<div className="conditions-page">
					<section className="conditions-page__condition">
						<em>
							Please review the{' '}
							<a href="https://compuzest.com/terms-of-service/" target="_blank">
								Terms of Service
							</a>{' '}
							and{' '}
							<a href="https://compuzest.com/privacy-policy/" target="_blank">
								Privacy Policy.
							</a>
						</em>
					</section>
					<form ref={formRef} noValidate onSubmit={e => e.preventDefault()} className="conditions-page__form">
						<section className="conditions-page__form__condition-form">
							<label>Name</label>
							<input
								name="name"
								required
								className="shadowy-input input"
								type="text"
								defaultValue={user.name}
							/>
							<label className="error-msg">Please enter a valid value</label>
						</section>
						<section className="conditions-page__form__condition-form">
							<label>Organization Name</label>
							<input
								name="organizationName"
								required
								pattern="^[a-zA-Z]+[a-zA-Z0-9]*(-[a-zA-Z0-9]+)*$"
								className="shadowy-input input"
								minLength={1}
								maxLength={63}
								type="text"
							/>
							<label className="error-msg">Organization name can only contain alphanumeric characters, '-', and should start and end with alphanumeric with at most 63 characters.</label>
						</section>
						<section className="conditions-page__form__condition-form">
							<label>Email</label>
							<input
								name="email"
								required
								className="shadowy-input input"
								type="email"
								pattern="[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$"
								defaultValue={user.email}
							/>
							<label className="error-msg">Please enter a valid value</label>
						</section>
						<section className="conditions-page__form__condition">
							<label>
								<input name="terms" required type="checkbox" /> I agree to the Terms of Service.
								<label className="error-msg">Please agree to Terms of Service</label>
							</label>
						</section>
						<section className="conditions-page__form__condition">
							<label>
								<input name="policy" required type="checkbox" /> I acknowledge the Privacy Policy.
								<label className="error-msg">Please agree to Privacy Policy</label>
							</label>
						</section>
						<button onClick={async () => await submitOrganizationForm()}>Submit</button>
					</form>
					<Notifications notifications={nm.notifications} />
				</div>
			) : (
				<div className="conditions-page">
					<em>
						Please wait while we provision "{user.selectedOrg.name}" organization (Auto refreshing in{' '}
						{retryCount} {retryCount === 1 ? 'second' : 'seconds'})
					</em>
					<div className="btn-container">
						<button className="base-btn m-t-10" onClick={() => checkOrganizationStatus()}>
							Refresh Now
						</button>
					</div>
				</div>
			)}
		</>
	);
};
