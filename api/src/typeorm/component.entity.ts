import { Column, Entity, Index, JoinColumn, ManyToOne, PrimaryGeneratedColumn, RelationId, UpdateDateColumn } from 'typeorm';
import { Organization } from './Organization.entity';
import { Environment } from './environment.entity';
import { CostResource } from 'src/component/dto/update-component.dto';
import { ColumnNumericTransformer } from './helper';

@Entity({ name: 'components' })
@Index(['organization', 'environment', 'name'], { unique: true })
export class Component {
  @PrimaryGeneratedColumn()
  id: number

  @ManyToOne(() => Environment, (environment) => environment.components)
  environment: Environment;

  @Column({
    name: 'component_name',
  })
  name: string;

  @Column({
    default: 'terraform'
  })
  type: string;

  @Column({
    name: 'status',
    default: null
  })
  status: string;

  @Column({
    name: 'estimated_cost',
    type: 'decimal',
    precision: 10,
    scale: 3,
    default: 0,
    transformer: new ColumnNumericTransformer()
  })
  estimatedCost: number;

  @UpdateDateColumn({
    name: 'last_reconcile_datetime'
  })
  lastReconcileDatetime: string;

  @Column({
    default: -1
  })
  duration: number;

  @Column({
    default: null,
    nullable: true
  })
  lastWorkflowRunId: string;

  @Column({
    default: false,
    type: 'boolean'
  })
  isDestroyed?: boolean

  @Column({
    name: 'cost_resources',
    default: null,
    type: 'json'
  })
  costResources: CostResource[];

  @ManyToOne(() => Organization, (org) => org.id, {
    onDelete: "CASCADE"
  })
  @JoinColumn({
    referencedColumnName: 'id'
  })
  organization: Organization

  @RelationId((comp: Component) => comp.environment)
  envId: number

  @RelationId((comp: Component) => comp.organization)
  orgId: number
}
