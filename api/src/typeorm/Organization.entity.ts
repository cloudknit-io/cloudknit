import {
  Column,
  CreateDateColumn,
  Entity,
  Index,
  JoinColumn,
  JoinTable,
  ManyToMany,
  OneToOne,
  PrimaryGeneratedColumn,
  UpdateDateColumn,
} from "typeorm";
import { User } from "./User.entity";

@Entity({ name: "organization" })
export class Organization {
  @PrimaryGeneratedColumn()
  id: number;

  @Index({ unique: true })
  @Column()
  name: string;

  @Column({
    name: "github_repo",
    default: null
  })
  githubRepo: string;

  @OneToOne(() => User, (user) => user.id)
  @Column({
    default: null
  })
  termsAgreedUserId: number;

  @ManyToMany(() => User, (user) => user.organizations)
  @JoinTable()
  users: User[];

  @UpdateDateColumn()
  updated: Date;

  @CreateDateColumn()
  created: Date;
}
