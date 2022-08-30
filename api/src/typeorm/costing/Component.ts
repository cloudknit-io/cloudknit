import { Resource } from 'src/typeorm/resources/Resource.entity';
import { Column, Entity, JoinColumn, ManyToOne, OneToMany } from 'typeorm';
import { Organization } from '../Organization.entity';
import { Environment } from '../reconciliation/environment.entity';

@Entity({ name: 'components' })
export class Component {
  // TODO : Get rid of this.
  @Column({
    primary: true,
    name: 'id',
  })
  id: string;

  @Column({
    name: 'team_name',
  })
  teamName: string;

  @ManyToOne(() => Environment, (environment) => environment.components, {
    eager: true
  })
  environment: Environment;

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

  @Column({
    default: false,
    type: 'boolean'
  })
  isDeleted?: boolean

  @OneToMany(() => Resource, resource => resource.component, {
    cascade: true,
  })
  resources: Resource[];

  @ManyToOne(() => Organization, (org) => org.id, {
    onDelete: "CASCADE"
  })
  @JoinColumn({
    referencedColumnName: 'id'
  })
  organization: Organization
}
