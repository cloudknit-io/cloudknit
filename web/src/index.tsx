import './assets/styles/main.scss';

import * as Sentry from '@sentry/react';
import { Integrations } from '@sentry/tracing';
import App from 'App';
import { AuthProvider } from 'context/auth-provider/AuthProvider';
import { createRoot } from 'react-dom/client';

import * as serviceWorker from './serviceWorker';

if (process.env.REACT_APP_SENTRY_DSN) {
	Sentry.init({
		dsn: process.env.REACT_APP_SENTRY_DSN,
		environment: process.env.REACT_APP_SENTRY_ENVIRONMENT?.toLowerCase() ?? 'unknown',
		integrations: [new Integrations.BrowserTracing()],
		sampleRate: 1.0,
		tracesSampleRate: 0.0,
	});
}

const root = document.getElementById('root') as HTMLElement;
createRoot(root).render(
	<AuthProvider>
		<App />
	</AuthProvider>
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
