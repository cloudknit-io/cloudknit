import { ReactComponent as Logo } from 'assets/images/icons/logo.svg';
import { ZText } from 'components/atoms/text/Text';
import React from 'react';
import { Link } from 'react-router-dom';

export const NotFound: React.FC = () => {
	return (
		<div className="anonymous-page">
			<div className="anonymous-page__section">
				<div className="anonymous-page__section__content">
					<div className="main-div-">
						<div className="panel d-flex align-center justify-center h-100 text-center">
							<div>
								<Logo style={{ width: '60px', marginBottom: '30px' }} />
								<ZText.Headline size="28" lineHeight="26">
									ZLIFECYCLE
								</ZText.Headline>
								<ZText.Body size="18" lineHeight="24">
									Page not found!
								</ZText.Body>
								<Link to="/dashboard">
									<ZText.Body size="18" lineHeight="24">
										Go back to dashboard
									</ZText.Body>
								</Link>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
};
