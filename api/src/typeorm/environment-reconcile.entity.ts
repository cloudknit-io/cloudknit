import { Column, Entity, JoinColumn, ManyToOne, OneToMany, PrimaryColumn, PrimaryGeneratedColumn } from 'typeorm'
import { Organization } from './Organization.entity'
import { ComponentReconcile } from './component-reconcile.entity'
import { Environment } from './environment.entity'
import { Team } from './team.entity'

@Entity({
  name: 'environment_reconcile',
})
export class EnvironmentReconcile {
  @PrimaryGeneratedColumn({
    name: 'id'
  })
  reconcileId: number

  @Column()
  status: string

  @Column({
      type: 'datetime'
  })
  start_date_time: string

  @Column({
    type: 'datetime',
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

  @ManyToOne(() => Team, (team) => team.id)
  team: Team;

  @ManyToOne(() => Organization, (org) => org.id, {
    onDelete: "CASCADE"
  })
  @JoinColumn({
    referencedColumnName: 'id'
  })
  organization: Organization
}
