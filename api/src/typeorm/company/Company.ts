import {
  Column,
  Entity,
  UpdateDateColumn,
} from "typeorm";

@Entity({ name: "company" })
export class Company {
  @Column({
    primary: true,
    name: "name",
  })
  name: string;

  @Column({
    name: "client_id",
  })
  clientId: string;

  @Column({
    name: "client_secret",
  })
  clientSecret: string;

  @Column({
    name: "github_repo",
  })
  githubRepo: string;

  @Column({
    name: "github_path",
  })
  githubPath: string;

  @Column({
    name: "github_source",
  })
  githubSource: string;


  @UpdateDateColumn({
    name: "timeStamp",
  })
  timeStamp: string;
}
