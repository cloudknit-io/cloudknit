import { Button } from 'components/atoms/button/Button';
import { Loader } from 'components/atoms/loader/Loader';
import { ZText } from 'components/atoms/text/Text';
import React, { FC, useState } from 'react';

type Props = {
	approved: boolean;
	comment?: string;
	onDecline: Function;
	onApprove: Function;
};

export const ZFeedbackModal: FC<Props> = ({ onApprove, onDecline }: Props) => {
	const [showLoader, setShowLoader] = useState<boolean>(false);
	const [showConfirmation, setConfirmation] = useState<boolean>(false);

	return (
		<>
			{!showLoader && (
				<div className="zlifecycle-feedback-modal__actions">
					{showConfirmation ? (
						<div className="zlifecycle-feedback-modal__confirmation-popup">
							Do you want to Proceed?
							<Button
								color="primary"
								onClick={() => {
									setShowLoader(true);
									onApprove();
								}}>
								OK
							</Button>
						</div>
					) : (
						<Button
							color="primary"
							onClick={() => {
								setConfirmation(true);
							}}>
							Approve
						</Button>
					)}
				</div>
			)}
			{showLoader && <Loader height={20} width={20} />}
		</>
	);
};
