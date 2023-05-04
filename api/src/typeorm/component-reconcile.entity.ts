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
import { ColumnNumericTransformer } from './helper';
import { CostResource } from 'src/component/dto/update-component.dto';

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
    onDelete: 'CASCADE',
  })
  component: Component;

  @Column()
  status: string;

  @Column({
    name: 'estimated_cost',
    type: 'decimal',
    precision: 10,
    scale: 3,
    default: 0,
    transformer: new ColumnNumericTransformer(),
  })
  estimatedCost: number;

  @Column({
    default: null,
    nullable: true,
  })
  lastWorkflowRunId: string;

  @Column({
    default: false,
    type: 'boolean',
  })
  isDestroyed?: boolean;

  @Column({
    default: false,
    type: 'boolean',
  })
  isSkipped?: boolean;

  @Column({
    name: 'cost_resources',
    default: null,
    type: 'json',
  })
  costResources: CostResource[];

  @Column({
    nullable: true,
  })
  approvedBy?: string;

  @Column({
    type: 'datetime',
    default: null
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
