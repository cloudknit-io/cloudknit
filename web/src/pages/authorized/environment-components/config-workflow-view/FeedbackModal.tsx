import { Button } from 'components/atoms/button/Button';
import { Loader } from 'components/atoms/loader/Loader';
import { ZText } from 'components/atoms/text/Text';
import React, { FC, useState } from 'react';
import { useEffect } from 'react';

type Props = {
	approved: boolean;
	comment?: string;
	onDecline: Function;
	onApprove: Function;
};

export const ZFeedbackModal: FC<Props> = ({ onApprove, onDecline }: Props) => {
	const [mode, setMode] = useState<string>('');
	const [showLoader, setShowLoader] = useState<boolean>(false);

	return (
		<div className="zlifecycle-feedback-modal">
			{mode === '' && (
				<>
					<div className="zlifecycle-feedback-modal__content">
						<ZText.Body size="14" lineHeight="18">
							In order to continue, please approve:
						</ZText.Body>
					</div>
					<div className="zlifecycle-feedback-modal__actions">
						<Button color="primary" onClick={() => setMode('approve')}>
							Approve
						</Button>
						{/* <Button color="secondary" onClick={() => setMode('decline')}>
							Decline
						</Button> */}
					</div>
				</>
			)}
			{mode === 'approve' && (
				<>
					<div className="zlifecycle-feedback-modal__content">
						<ZText.Body size="16" lineHeight="18">
							In order to continue, please accept the changes.
						</ZText.Body>
					</div>
					<div className="zlifecycle-feedback-modal__actions">
						{!showLoader && (
							<>
								<Button
									color="primary"
									onClick={() => {
										setShowLoader(true);
										onApprove();
									}}>
									Confirm
								</Button>
								{/* <Button color="secondary" onClick={() => setMode('')}>
									Cancel
								</Button> */}
							</>
						)}
						{showLoader && <Loader />}
					</div>
				</>
			)}
			{mode === 'decline' && (
				<>
					<div className="zlifecycle-feedback-modal__content">
						<ZText.Body weight="bold" size="14" lineHeight="18">
							Changes Declined
						</ZText.Body>
					</div>
					<div className="zlifecycle-feedback-modal__actions">
						{!showLoader && (
							<>
								<Button
									color="primary"
									onClick={() => {
										setShowLoader(true);
										onDecline();
									}}>
									Confirm
								</Button>
								<Button color="secondary" onClick={() => setMode('')}>
									Cancel
								</Button>
							</>
						)}
						{showLoader && <Loader />}
					</div>
				</>
			)}
		</div>
	);
};
