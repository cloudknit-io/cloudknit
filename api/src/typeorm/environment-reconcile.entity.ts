import { EnvSpecComponentDto } from 'src/environment/dto/env-spec.dto';
import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  OneToMany, PrimaryGeneratedColumn,
  RelationId
} from 'typeorm';
import { ComponentReconcile } from './component-reconcile.entity';
import { Environment } from './environment.entity';
import { ColumnNumericTransformer } from './helper';
import { Organization } from './Organization.entity';
import { Team } from './team.entity';

@Entity({
  name: 'environment_reconcile',
})
export class EnvironmentReconcile {
  @PrimaryGeneratedColumn({
    name: 'id',
  })
  reconcileId: number;

  @Column()
  status: string;

  @Column()
  gitSha: string;

  @Column({
    type: 'datetime',
  })
  startDateTime: string;

  @Column({
    type: 'datetime',
    nullable: true,
  })
  endDateTime?: string;

  @Column({
    name: 'estimated_cost',
    type: 'decimal',
    precision: 10,
    scale: 3,
    default: 0,
    transformer: new ColumnNumericTransformer(),
  })
  estimatedCost: number;

  @Column({
    type: 'json',
    default: null,
  })
  dag: EnvSpecComponentDto[];

  @Column({
    default: null,
    type: 'json',
  })
  errorMessage: string[];

  @OneToMany(
    () => ComponentReconcile,
    (component) => component.environmentReconcile,
    {
      eager: true,
      cascade: true,
    }
  )
  componentReconciles?: ComponentReconcile[];

  @ManyToOne(() => Environment, (environment) => environment.id)
  environment: Environment;

  @ManyToOne(() => Team, (team) => team.id)
  team: Team;

  @ManyToOne(() => Organization, (org) => org.id, {
    onDelete: 'CASCADE',
    eager: true
  })
  @JoinColumn({
    referencedColumnName: 'id',
  })
  organization: Organization;

  @RelationId((envRecon: EnvironmentReconcile) => envRecon.environment)
  envId: number;

  @RelationId((envRecon: EnvironmentReconcile) => envRecon.organization)
  orgId: number;

  @RelationId((envRecon: EnvironmentReconcile) => envRecon.team)
  teamId: number;
}
