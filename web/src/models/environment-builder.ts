import { uniqueId } from 'lodash';
import 'reflect-metadata';

const ignoreKey = Symbol('Ignore');
const valueResolver = Symbol('ValueResolver');
const setResolver = <T>(val: Set<T>) => (val.size > 0 ? [...val.values()] : undefined);
const dependsOnResolver = (val: Map<string, YComponent>) =>
	val.size > 0 ? [...val.values()].map(e => e.Name.value) : undefined;
const mapResolver = (val: Map<any, any>) => (val.size > 0 ? [...val.values()] : undefined);
const boolResolver = (val: Set<any>) => val || undefined;
const emptyStringResolver = (val: string) => val.trim() || undefined;
const nullResolver = (val: any) => (val === null ? undefined : val);

export type PropertyBinding = {
	value: any;
	label: string;
	set: (val: any) => void;
	key: string;
};

export type Nullable<T> = T | null;

export class YMetaData {
	name = 'zmart-checkout-demo';
	readonly namespace: string = 'zlifecycle';
	constructor(teamName: string, environmentName: string) {
		this.name = `zmart-${teamName}-${environmentName}`;
	}
}

export class YSpec {
	private teamName = '';
	private envName = '';
	@ValueResolver(nullResolver)
	private teardown: Nullable<boolean> = null;
	@ValueResolver(nullResolver)
	private autoApprove: Nullable<boolean> = null;
	@ValueResolver(mapResolver)
	components: Map<string, YComponent> = new Map<string, YComponent>();
	constructor(teamName: string, environmentName: string) {
		this.teamName = teamName;
		this.envName = environmentName;
	}

	@Ignore(true)
	get TeamName(): PropertyBinding {
		return {
			value: this.teamName,
			label: 'Team Name',
			set: (val: any) => (this.teamName = val),
			key: 'teamName',
		};
	}

	@Ignore(true)
	get EnvName(): PropertyBinding {
		return {
			value: this.envName,
			label: 'Environment Name',
			set: (val: any) => (this.envName = val),
			key: 'envName',
		};
	}

	@Ignore(true)
	get AutoApprove(): PropertyBinding {
		return {
			value: this.autoApprove,
			label: 'Auto Approve',
			set: (val: any) => (this.autoApprove = val),
			key: 'autoApprove',
		};
	}

	@Ignore(true)
	get Teardown(): PropertyBinding {
		return {
			value: this.teardown,
			label: 'Teardown',
			set: (val: any) => (this.teardown = val),
			key: 'teardown',
		};
	}
}

export class YModule {
	source: string;
	private name: string;
	@ValueResolver(emptyStringResolver)
	private path: string;
	constructor(source: string, name: string, path: string = '') {
		this.source = source;
		this.name = name;
		this.path = path;
	}

	@Ignore(true)
	get Name(): PropertyBinding {
		return {
			value: this.name,
			label: 'Module Source Name',
			set: (val: any) => (this.name = val),
			key: 'moduleName',
		};
	}

	@Ignore(true)
	get Path(): PropertyBinding {
		return {
			value: this.path,
			label: 'Module Source Path',
			set: (val: any) => (this.path = val),
			key: 'modulePath',
		};
	}
}

export class YVariableFile {
	private source: string;
	private path: string;
	@Ignore(true)
	id = uniqueId('YVariableFile');
	constructor(source: string, path: string) {
		this.source = source;
		this.path = path;
	}

	@Ignore(true)
	get Source(): PropertyBinding {
		return {
			value: this.source,
			label: 'Source',
			set: (val: any) => (this.source = val),
			key: 'source',
		};
	}

	@Ignore(true)
	get Path(): PropertyBinding {
		return {
			value: this.path,
			label: 'Path',
			set: (val: any) => (this.path = val),
			key: 'path',
		};
	}
}

export class YSecret {
	@Ignore(true)
	id = uniqueId('ytag');
	private name: string;
	private key: string;
	private scope: string;
	constructor(name: string = '', key: string = '', scope: string = '') {
		this.name = name;
		this.key = key;
		this.scope = scope;
	}

	@Ignore(true)
	get Name(): PropertyBinding {
		return {
			key: 'Name',
			value: this.name,
			set: (val: any) => (this.name = val),
			label: 'Secret Name',
		};
	}

	@Ignore(true)
	get Key(): PropertyBinding {
		return {
			key: 'Key',
			value: this.key,
			set: (val: any) => (this.key = val),
			label: 'Secret Key',
		};
	}

	@Ignore(true)
	get Scope(): PropertyBinding {
		return {
			key: 'Scope',
			value: this.scope,
			set: (val: any) => (this.scope = val),
			label: 'Secret Scope',
		};
	}
}

export class YVariable {
	private name: string;
	@ValueResolver(emptyStringResolver)
	private value: string;
	@ValueResolver(emptyStringResolver)
	private valueFrom: string;
	@Ignore(true)
	public choosenDisposition: 'Value' | 'ValueFrom' = 'Value';
	constructor(name: string = '', value: string = '', valueFrom: string = '') {
		this.name = name;
		this.value = value;
		this.valueFrom = valueFrom;
	}
	@Ignore(true)
	get Name(): PropertyBinding {
		return {
			key: 'Name',
			value: this.name,
			set: (val: any) => (this.name = val),
			label: 'Name',
		};
	}

	@Ignore(true)
	private valuePropertyBinding = {
		key: 'Value',
		value: () => this.value,
		set: (val: any) => {
			this.valueFrom = '';
			this.value = val;
			this.ValueFrom.key = uniqueId('valueFrom-');
		},
		label: 'Value',
	};

