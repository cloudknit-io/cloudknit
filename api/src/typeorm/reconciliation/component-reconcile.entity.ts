import { Column, Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn } from 'typeorm'
import { EnvironmentReconcile } from './environment-reconcile.entity'

@Entity({
  name: 'component_reconcile',
})
export class ComponentReconcile {
  @PrimaryGeneratedColumn()
  reconcile_id?: number

  @ManyToOne(() => EnvironmentReconcile, environmentReconcile => environmentReconcile.componentReconciles, {
    onDelete: 'CASCADE',
  })
  @JoinColumn({
      referencedColumnName: 'reconcile_id',
  })
  environmentReconcile: EnvironmentReconcile

  @Column()
  name: string

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
}
