import { MiddlewareConsumer, Module, NestModule, RequestMethod } from "@nestjs/common";
import { TypeOrmModule } from "@nestjs/typeorm";
import { OrganizationMiddleware } from "src/middleware/organizations.middle";
import { Organization } from "src/typeorm";
import { OrganizationController } from "./organization.controller";
import { OrganizationService } from "./organization.service";

@Module({
  imports: [TypeOrmModule.forFeature([Organization])],
  controllers: [OrganizationController],
  providers: [OrganizationService],
})
export class OrganizationModule implements NestModule {
  
  configure(consumer: MiddlewareConsumer) {
    consumer.apply(OrganizationMiddleware).forRoutes({
      /**
       * This is confusing but I don't know how else to do it.
       * 
       * The goal is for _all_ `/[ver]/orgs/:id` routes to have req.org which is handled
       * via OrgnizationMiddleware.
       * 
       * Since this path value is used to evaluate the path from root
       * we need the initial '*' to account for `/v1` or any other version.
       * 
       * This sheds some light on the issue:
       * https://github.com/nestjs/nest/issues/4210#issuecomment-595037635
       * 
       * The real issue here is if our global routing gets updated then it's likely
       * that this middleware will no longer work until it's updated to match the 
       * global changes. Nest is forcing us to make sure our routes are 
       * updated in two places.
       * 
       * - Brad
       */
      path: '/*/orgs/:orgId*',
      method: RequestMethod.ALL
    });
  }
}
