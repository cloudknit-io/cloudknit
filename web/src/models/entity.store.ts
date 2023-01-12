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
		environments: []
	});
  public emitterComp = new Subject<Component[]>();
	private emitterCompAudit = new Subject<CompAuditData>();
	private emitterEnvAudit = new Subject<EnvAuditData>();
	private componentAuditListeners = new Set<number>();
	private environmentAuditListeners = new Set<number>();

	public get Teams() {
		return [...this.teamMap.values()];
	}

	public get Environments() {
		return [...this.envMap.values()];
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
		});
	}

	private async getTeams() {
		const teams = await this.entityService.getTeams();
		teams.forEach(e => this.teamMap.set(e.id, e));
		await this.getEnvironments();
	}

	private async getEnvironments() {
		const envCalls = this.Teams.map(team =>
			this.entityService.getEnvironments(team.id)
		);
		const resps = await Promise.all(envCalls);
		resps.flat().forEach(e => this.envMap.set(e.id, {
			...e,
			argoId: `${this.getTeam(e.teamId)?.name}-${e.name}`
		}));
		this.emit();
	}

	private startStreaming() {
		this.entityService.streamEnvironments().subscribe((environment: Environment) => {
			const present = this.envMap.has(environment.id);
			
			if (present) {
				const currEnv = this.envMap.get(environment.id);
				this.envMap.set(environment.id, {
					...currEnv,
					...environment
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
					...component
				});
				this.emitterComp.next(this.getComponentsByEnvId(component.envId));
			}
		});

		this.entityService.streamAudit().subscribe((data: CompAuditData | EnvAuditData) => {
			if ((data as CompAuditData).compId) {
				const compData = data as CompAuditData;
				this.componentAuditListeners.has(compData.compId) && this.emitterCompAudit.next(compData);
			}
			if ((data as EnvAuditData).envId) {
				const envData = data as EnvAuditData;
				this.environmentAuditListeners.has(envData.envId) && this.emitterEnvAudit.next(envData);
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

	public getComponentsByEnvId(envId: number) : Component[] {
		return [...this.compMap.values()].filter(e => e.envId === envId);
	}

	public async getComponents(teamId: number, envId: number) {
		const components = await this.entityService.getComponents(teamId, envId);
    const currEnv = this.getEnvironmentById(envId);

		if (components.length > 0 && currEnv) {
			components.forEach(c => {
				const compDag = currEnv.dag.find(d => d.name === c.name);
				
				c.dependsOn = compDag?.dependsOn || [];
        c.argoId = `${currEnv.argoId}-${c.name}`;
				c.teamId = currEnv.teamId;

				this.compMap.set(c.id, c);
			});
		}
		
		this.emitterComp.next(this.getComponentsByEnvId(envId))
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
};

export type Update = {
	teams: Team[];
	environments: Environment[];
};

export type AuditData = {
	reconcileId: number;
	duration: number;
	status: AuditStatus;
	startDateTime: string;
	operation?: string;
	approvedBy?: string;
}

export type EnvAuditData = {
	envId: number;
} & AuditData

export type CompAuditData = {
	compId: number;
} & AuditData