	@Ignore(true)
	get Value(): PropertyBinding {
		return {
			key: 'Value',
			value: this.value,
			set: (val: any) => {
				this.valueFrom = '';
				this.value = val;
			},
			label: 'Value',
		};
	}

	@Ignore(true)
	get ValueFrom(): PropertyBinding {
		return {
			key: 'valueFrom',
			value: this.valueFrom,
			set: (val: any) => {
				this.value = '';
				this.valueFrom = val;
			},
			label: 'Value From',
		};
	}
}

export class YTag {
	@Ignore(true)
	id = uniqueId('ytag');
	private name: string;
	private value: string;
	constructor(name: string = '', value: string = '') {
		this.name = name;
		this.value = value;
	}

	@Ignore(true)
	get Name(): PropertyBinding {
		return {
			key: 'Name',
			value: this.name,
			set: (val: any) => (this.name = val),
			label: 'Tag Name',
		};
	}

	@Ignore(true)
	get Value(): PropertyBinding {
		return {
			key: 'Value',
			value: this.value,
			set: (val: any) => (this.value = val),
			label: 'Tag Value',
		};
	}
}

export class YOutput {
	name: string;
	@ValueResolver(nullResolver)
	sensitive: Nullable<boolean>;
	constructor(name: string = '', sensitive: Nullable<boolean> = null) {
		this.name = name;
		this.sensitive = sensitive;
	}
}

export class YComponent {
	private name: string;
	@ValueResolver(nullResolver)
	private autoApprove: Nullable<boolean> = null;
	type = 'terraform';
	module: YModule;
	variablesFile: YVariableFile | undefined;
	@ValueResolver(dependsOnResolver)
	dependsOn: Map<string, YComponent>;
	@ValueResolver(setResolver)
	secrets: Set<YSecret> = new Set<YSecret>();
	@ValueResolver(mapResolver)
	variables: Map<string, YVariable> = new Map<string, YVariable>();
	@ValueResolver(setResolver)
	outputs: Map<string, YOutput> = new Map<string, YOutput>();
	@ValueResolver(setResolver)
	tags: Set<YTag> = new Set<YTag>();
	@ValueResolver(nullResolver)
	destroyProtection: Nullable<boolean> = null;
	@Ignore(true)
	id: string = uniqueId();
	@Ignore(true)
	outputList: Set<string> = new Set<string>();
	@Ignore(true)
	selectedOutputs: Set<string> = new Set<string>();
	//TODO: Change implementation with new DDL
	@Ignore(true)
	inputList: Map<string, any> = new Map<string, any>();
	@Ignore(true)
	selectedInputs: Map<string, any> = new Map<string, any>();
	@Ignore(true)
	get Name(): PropertyBinding {
		return {
			value: this.name,
			label: 'Name',
			set: (val: any) => (this.name = val),
			key: 'name',
		};
	}
	@Ignore(true)
	get AutoApprove(): PropertyBinding {
		return {
			value: this.autoApprove,
			label: 'Auto Approve',
			set: (val: any) => {
				this.autoApprove = val;
			},
			key: 'autoApprove',
		};
	}

	@Ignore(true)
	get DestroyProtection(): PropertyBinding {
		return {
			value: this.destroyProtection,
			label: 'Destroy Protection',
			set: (val: any) => {
				this.destroyProtection = val;
			},
			key: 'destroyProtection',
		};
	}

	constructor(
		name: string = '',
		source: string = '',
		sourceName: string = '',
		variableFile: string = '',
		path: string = ''
	) {
		this.name = name || `node-${this.id}`;
		this.module = new YModule(source, sourceName);
		this.variablesFile = new YVariableFile(variableFile, path);
		this.dependsOn = new Map<string, YComponent>();
	}
}

export class YEnvironmentBuilder {
	apiVersion = 'stable.compuzest.com/v1';
	kind = 'Environment';
	metadata: YMetaData;
	spec: YSpec;
	@Ignore(true)
	id: string = uniqueId();
	@Ignore(true)
	repository: string;

	constructor(teamName: string = '{team}', environmentName: string = '{env}', repository: string = '') {
		this.metadata = new YMetaData(teamName, environmentName);
		this.spec = new YSpec(teamName, environmentName);
		this.repository = repository;
		this.replacer = this.replacer;
		this.yamlObject = this.yamlObject;
	}

	private replacer(key: string, value: any) {
		const decorators = Reflect.getMetadataKeys(this, key);
		if (decorators.includes(ignoreKey)) {
			return undefined;
		}

		if (decorators.includes(valueResolver)) {
			return Reflect.getMetadata(valueResolver, this, key)(value);
		}

		if (this.constructor.name === YVariable.name) {
			if ((key === 'value' || key === 'valueFrom') && value.length === 0) {
				return undefined;
			}
		}

		return value;
	}

	yamlObject() {
		this.spec.components.forEach(e => {
			if (e.variables.size > 0) {
				//if (e.variables.size > 0) {
				// e.variablesFile = undefined;
			} else {
				// if (!e.variablesFile) {
				// 	e.variablesFile.Source.set(this.repository);
				// 	e.variablesFile.path = `${this.spec.EnvName.value}/${e.Name.value}.tfvars`;
				// }
			}
			e.module.source = 'aws';
		});
		return JSON.parse(JSON.stringify(this, this.replacer));
	}
}

function Ignore(value: boolean): any {
	return Reflect.metadata(ignoreKey, value);
}

function ValueResolver(resolver: (val: any) => any): any {
	return Reflect.metadata(valueResolver, resolver);
}
