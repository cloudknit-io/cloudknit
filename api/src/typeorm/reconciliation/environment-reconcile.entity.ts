import { Column, Entity, JoinColumn, ManyToOne, OneToMany, PrimaryColumn, PrimaryGeneratedColumn } from 'typeorm'
import { Organization } from '../Organization.entity'
import { ComponentReconcile } from './component-reconcile.entity'
import { Environment } from './environment.entity'

@Entity({
  name: 'environment_reconcile',
})
export class EnvironmentReconcile {
  @PrimaryGeneratedColumn()
  reconcile_id: number

  @Column()
  name: string

  @Column()
  team_name: string

  @Column()
  status: string

  @Column({
      type: 'datetime'
  })
  start_date_time: string

  @Column({
      nullable: true
  })
  end_date_time?: string

  @OneToMany(() => ComponentReconcile, (component) => component.environmentReconcile, {
    eager: true,
    cascade: true,
  })
  componentReconciles?: ComponentReconcile[]

  @ManyToOne(() => Environment, (environment) => environment.id, {
    eager: true
  })
  environment: Environment;

  @ManyToOne(() => Organization, (org) => org.id, {
    onDelete: "CASCADE"
  })
  @JoinColumn({
    referencedColumnName: 'id'
  })
  organization: Organization
}
