import ApiClient from 'utils/apiClient';
export class SecretsService {
	private static instance: SecretsService | null = null;
	private constructUri = (path: string) => `/secrets/${path}`;
	private sanitizeScope = (scope: string) => scope.split('/').slice(1).join('/');

	static getInstance() {
		if (!SecretsService.instance) {
			SecretsService.instance = new SecretsService();
		}
		return SecretsService.instance;
	}

	updateAWSSecret(secrets: { key: string; value: string }[], scope: string): Promise<any> {
		
		const awsSecrets = secrets.map(s => ({
			path: `${this.sanitizeScope(scope.toLowerCase())}/${s.key}`,
			value: s.value,
		}));
		const url = this.constructUri(SecretsUriType.awsSecret);

		return ApiClient.post(url, {
			awsSecrets,
		});
	}

	secretsExists(secrets: string[], scope: string) {
		const pathNames = secrets.map(s => `${this.sanitizeScope(scope.toLowerCase())}/${s}`);
		const url = this.constructUri(SecretsUriType.existsAwsSecret);
	
		return ApiClient.post(url, {
			pathNames,
		});
	}

	async getSsmSecrets(recursive: boolean, path?: string) {
		const url = this.constructUri(SecretsUriType.getSSMSecrets);
		const resp = await ApiClient.post(url, {
			path: path ? this.sanitizeScope(path) : null,
			recursive
		});

		return resp;
	}

	getEnvironments(path: string) {
		const url = this.constructUri(SecretsUriType.getEnvironments);
		return ApiClient.post(url, { path });
	}

	putSsmSecret(path: string, value: string) {
		const url = this.constructUri(SecretsUriType.awsSecret);
		return ApiClient.post(url, {
			awsSecrets: [{
				path: this.sanitizeScope(path.toLowerCase()),
				value: value
			}],
		});
	}

	deleteSsmSecret(path: string) {
		const url = this.constructUri(SecretsUriType.deleteSSMSecrets(encodeURIComponent(path)));
		return ApiClient.delete(url);
	}
}

class SecretsUriType {
	static awsSecret = `update/aws-secret`;
	static existsAwsSecret = `exists/aws-secret`;
	static getSSMSecrets = `get/ssm-secrets`;
	static getEnvironments = `get/environments`;
	static deleteSSMSecrets = (path: string) => `delete/ssm-secret?path=${path}`;
}
