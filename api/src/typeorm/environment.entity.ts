import { Column, Entity, Index, JoinColumn, ManyToOne, OneToMany, PrimaryGeneratedColumn, UpdateDateColumn } from "typeorm";
import { Organization } from "./Organization.entity";
import { Component } from "./component.entity";
import { Team } from "./team.entity";
import { DagDto } from "src/reconciliation/dtos/environment-dag.dto";

@Entity({
  name: "environment",
})
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

  @Column()
  duration: number;

  @OneToMany(() => Component, (component) => component.environment)
  components: Component[];

  @Column({
    type: 'json',
    default: null
  })
  dag: DagDto;

  @ManyToOne(() => Team, (team) => team.id, {
    onDelete: "CASCADE"
  })
  @JoinColumn({
    referencedColumnName: 'id'
  })
  team: Team

  @ManyToOne(() => Organization, (org) => org.id, {
    onDelete: "CASCADE"
  })
  @JoinColumn({
    referencedColumnName: 'id'
  })
  organization: Organization
}
