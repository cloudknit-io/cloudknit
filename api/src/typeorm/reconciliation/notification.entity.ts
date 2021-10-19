import {
  Column,
  Entity,
  PrimaryGeneratedColumn,
} from "typeorm";

@Entity({
  name: "notification",
})
export class Notification {
  @PrimaryGeneratedColumn()
  notification_id?: number;

  @Column()
  message: string;

  @Column()
  team_name: string;

  @Column()
  environment_name: string;

  @Column()
  company_id: string;

  @Column({
    type: "datetime",
  })
  timestamp: string;

  @Column({
    default: ''
  })
  message_type?: string;

  @Column({
      default: false
  })
  seen?: boolean;

  @Column({
    default: null,
    type: 'json',
    nullable: true
  })
  debug?: {}
}
