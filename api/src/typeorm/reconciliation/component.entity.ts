import { Column, Entity, Index, ManyToOne, PrimaryGeneratedColumn, UpdateDateColumn } from "typeorm";
import { Environment } from "./environment.entity";

@Entity({
  name: "component",
})
export class Component {
  @PrimaryGeneratedColumn()
  id: number

  @Column({
    name: "component_name"
  })
  @Index()
  componentName: string;

  @UpdateDateColumn({
    name: 'last_reconcile_datetime'
  })
  lastReconcileDatetime: string;

  @Column()
  duration: number;

  @ManyToOne(() => Environment, (environment) => environment.components, {
    eager: true
  })
  environment: Environment;
}
