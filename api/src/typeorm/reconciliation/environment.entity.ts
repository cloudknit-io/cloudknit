import { Column, Entity, OneToMany, UpdateDateColumn } from "typeorm";
import { Component } from "./component.entity";

@Entity({
  name: "environment",
})
export class Environment {
 
  @Column({
    name: "environment_name",
    primary: true,
  })
  environmentName: string;

  @UpdateDateColumn({
    name: 'last_reconcile_datetime'
  })
  lastReconcileDatetime: string;

  @Column()
  duration: number;

  @OneToMany(type => Component, component => component.environment) components: Component[]
}
