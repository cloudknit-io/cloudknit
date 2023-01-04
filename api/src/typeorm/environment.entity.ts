import { Column, Entity, Index, JoinColumn, ManyToOne, OneToMany, PrimaryGeneratedColumn, UpdateDateColumn } from "typeorm";
import { Organization } from "./Organization.entity";
import { Component } from "./component.entity";
import { Team } from "./team.entity";
import { EnvSpecComponentDto } from "src/environment/dto/env-spec.dto";

@Entity({
  name: "environment",
})
@Index(['organization', 'team', 'name'], { unique: true })
export class Environment {
  @PrimaryGeneratedColumn()
  id: number

  @Column()
  @Index()
  name: string;

  @UpdateDateColumn({
    name: "last_reconcile_datetime",
  })
  lastReconcileDatetime: string;

  @Column({
    default: -1
  })
  duration: number;

  @Column({
    name: 'status',
    default: null
  })
  status: string;

  @OneToMany(() => Component, (component) => component.environment)
  components: Component[];

  @Column({
    type: 'json',
    default: null
  })
  dag: EnvSpecComponentDto[];

  @Column({
    default: false
  })
  isDeleted: boolean;

  @ManyToOne(() => Team, (team) => team.id, {
    onDelete: "CASCADE"
  })
  @JoinColumn({
    referencedColumnName: 'id'
  })
  team: Team;

  @ManyToOne(() => Organization, (org) => org.id, {
    onDelete: "CASCADE"
  })
  @JoinColumn({
    referencedColumnName: 'id'
  })
  organization: Organization;
}
