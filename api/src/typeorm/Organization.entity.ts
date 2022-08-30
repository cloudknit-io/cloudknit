import {
  Column,
  CreateDateColumn,
  Entity,
  Index,
  JoinTable,
  ManyToMany,
  PrimaryGeneratedColumn,
  UpdateDateColumn,
} from "typeorm";
import { User } from "./User.entity";

@Entity({ name: "organization" })
export class Organization {
  @PrimaryGeneratedColumn()
  id: number

  @Index({ unique: true })
  @Column()
  name: string

  @Column({
    name: "client_id",
    default: null
  })
  clientId: string;

  @Column({
    name: "client_secret",
    default: null
  })
  clientSecret: string;

  @Column({
    name: "github_repo",
    default: null
  })
  githubRepo: string;

  @Column({
    name: "github_path",
    default: null
  })
  githubPath: string;

  @Column({
    name: "github_source",
    default: null
  })
  githubSource: string;

  @ManyToMany(() => User, (user) => user.organizations)
  @JoinTable()
  users: User[]

  @UpdateDateColumn()
  updated: Date;

  @CreateDateColumn()
  created: Date;
}
