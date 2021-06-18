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
  @Column()
  hourlyCost?: string
  @Column()
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
  @Column()
  hourlyCost?: string
  @Column()
  hourlyQuantity?: string
  @Column()
  monthlyCost?: string
  @Column()
  monthlyQuantity?: string
  @Column()
  name?: string
  @Column()
  price?: string
  @Column()
  unit?: string

  @ManyToOne(() => Resource, resource => resource.resource)
  @JoinColumn({
    referencedColumnName: 'name'
  })
  resource?: Resource
}
