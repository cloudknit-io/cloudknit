import { EntityService } from 'services/entity/entity.service';
import { BehaviorSubject } from 'rxjs';

export class EntityStore {
	private static instance: EntityStore;
	private entityService = EntityService.getInstance();
	public emitter = new BehaviorSubject<Team[]>([]);
	private envEmitters = new Map<string, BehaviorSubject<Environment[]>>();

	private teams: Team[] = [];

	private constructor() {
		this.emitter.subscribe(data => {
			data.forEach(e =>
				(this.envEmitters.get(e.name) as BehaviorSubject<Environment[]>)?.next(e?.environments || [])
			);
		});
		Promise.resolve(this.getTeams());
	}

	static getInstance() {
		if (!EntityStore.instance) {
			EntityStore.instance = new EntityStore();
		}
		return EntityStore.instance;
	}

	private async getTeams() {
		this.teams = await this.entityService.getTeams();
		await this.getEnvironments();
	}

	private async getEnvironments() {
		const envCalls = this.teams.map(team => this.entityService.getEnvironments(team.id));
		const resps = await Promise.all(envCalls);
		resps.forEach((resp, _i) => {
            resp.forEach(r => r.teamId = this.teams[_i].id);
            this.teams[_i].environments = resp;
		});
		this.emitter.next(this.teams);
	}

    public getEnvironmentsEmitter(teamName: string) {
		if (!this.envEmitters.has(teamName)) {
			this.envEmitters.set(teamName, new BehaviorSubject<Environment[]>([]));
		}
		const emitter = this.envEmitters.get(teamName) as BehaviorSubject<Environment[]>;
        emitter.next(this.getTeam(teamName)?.environments || []);
        return emitter;
	};

    public getAllEnvironments(team: Team[]): Environment[] {
        return team.map(t => t.environments).flat();
    }

    public getTeam(id: number | string) {
        return this.teams.find(e => e.id === id || e.name === id);
    }
}

export type Team = {
	id: number;
	name: string;
	cost?: number;
	environments: Environment[];
};

export type Environment = {
	id: number;
	name: string;
	lastReconcileDatetime: Date;
	duration: number;
	dag: Component[];
    teamId: number;
};

export type Component = {
	name: string;
	type: string;
	dependsOn: string[];
};
