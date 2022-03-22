import { Column, Entity, JoinColumn, ManyToOne, UpdateDateColumn } from "typeorm";
import { Environment } from "./environment.entity";

@Entity({
  name: "component",
})
export class Component {
  @Column({
    name: "component_name",
  })
  componentName: string;

  @UpdateDateColumn({
    name: 'last_reconcile_datetime'
  })
  lastReconcileDatetime: string;

  @Column()
  duration: number;

  @ManyToOne((type) => Environment, (environment) => environment.components)
  environment: Environment;
}
