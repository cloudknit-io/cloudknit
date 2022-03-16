import { ApiProperty } from "@nestjs/swagger";
import { Column, CreateDateColumn, Entity, PrimaryGeneratedColumn } from "typeorm";

@Entity({ name: "users" })
export class User {
  @PrimaryGeneratedColumn()
  id: number;

  @ApiProperty()
  @Column()
  username: string;

  @ApiProperty()
  @Column()
  email: string;

  @Column()
  company: string;

  @Column()
  termAgreementStatus: boolean;

  @CreateDateColumn()
  timeStamp: string;
}
