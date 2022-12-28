import { Column, Entity, Index, JoinColumn, ManyToOne, OneToMany, PrimaryGeneratedColumn } from "typeorm";
import { Organization } from "./Organization.entity";
import { Environment } from "./environment.entity";

@Entity({
  name: "team",
})
export class Team {
  @PrimaryGeneratedColumn()
  id: number

  @Column()
  @Index()
  name: string;

  @Column({
    default: false
  })
  isDeleted: boolean

  @OneToMany(() => Environment, (env) => env.team)
  environments: Environment[];

  @ManyToOne(() => Organization, (org) => org.id, {
    onDelete: "CASCADE"
  })
  @JoinColumn({
    referencedColumnName: 'id'
  })
  organization: Organization
}
