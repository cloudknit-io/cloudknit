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
	const [showLoader, setShowLoader] = useState<boolean>(false);

	return (
		<div className="zlifecycle-feedback-modal">
			{showLoader ? <Loader /> : (
				<>
					<div className="zlifecycle-feedback-modal__content">
						<ZText.Body size="14" lineHeight="18">
							In order to continue, please approve:
						</ZText.Body>
					</div>
					<div className="zlifecycle-feedback-modal__actions">
						<Button
							color="primary"
							onClick={() => {
								setShowLoader(true);
								onApprove();
							}}>
							Approve
						</Button>
					</div>
				</>
			)}
		</div>
	);
};
