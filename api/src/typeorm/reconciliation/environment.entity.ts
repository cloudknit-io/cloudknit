import { Column, Entity, Index, JoinColumn, ManyToOne, OneToMany, PrimaryGeneratedColumn, UpdateDateColumn } from "typeorm";
import { Organization } from "../Organization.entity";
import { Component } from "./component.entity";

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

  @ManyToOne(() => Organization, (org) => org.id, {
    onDelete: "CASCADE"
  })
  @JoinColumn({
    referencedColumnName: 'id'
  })
  organization: Organization
}
