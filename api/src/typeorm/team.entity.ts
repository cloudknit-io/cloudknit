import {
  Column,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  OneToMany,
  PrimaryGeneratedColumn,
  RelationId,
} from 'typeorm';
import { Organization } from './Organization.entity';
import { Environment } from './environment.entity';
import { ColumnNumericTransformer } from './helper';

@Entity({
  name: 'team',
})
@Index(['organization', 'name'], { unique: true })
export class Team {
  @PrimaryGeneratedColumn()
  id: number;

  @Column()
  name: string;

  @Column()
  repo: string;

  @Column()
  repo_path: string;

  @Column({
    default: false,
  })
  isDeleted: boolean;

  @Column({
    default: false,
  })
  teardownProtection: boolean;

  @Column({
    name: 'estimated_cost',
    type: 'decimal',
    precision: 10,
    scale: 3,
    default: 0,
    transformer: new ColumnNumericTransformer(),
  })
  estimatedCost: number;

  @OneToMany(() => Environment, (env) => env.team)
  environments: Environment[];

  @ManyToOne(() => Organization, (org) => org.id, {
    onDelete: 'CASCADE',
  })
  @JoinColumn({
    referencedColumnName: 'id',
  })
  organization: Organization;

  @RelationId((team: Team) => team.organization)
  orgId: number;
}
