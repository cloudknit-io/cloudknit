import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  OneToMany,
  PrimaryColumn,
} from 'typeorm'
import { Component } from '../costing/entities/Component'

@Entity({
  name: 'resources',
})
export class Resource {

  @PrimaryColumn()
  id: string;

  @Column()
  name: string

  @Column({
    nullable: true
  })
  hourlyCost?: string
  @Column({
    nullable: true
  })
  monthlyCost?: string

  @Column()
  componentId: string;

  @Column({
    nullable: true
  })
  parentId?: string;
  
  @OneToMany((type) => Resource, (resource) => resource.resource, {
    cascade: true,
  })
  subresources?: Resource[]

  @ManyToOne((type) => Resource, (resource) => resource.subresources)
  @JoinColumn({
    referencedColumnName: 'id'
  })
  resource?: Resource

  @ManyToOne((type) => Component, (component) => component.resources)
  @JoinColumn({
    referencedColumnName: 'id'
  })
  component?: Component


  @OneToMany(() => CostComponent, (component) => component.resource, {
    cascade: true,
    eager: true
  })
  costComponents?: CostComponent[]
}

@Entity({
  name: 'costcomponents',
})
export class CostComponent {
  @PrimaryColumn()
  id: string

  @Column({
    nullable: true
  })
  hourlyCost?: string
  @Column({
    nullable: true
  })
  hourlyQuantity?: string
  @Column({
    nullable: true
  })
  monthlyCost?: string
  @Column({
    nullable: true
  })
  monthlyQuantity?: string
  @Column({
    nullable: true
  })
  name?: string
  @Column({
    nullable: true
  })
  price?: string
  @Column({
    nullable: true
  })
  unit?: string

  @ManyToOne(() => Resource, resource => resource.resource)
  @JoinColumn({
    referencedColumnName: 'id'
  })
  resource?: Resource
}
