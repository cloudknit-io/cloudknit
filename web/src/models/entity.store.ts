import { EntityService } from 'services/entity/entity.service';
import { BehaviorSubject, Subject, Subscription } from 'rxjs';
import { ErrorStateService } from 'services/error/error-state.service';
import { CompAuditData, Component, EnvAuditData, Environment, StreamTypeEnum, Team, Update } from './entity.type';
import { ZSyncStatus } from './argo.models';

export class EntityStore {
	private static instance: EntityStore;
	private entityService = EntityService.getInstance();
	private teamMap = new Map<number, Team>();
	private compMap = new Map<number, Component>();
	private envMap = new Map<number, Environment>();
	public emitter = new BehaviorSubject<Update>({
		teams: [],
		environments: [],
		components: [],
	});
	public emitterComp = new Subject<Component[]>();
	public emitterCompAudit = new Subject<CompAuditData>();
	private emitterEnvAudit = new Subject<EnvAuditData>();
	private emitterMap = new Map<StreamTypeEnum, Function>();
	private componentAuditListeners = new Set<number>();
	private environmentAuditListeners = new Set<number>();
	private allDataFetched = false;

	public get Teams() {
		return [...this.teamMap.values()];
	}

	public get Environments() {
		return [...this.envMap.values()].map(e => {
			if (e.errorMessage) {
				e.status = ZSyncStatus.ValidationFailed;
			}
			return e;
		});
	}

	public get Components() {
		return [...this.compMap.values()];
	}

	public get AllDataFetched() {
		return this.allDataFetched;
	}

	private constructor() {
		ErrorStateService.getInstance();
		this.generateEmitterMap();
		Promise.resolve(this.getTeams());
		this.startStreaming();
	}

	static getInstance() {
		if (!EntityStore.instance) {
			EntityStore.instance = new EntityStore();
		}
		return EntityStore.instance;
	}

	private generateEmitterMap() {
		this.emitterMap.set(StreamTypeEnum.Component, this.streamComponent.bind(this));
		this.emitterMap.set(StreamTypeEnum.ComponentReconcile, this.streamAudit.bind(this));
		this.emitterMap.set(StreamTypeEnum.EnvironmentReconcile, this.streamAudit.bind(this));
		this.emitterMap.set(StreamTypeEnum.Environment, this.streamEnvironment.bind(this));
		this.emitterMap.set(StreamTypeEnum.Team, this.streamTeam.bind(this));
	}

	private emit() {
		this.emitter.next({
			teams: this.Teams,
			environments: this.Environments,
			components: this.Components,
		});
	}

	public async getTeams(withComponents: boolean = false) {
		if (!this.allDataFetched) {
			const teams = await this.entityService.getTeams(withComponents);
			teams.forEach(e => {
				this.teamMap.set(e.id, e);
				e.environments.forEach(env => {
					this.envMap.set(env.id, {
						...env,
						...this.mergeEnvReconToEnv(env, env.latestEnvRecon),
						argoId: `${e.name}-${env.name}`,

					});
					env.components?.forEach(c => {
						const compDag = env.dag.find(d => d.name === c.name);
						this.compMap.set(c.id, {
							...(this.compMap.has(c.id) ? this.compMap.get(c.id) : {}),
							...c,
							argoId: `${e.name}-${env.name}-${c.name}`,
							dependsOn: compDag?.dependsOn || [],
							teamId: e.id,
						});
					});
				});
			});
		}
		this.allDataFetched = this.allDataFetched || withComponents;
		this.emit();
	}

	private startStreaming() {
		this.entityService
			.stream([
				StreamTypeEnum.Component,
				StreamTypeEnum.ComponentReconcile,
				StreamTypeEnum.Empty,
				StreamTypeEnum.Environment,
				StreamTypeEnum.EnvironmentReconcile,
				StreamTypeEnum.Team,
			])
			.subscribe(({ data, type }) => {
				if (this.emitterMap.has(type)) {
					this.emitterMap.get(type)?.call(this, data);
				}
			});
	}

	private streamTeam(team: Team) {
		if (this.teamMap.has(team.id)) {
			const currTeam = this.teamMap.get(team.id);
			this.teamMap.set(team.id, {
				...currTeam,
				...team,
			});
		} else {
			this.teamMap.set(team.id, team);
		}
		this.emit();
	}

	private streamEnvironment(environment: Environment) {
		if (this.envMap.has(environment.id)) {
			const currEnv = this.envMap.get(environment.id);
			this.envMap.set(environment.id, {
				...currEnv,
				...environment,
				...this.mergeEnvReconToEnv(environment, environment.latestEnvRecon)
			});
		} else {
			this.envMap.set(environment.id, {
				...environment,
				...this.mergeEnvReconToEnv(environment, environment.latestEnvRecon),
				argoId: `${this.getTeam(environment.teamId)?.name}-${environment.name}`,
			});
		}

		this.emit();
	}

	private streamComponent(component: Component) {
		const present = this.compMap.has(component.id);
		if (present) {
			const currComp = this.compMap.get(component.id);
			this.compMap.set(component.id, {
				...currComp,
				...component,
			});
		} else {
			this.compMap.set(component.id, this.mapComponent(component));
		}
		this.emitterComp.next(this.getComponentsByEnvId(component.envId));
	}

	private streamAudit(data: CompAuditData | EnvAuditData) {
		if ((data as CompAuditData).compId) {
			const compData = data as CompAuditData;
			this.emitterCompAudit.next(compData);
		}
		if ((data as EnvAuditData).envId) {
			const envData = data as EnvAuditData;
			this.emitterEnvAudit.next(envData);
		}
	}

	private mapComponent(component: Component) {
		const currEnv = this.getEnvironmentById(component.envId);
		if (currEnv) {
			const compDag = currEnv.dag.find(d => d.name === component.name);
			component.dependsOn = compDag?.dependsOn || [];
			component.argoId = `${currEnv.argoId}-${component.name}`;
			component.teamId = currEnv.teamId;
		}
		return component;
	}

	private mergeEnvReconToEnv(env: Environment, envRecon: EnvAuditData): Environment {
		if (!envRecon) return env;
		env.estimatedCost = envRecon.estimatedCost;
		env.dag = envRecon.dag;
		env.errorMessage = envRecon.errorMessage;
		env.lastReconcileDatetime = envRecon.startDateTime;
		env.status = envRecon.status;
		return env;
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
		const cachedComps = this.getComponentsByEnvId(envId);
		if (cachedComps.length === this.getEnvironmentById(envId)?.dag?.length) {
			this.emitterComp.next(this.getComponentsByEnvId(envId));
			return cachedComps;
		}

		const components: any = await this.entityService.getComponents(teamId, envId, true);
		const currEnv = this.getEnvironmentById(envId);

		if (components.length > 0 && currEnv) {
			components.forEach((c: Component) => {
				const component = this.mapComponent(c);
				this.compMap.set(component.id, component);
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
