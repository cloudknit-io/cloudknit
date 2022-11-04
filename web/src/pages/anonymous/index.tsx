import './Anonymous.scss';

import React from 'react';

const Anonymous: React.FC = ({ children }) => {
	return <div className="main-wrapper">{children}</div>;
};

export default Anonymous;
