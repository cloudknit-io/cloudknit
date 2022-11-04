import { ZDropdownMenuJSX } from 'components/molecules/dropdown-menu/DropdownMenu';
import { Context } from 'context/argo/ArgoUi';
import React, { useState } from 'react';
import { ReactComponent as Avatar } from 'assets/images/icons/notification.svg';
import './styles.scss';
import { getTime, momentHumanizer } from 'pages/authorized/environment-components/helpers';

export const ZNotifications = () => {
	const notifications: any[] = [];
	const unseenCount = 0;
	const [open, setOpen] = useState<boolean>(false);

	return (
		<div className="top-bar__notifications">
			<Avatar onClick={() => {
				setOpen(!open);
			}} />
			{unseenCount > 0 && <span className={`top-bar__notifications--unseen-bubble`}>{unseenCount}</span>}
			<ZDropdownMenuJSX
				className="top-bar__notifications__dropdown"
				isOpened={open}
				items={[
					...(notifications.sort((a: any, b: any) => b.notification_id - a.notification_id)).map(n => ({
						text: '',
						action: () => {
							
						},
						jsx: (
							<span className={`top-bar__notifications__dropdown--item message-type-${n.message_type} ${n.seen ? 'seen' : ''}`}>
								<span>{n.message}</span>
								<span className="timestamp">{new Date(n.timestamp).toLocaleString()}</span>
							</span>
						),
					})),
				]}
			/>
		</div>
	);
};
