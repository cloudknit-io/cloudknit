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

  @UpdateDateColumn({
    name: "timeStamp",
  })
  timeStamp: string;
}
