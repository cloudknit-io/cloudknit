import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  PrimaryGeneratedColumn,
  RelationId,
} from 'typeorm';
import { Organization } from './Organization.entity';
import { EnvironmentReconcile } from './environment-reconcile.entity';
import { Component } from './component.entity';

@Entity({
  name: 'component_reconcile',
})
export class ComponentReconcile {
  @PrimaryGeneratedColumn({
    name: 'id',
  })
  reconcileId: number;

  @ManyToOne(
    () => EnvironmentReconcile,
    (environmentReconcile) => environmentReconcile.componentReconciles,
    {
      onDelete: 'CASCADE',
    }
  )
  @JoinColumn({
    referencedColumnName: 'reconcileId',
  })
  environmentReconcile: EnvironmentReconcile;

  @ManyToOne(() => Component, (component) => component.id, {
    eager: true,
  })
  component: Component;

  @Column()
  status: string;

  @Column({
    nullable: true,
  })
  approvedBy?: string;

  @Column({
    type: 'datetime',
  })
  startDateTime: string;

  @Column({
    nullable: true,
  })
  endDateTime?: string;

  @ManyToOne(() => Organization, (org) => org.id, {
    onDelete: 'CASCADE',
  })
  @JoinColumn({
    referencedColumnName: 'id',
  })
  organization: Organization;

  @RelationId((compRecon: ComponentReconcile) => compRecon.component)
  compId: number;

  @RelationId((compRecon: ComponentReconcile) => compRecon.environmentReconcile)
  envReconId: number;

  @RelationId((compRecon: ComponentReconcile) => compRecon.organization)
  orgId: number;
}
