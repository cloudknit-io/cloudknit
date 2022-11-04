import { NotificationType } from 'components/argo-core';
import { Loader } from 'components/atoms/loader/Loader';
import { Context } from 'context/argo/ArgoUi';
import React, { FormEvent, useEffect, useState } from 'react';
import ApiClient from 'utils/apiClient';
import { ENVIRONMENT_VARIABLES } from 'utils/environmentVariables';
import { ReactComponent as DeleteIcon } from 'assets/images/icons/card-status/sync/delete.svg';
import './styles.scss';
import AuthStore from 'auth/AuthStore';

export type User = {
	githubId: string;
	email: string;
	timestamp: string;
	role: 'Admin' | 'User' | string;
};

export const AddUser: React.FC = () => {
	const notificationsApi = React.useContext(Context)?.notifications;
	const [addedUsers, setAddedUsers] = useState<User[] | null>(null);
	const currentUser = AuthStore.getUser();
	const columns = [
		{
			id: 'email',
			value: 'Email Id',
		},
		{
			id: 'githubId',
			value: 'Github Id',
		},
		{
			id: 'role',
			value: 'Role',
		},
		{
			id: 'timestamp',
			value: 'Created On',
		},
	];

	useEffect(() => {
		ApiClient.get(`/users/get`)
			.then(({ data }) => {
				const users = data as any;
				setAddedUsers(
					users.map((e: any) => ({
						githubId: e.username,
						email: e.email,
						timestamp: e.created,
						role: e.role || 'User',
					}))
				);
			})
			.catch(e => {
				setAddedUsers([]);
			});
	}, []);

	const addUser = (e: FormEvent) => {
		e.preventDefault();
		const form = e.target as HTMLFormElement;
		const formData = new FormData(form);
		const email = formData.get('email')?.toString().trim();
		const username = formData.get('githubId')?.toString().trim();
		const organizationId = ENVIRONMENT_VARIABLES.REACT_APP_CUSTOMER_NAME;
		const role = formData.get('role')?.toString().trim();
		form.style.pointerEvents = 'none';

		ApiClient.post(`/users/add`, {
			email,
			username,
			role,
		})
			.then(() => {
				addedUsers?.push({
					email: email || '',
					githubId: username || '',
					role: role || 'User',
					timestamp: new Date().toISOString(),
				});
				setAddedUsers([...(addedUsers as any)]);
				form.reset();
				form.style.pointerEvents = 'auto';
			})
			.catch((err: any) => {
				const { message } = err;
				notificationsApi?.show({
					content: message || 'There was an error adding the user',
					type: NotificationType.Error,
				});
				form.style.pointerEvents = 'auto';
			});
	};

	return (
		<div className="add-user">
			<h4>Create User</h4>
			<form className="add-user__form" onSubmit={e => addUser(e)}>
				<div className="add-user__form-pair">
					<label>Email Id</label>
					<input name="email" type="email" required />
				</div>
				<div className="add-user__form-pair">
					<label>Github Id</label>
					<input name="githubId" type="text" required />
				</div>
				<div className="add-user__form-pair">
					<label>Role</label>
					<select name="role" required defaultValue={'User'}>
						<option value={'Admin'}>Admin</option>
						<option value={'User'}>User</option>
					</select>
				</div>
				<button type="submit">Create</button>
			</form>
			<div className="add-user__list">
				<h4>Existing Users</h4>
				<div className="add-user__list--columns">
					{columns.map(e => (
						<span key={e.id}>{e.value}</span>
					))}
					<span>Delete</span>
				</div>
				{addedUsers ? (
					addedUsers.map((r: User) => (
						<div
							key={r.email}
							className={`add-user__list--row ${currentUser?.username === r.githubId ? 'disabled' : ''}`}>
							{columns.map(c => (
								<span>
									{c.id === 'timestamp' ? new Date(r[c.id]).toDateString() : (r as any)[c.id]}
								</span>
							))}

							<span>
								{currentUser?.username !== r.githubId && (
									<DeleteIcon
										style={{ cursor: 'pointer' }}
										height={16}
										width={16}
										onClick={e => {
											if (window.confirm('Are you sure, you want to remove this user?')) {
												const ele = e.target as HTMLSpanElement;
												ele.style.pointerEvents = 'none';
												ApiClient.delete(`/users/delete/${r.githubId}`)
													.then(() => {
														notificationsApi?.show({
															content: 'Deleted succesfully',
															type: NotificationType.Success,
														});
														ele.style.pointerEvents = 'auto';
														addedUsers.splice(addedUsers.indexOf(r), 1);
														setAddedUsers([...addedUsers]);
													})
													.catch(err => {
														notificationsApi?.show({
															content: 'Some error occurred.',
															type: NotificationType.Error,
														});
														ele.style.pointerEvents = 'auto';
													});
											}
										}}
									/>
								)}
							</span>
						</div>
					))
				) : (
					<span className="d-flex justify-center align-center">
						Loading Users <Loader height={16} />
					</span>
				)}
			</div>
		</div>
	);
};
