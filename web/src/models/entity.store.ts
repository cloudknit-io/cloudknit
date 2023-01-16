import { EntityService } from 'services/entity/entity.service';
import { BehaviorSubject, Subject, Subscription } from 'rxjs';
import { AuditStatus } from './argo.models';
import { ErrorStateService } from 'services/error/error-state.service';

export class EntityStore {
	private static instance: EntityStore;
	private entityService = EntityService.getInstance();
	private teamMap = new Map<number, Team>();
	private compMap = new Map<number, Component>();
	private envMap = new Map<number, Environment>();
	public emitter = new BehaviorSubject<Update>({
		teams: [],
		environments: [],
		components: []
	});
	public emitterComp = new Subject<Component[]>();
	public emitterCompAudit = new Subject<CompAuditData>();
	private emitterEnvAudit = new Subject<EnvAuditData>();
	private componentAuditListeners = new Set<number>();
	private environmentAuditListeners = new Set<number>();

	public get Teams() {
		return [...this.teamMap.values()];
	}

	public get Environments() {
		return [...this.envMap.values()];
	}

	public get Components() {
		return [...this.compMap.values()];
	}

	private constructor() {
		ErrorStateService.getInstance();
		Promise.resolve(this.getTeams());
		this.startStreaming();
	}

	static getInstance() {
		if (!EntityStore.instance) {
			EntityStore.instance = new EntityStore();
		}
		return EntityStore.instance;
	}

	private emit() {
		this.emitter.next({
			teams: this.Teams,
			environments: this.Environments,
			components: this.Components,
		});
	}

	public async getTeams(withComponents: boolean = false) {
		const teams = await this.entityService.getTeams(withComponents);
		teams.forEach(e => {
			this.teamMap.set(e.id, e);
			e.environments.forEach(env => {
				this.envMap.set(env.id, {
					...env,
					argoId: `${e.name}-${env.name}`,
				});
				env.components?.forEach(c => {
					const compDag = env.dag.find(d => d.name === c.name);
					this.compMap.set(c.id, {
						...(this.compMap.has(c.id) ? this.compMap.get(c.id) : {}),
						...c,
						argoId: `${e.name}-${env.name}-${c.name}`,
						dependsOn: compDag?.dependsOn || [],
						teamId: e.id
					});
				})
			});
		});
		this.emit();
	}

	private startStreaming() {
		this.entityService.streamEnvironments().subscribe((environment: Environment) => {
			const present = this.envMap.has(environment.id);

			if (present) {
				const currEnv = this.envMap.get(environment.id);
				this.envMap.set(environment.id, {
					...currEnv,
					...environment,
				});

				this.emit();
			}
		});

		this.entityService.streamComponents().subscribe((component: Component) => {
			const present = this.compMap.has(component.id);
			if (present) {
				const currComp = this.compMap.get(component.id);
				this.compMap.set(component.id, {
					...currComp,
					...component,
				});
				this.emitterComp.next(this.getComponentsByEnvId(component.envId));
			}
		});

		this.entityService.streamAudit().subscribe((data: CompAuditData | EnvAuditData) => {
			if ((data as CompAuditData).compId) {
				const compData = data as CompAuditData;
				this.emitterCompAudit.next(compData);
			}
			if ((data as EnvAuditData).envId) {
				const envData = data as EnvAuditData;
				this.emitterEnvAudit.next(envData);
			}
		});
	}

	public getTeam(id: number) {
		return this.teamMap.get(id);
	}

	public getAllEnvironmentsByTeamName(teamName: string): Environment[] {
		const teamId = this.Teams.find(e => e.name === teamName)?.id;
		return this.Environments.filter(env => env.teamId === teamId);
	}

	public getAllEnvironmentsByTeamId(teamId: number): Environment[] {
		return this.Environments.filter(env => env.teamId === teamId);
	}

	public getEnvironmentByName(teamName: string, envName: string): Environment | null | undefined {
		const teamId = this.Teams.find(e => e.name === teamName)?.id;
		return this.Environments.find(e => e.name === envName && e.teamId === teamId);
	}

	public getEnvironmentById(envId: number): Environment | null | undefined {
		return this.envMap.get(envId);
	}

	public getComponentsByEnvId(envId: number): Component[] {
		return [...this.compMap.values()].filter(e => e.envId === envId);
	}

	public async getComponents(teamId: number, envId: number) {
		const components: any = await this.entityService.getComponents(teamId, envId, true);
		const currEnv = this.getEnvironmentById(envId);

		if (components.length > 0 && currEnv) {
			components.forEach((c: any) => {
				const compDag = currEnv.dag.find(d => d.name === c['component_name']);
				const rd: Component = {
					costResources: c['cost_resources'],
					name: c['component_name'],
					duration: c['duration'],
					envId: c['environmentId'],
					estimatedCost: c['estimated_cost'],
					isDestroyed: c['isDestroyed'],
					lastReconcileDatetime: c['last_reconcile_datetime'],
					status: c['status'],
					id: c['id'],
					lastWorkflowRunId: c['lastWorkflowRunId'],
					lastAuditStatus: c['lastAuditStatus'],
					dependsOn: compDag?.dependsOn || [],
					argoId: `${currEnv.argoId}-${c.component_name}`,
					teamId: currEnv.teamId,
					changeId: Symbol(),
					type: c['type'],
				};

				// c.dependsOn = compDag?.dependsOn || [];
				// c.argoId = `${currEnv.argoId}-${rd.name}`;
				// c.teamId = currEnv.teamId;

				this.compMap.set(c.id, rd);
			});
		}

		this.emitterComp.next(this.getComponentsByEnvId(envId));
		return components;
	}

	public setComponentAuditLister(componentId: number) {
		this.componentAuditListeners.add(componentId);
		return this.emitterCompAudit;
	}

	public removeComponentAuditLister(componentId: number, sub: Subscription) {
		this.componentAuditListeners.delete(componentId);
		sub.unsubscribe();
	}

	public setEnvironmentAuditLister(envId: number) {
		this.environmentAuditListeners.add(envId);
		return this.emitterEnvAudit;
	}

	public removeEnvironmentAuditLister(envId: number, sub: Subscription) {
		this.environmentAuditListeners.delete(envId);
		sub.unsubscribe();
	}
}

export type Team = {
	id: number;
	name: string;
	estimatedCost?: number;
	environments: Environment[];
};

export type Environment = {
	argoId: string;
	id: number;
	name: string;
	lastReconcileDatetime: Date;
	duration: number;
	dag: DAG[];
	teamId: number;
	status: string;
	isDeleted: boolean;
	estimatedCost: number;
	components: Component[];
};

export type DAG = {
	name: string;
	type: string;
	dependsOn: string[];
};

export type Component = {
	changeId: Symbol;
	argoId: string;
	teamId: number;
	id: number;
	name: string;
	type: string;
	status: string;
	estimatedCost: number;
	lastReconcileDatetime: Date;
	duration: number;
	isDestroyed: boolean;
	costResources: any;
	dependsOn: string[];
	envId: number;
	lastWorkflowRunId: string;
	lastAuditStatus: AuditStatus;
};

export type Update = {
	teams: Team[];
	environments: Environment[];
	components: Component[]
};

export type AuditData = {
	reconcileId: number;
	duration: number;
	status: AuditStatus;
	startDateTime: string;
	operation?: string;
	approvedBy?: string;
};

export type EnvAuditData = {
	envId: number;
} & AuditData;

export type CompAuditData = {
	compId: number;
} & AuditData;
