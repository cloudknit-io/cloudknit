import { EnvSpecComponentDto } from 'src/environment/dto/env-spec.dto';
import {
  Column,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  OneToMany,
  OneToOne,
  PrimaryGeneratedColumn,
  RelationId
} from 'typeorm';
import { Component } from './component.entity';
import { EnvironmentReconcile } from './environment-reconcile.entity';
import { Organization } from './Organization.entity';
import { Team } from './team.entity';

@Entity({
  name: 'environment',
})
@Index(['organization', 'team', 'name'], { unique: true })
export class Environment {
  @PrimaryGeneratedColumn()
  id: number;

  @Column()
  @Index()
  name: string;

  @Column({
    name: 'last_reconcile_datetime',
    type: 'datetime',
  })
  lastReconcileDatetime: string;

  @OneToOne(() => EnvironmentReconcile, (envRecon) => envRecon.environment, {
    eager: true
  })
  @JoinColumn({
    referencedColumnName: 'reconcileId',
    name: 'latest_env_recon_id'
  })
  latestEnvRecon: EnvironmentReconcile

  @OneToMany(() => Component, (component) => component.environment)
  components: Component[];

  @Column({
    default: false,
  })
  isDeleted: boolean;

  @Column({
    type: 'json',
    default: null,
  })
  dag: EnvSpecComponentDto[];

  @ManyToOne(() => Team, (team) => team.id, {
    onDelete: 'CASCADE',
  })
  @JoinColumn({
    referencedColumnName: 'id',
  })
  team: Team;

  @ManyToOne(() => Organization, (org) => org.id, {
    onDelete: 'CASCADE',
  })
  @JoinColumn({
    referencedColumnName: 'id',
  })
  organization: Organization;

  @RelationId((env: Environment) => env.team)
  teamId: number;

  @RelationId((env: Environment) => env.organization)
  orgId: number;
}
