import * as React from 'react';
import { toast, ToastContainer, ToastOptions } from 'react-toastify';
import { Observable } from 'rxjs';
import { useEffect } from 'react';

require('react-toastify/dist/ReactToastify.css');

const AUTO_CLOSE_TIMEOUT = 10000;

export enum NotificationType {
	Success,
	Warning,
	Error,
}

export interface NotificationInfo {
	type: NotificationType;
	content: React.ReactNode;
	toastOptions?: ToastOptions;
}

export interface NotificationsProps {
	notifications: Observable<NotificationInfo>;
}

export const Notifications = (props: NotificationsProps) => {
	useEffect(() => {
		props.notifications.subscribe(next => {
			let toastMethod = toast.success;
			switch (next.type) {
				case NotificationType.Error:
					toastMethod = toast.error;
					break;
				case NotificationType.Warning:
					toastMethod = toast.warn;
					break;
			}
			toastMethod(
				<div
					onClick={e => {
						const sel = window.getSelection();

						if (sel) {
							const range = document.createRange();

							range.selectNode(e.target as Node);
							sel.removeAllRanges();
							sel.addRange(range);
						}
					}}>
					{next.content}
				</div>,
				{
					position: toast.POSITION.BOTTOM_RIGHT,
					closeOnClick: false,
					pauseOnHover: true,
					pauseOnFocusLoss: true,
					draggable: true,
					closeButton: false,
					autoClose: AUTO_CLOSE_TIMEOUT,
					...next.toastOptions
					
				}
			);
		});
	}, [props.notifications]);

	return <ToastContainer className="zlifecycle-toast" />;
};
