import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  OneToMany,
  PrimaryColumn,
  PrimaryGeneratedColumn,
} from 'typeorm'
import { Component } from '../costing/entities/Component'

@Entity({
  name: 'resources',
})
export class Resource {
  @PrimaryColumn()
  name: string

  @Column({
    nullable: true
  })
  hourlyCost?: string
  @Column({
    nullable: true
  })
  monthlyCost?: string
  @OneToMany((type) => Resource, (resource) => resource.resource, {
    cascade: true,
  })
  subresources?: Resource[]

  @ManyToOne((type) => Resource, (resource) => resource.subresources)
  @JoinColumn({
    referencedColumnName: 'name'
  })
  resource?: Resource

  @ManyToOne((type) => Component, (component) => component.resources)
  @JoinColumn({
    referencedColumnName: 'id'
  })
  component?: Component


  @OneToMany(() => CostComponent, (component) => component.resource, {
    cascade: true,
  })
  costComponents?: CostComponent[]
}

@Entity({
  name: 'costcomponents',
})
export class CostComponent {
  @PrimaryGeneratedColumn()
  id: number

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
    referencedColumnName: 'name'
  })
  resource?: Resource
}
