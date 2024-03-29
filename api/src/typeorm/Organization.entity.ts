import {
  Column,
  CreateDateColumn,
  Entity,
  Index,
  JoinTable,
  ManyToMany,
  OneToOne,
  PrimaryGeneratedColumn,
  UpdateDateColumn,
} from 'typeorm';
import { User } from './User.entity';

@Entity({ name: 'organization' })
export class Organization {
  @PrimaryGeneratedColumn()
  id: number;

  @Index({ unique: true })
  @Column()
  name: string;

  @Column({
    name: 'github_repo',
    default: null,
  })
  githubRepo: string;

  @Column({
    name: 'github_org_name',
    default: null,
    unique: true,
  })
  githubOrgName: string;

  @OneToOne(() => User, (user) => user.id)
  @Column({
    default: null,
  })
  termsAgreedUserId: number;

  @ManyToMany(() => User, (user) => user.organizations)
  @JoinTable()
  users: User[];

  @Column({
    default: false,
  })
  provisioned: boolean;

  @UpdateDateColumn()
  updated: Date;

  @CreateDateColumn()
  created: Date;
}
