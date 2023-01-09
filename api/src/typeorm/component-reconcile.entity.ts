import { Column, Entity, JoinColumn, ManyToOne, PrimaryGeneratedColumn, RelationId } from 'typeorm'
import { Organization } from './Organization.entity'
import { EnvironmentReconcile } from './environment-reconcile.entity'

@Entity({
  name: 'component_reconcile',
})
export class ComponentReconcile {
  @PrimaryGeneratedColumn({
    name: 'id'
  })
  reconcileId: number

  @ManyToOne(() => EnvironmentReconcile, environmentReconcile => environmentReconcile.componentReconciles, {
    onDelete: 'CASCADE',
  })
  @JoinColumn({
      referencedColumnName: 'reconcileId',
  })
  environmentReconcile: EnvironmentReconcile

  @Column()
  name: string

  @Column()
  status: string

  @Column({
    nullable: true
  })
  approved_by?: string;

  @Column({
      type: 'datetime'
  })
  startDateTime: string;

  @Column({
      nullable: true
  })
  endDateTime?: string;

  @ManyToOne(() => Organization, (org) => org.id, {
    onDelete: "CASCADE",
  })
  @JoinColumn({
    referencedColumnName: 'id',
  })
  organization: Organization

  @RelationId((compRecon: ComponentReconcile) => compRecon.environmentReconcile)
  envReconId: number

  @RelationId((compRecon: ComponentReconcile) => compRecon.organization)
  orgId: number
}
