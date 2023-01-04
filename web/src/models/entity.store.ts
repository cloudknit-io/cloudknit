import { EntityService } from 'services/entity/entity.service';
import { BehaviorSubject } from 'rxjs';

export class EntityStore {
	private static instance: EntityStore;
	private entityService = EntityService.getInstance();
	private teamMap = new Map<number, Team>();
	public emitter = new BehaviorSubject<Team[]>([]);

	private teams: Team[] = [];
	private envs: Environment[] = [];

    public get Teams() {
        return this.teams;
    }

    public get Environments() {
        return this.envs;
    }

	private constructor() {
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
		this.teams.forEach(e => this.teamMap.set(e.id, e));
		await this.getEnvironments();
	}

	private async getEnvironments() {
		const envCalls = this.teams.map(team =>
			this.entityService.getEnvironments(team.id).then(e => ({ data: e, teamId: team.id }))
		);
		const resps = await Promise.all(envCalls);
        
        resps.forEach(r => {
            if (this.teamMap.has(r.teamId)) {
                (this.teamMap.get(r.teamId) as Team).environments = r.data;
            }
        });
        
        this.envs = this.teams.map(e => e.environments).flat();
		this.emitter.next(this.teams);
	}

    public getAllEnvironmentsByTeamName(teamName: string): Environment[] {
        const teamId = this.teams.find(e => e.name === teamName)?.id;
        if (teamId && this.teamMap.has(teamId)) {
            return this.getTeam(teamId)?.environments as Environment[];
        }
		return [];
	}

    public getEnvironment(teamName: string, envName: string): Environment | null | undefined {
        const teamId = this.teams.find(e => e.name === teamName)?.id;
        if (!teamId || !this.teamMap.has(teamId)) return null;
        return this.getTeam(teamId)?.environments.find(e => e.name === envName);
	}

	public getTeam(id: number) {
		return this.teamMap.get(id);
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
