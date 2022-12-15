import { BadRequestException, Controller, Get, Logger, NotFoundException, Query } from "@nestjs/common";
import { SystemService } from "./system.service";

@Controller({
  version: '1'
})
export class SystemController {
  private readonly logger = new Logger(SystemController.name);

  constructor(
      private readonly systemService: SystemService
  ){}

  // GET /ssmsecret?path=/some/secret
  @Get('/ssmsecret')
  public async getSsmSecret(@Query("path") path: string) {
    if (!path) {
      throw new BadRequestException('path is required');
    }

    const value = await this.systemService.getSsmSecret(path);

    if (!value) {
      throw new NotFoundException();
    }

    return { value };
  }
}
