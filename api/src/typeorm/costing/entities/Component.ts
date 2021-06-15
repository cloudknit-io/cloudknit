import { Column, Entity, ManyToOne } from 'typeorm';

@Entity({ name: 'Components' })
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
}
