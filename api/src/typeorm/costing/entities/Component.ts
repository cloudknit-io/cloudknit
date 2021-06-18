import { Resource } from 'src/typeorm/resources/Resource.entity';
import { Column, Entity, JoinColumn, ManyToOne, OneToMany } from 'typeorm';

@Entity({ name: 'components' })
export class Component {
  @Column({
    primary: true,
    name: 'id',
  })
  id: string;

  @Column({
    name: 'team_name',
  })
  teamName: string;

  @Column({
    name: 'environment_name',
  })
  environmentName: string;

  @Column({
    name: 'component_name',
  })
  componentName: string;

  @Column({
    name: 'cost',
    type: 'decimal',
    precision: 10,
    scale: 3,
  })
  cost: number = 0;

  @OneToMany(type => Resource, resource => resource.component, {
    cascade: true,
  })
  resources: Resource[];
}
