import { Column, Entity, OneToMany, PrimaryColumn, PrimaryGeneratedColumn } from 'typeorm'
import { ComponentReconcile } from './component-reconcile.entity'

@Entity({
  name: 'environment_reconcile',
})
export class EnvironmentReconcile {
  @PrimaryGeneratedColumn()
  reconcile_id?: number

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
}
