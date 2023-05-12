import { ApiProperty } from '@nestjs/swagger';
import {
  Column,
  CreateDateColumn,
  Entity,
  ManyToMany,
  PrimaryGeneratedColumn,
  UpdateDateColumn,
} from 'typeorm';
import { Organization } from './Organization.entity';

@Entity({ name: 'users' })
export class User {
  @PrimaryGeneratedColumn()
  id: number;

  @ApiProperty()
  @Column({
    unique: true,
  })
  username: string;

  @ApiProperty()
  @Column({
    unique: true,
  })
  email: string;

  @ApiProperty()
  @Column({
    default: null,
  })
  name: string;

  @Column({
    default: 'User',
  })
  role: string;

  @Column({
    default: false,
  })
  archived: boolean;

  @Column({
    default: null,
    unique: true
  })
  ipv4: string;

  @ManyToMany(() => Organization, (org) => org.users)
  organizations: Organization[];

  @CreateDateColumn()
  created: Date;

  @UpdateDateColumn()
  updated: Date;
}
