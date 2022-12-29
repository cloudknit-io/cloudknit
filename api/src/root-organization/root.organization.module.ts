import { Module } from "@nestjs/common";
import { TypeOrmModule } from "@nestjs/typeorm";
import { Organization, User } from "src/typeorm";
import { UsersModule } from "src/users/users.module";
import { UsersService } from "src/users/users.service";
import { RootOrganizationsController } from "./root.organization.controller";
import { RootOrganizationsService } from "./root.organization.service";

@Module({
  imports: [TypeOrmModule.forFeature([Organization, User]), UsersModule],
  controllers: [RootOrganizationsController],
  providers: [RootOrganizationsService, UsersService],
})
export class RootOrganizationsModule {}
