import './style.scss';

import { ReactComponent as AWSIcon } from 'assets/images/icons/AWS.svg';
import { ZText } from 'components/atoms/text/Text';
import React, { FC } from 'react';

type Props = {
	labels?: any;
	items: any;
	title: string;
	teamName: string;
	envName: string;
	estimatedCost: string | JSX.Element;
	model: string;
	classNames: string;
	onClick: () => void;
};

export const ZModelCard: FC<Props> = ({
	title,
	teamName,
	model,
	envName,
	estimatedCost,
	items,
	labels,
	classNames,
	onClick,
}: Props) => {
	return (
		<div className={`com-card ${classNames}`} onClick={onClick}>
			<div className="com-card__header">
				<div className="com-card__header__title">
					<div>
						<h4>{title}</h4>
						<div className="com-card__component-descriptors">
							<div>
								<ZText.Body className="color-gray" size="14" lineHeight="18" weight="bold">
									Team
								</ZText.Body>
								<h5 className="color-gray">{teamName}</h5>
							</div>
							<div>
								<ZText.Body className="color-gray" size="14" lineHeight="18" weight="bold">
									Environment
								</ZText.Body>
								<h5 className="color-gray">{envName}</h5>
							</div>
							{model === 'Application' ? null : (
								<div>
									<ZText.Body className="color-gray" size="14" lineHeight="18" weight="bold">
										Est. Monthly Cost
									</ZText.Body>
									<h5 className="color-gray">
										{typeof estimatedCost === 'string' ? `$${estimatedCost}` : estimatedCost}
									</h5>
								</div>
							)}
							<div>{items}</div>
						</div>
					</div>
				</div>
				<div className="com-card__more-options">
					<AWSIcon />
				</div>
			</div>
			<div className="com-card__body">
				<div className="com-card__labels">{labels}</div>
			</div>
		</div>
	);
};
