import {
  AfterInsert,
  Column,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  OneToMany,
  PrimaryGeneratedColumn,
  RelationId,
  UpdateDateColumn,
} from "typeorm";
import { Organization } from "./Organization.entity";
import { Component } from "./component.entity";
import { Team } from "./team.entity";
import { EnvSpecComponentDto } from "src/environment/dto/env-spec.dto";
import { ColumnNumericTransformer } from "./helper";

@Entity({
  name: "environment",
})
@Index(["organization", "team", "name"], { unique: true })
export class Environment {
  @PrimaryGeneratedColumn()
  id: number;

  @Column()
  @Index()
  name: string;

  @UpdateDateColumn({
    name: "last_reconcile_datetime",
  })
  lastReconcileDatetime: string;

  @Column({
    default: -1,
  })
  duration: number;

  @Column({
    name: "status",
    default: "initializing",
  })
  status: string;

  @OneToMany(() => Component, (component) => component.environment)
  components: Component[];

  @Column({
    name: "estimated_cost",
    type: "decimal",
    precision: 10,
    scale: 3,
    default: 0,
    transformer: new ColumnNumericTransformer(),
  })
  estimatedCost: number;

  @Column({
    type: "json",
    default: null,
  })
  dag: EnvSpecComponentDto[];

  @Column({
    default: false,
  })
  isDeleted: boolean;

  @ManyToOne(() => Team, (team) => team.id, {
    onDelete: "CASCADE",
  })
  @JoinColumn({
    referencedColumnName: "id",
  })
  team: Team;

  @ManyToOne(() => Organization, (org) => org.id, {
    onDelete: "CASCADE",
  })
  @JoinColumn({
    referencedColumnName: "id",
  })
  organization: Organization;

  @RelationId((env: Environment) => env.team)
  teamId: number;

  @RelationId((env: Environment) => env.organization)
  orgId: number;
}
