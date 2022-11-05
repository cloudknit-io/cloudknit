import { Body, Controller, Delete, Get, Logger, NotFoundException, Param, Patch, Post, Request } from "@nestjs/common";
import { CreateUserDto, PatchUserDto } from "src/users/User.dto";
import { AuthService } from "./auth.service";

@Controller({
  version: '1'
})
export class AuthController {
  private readonly logger = new Logger(AuthController.name);

  constructor(private readonly authService: AuthService) {}

  @Get("users")
  public async getUsers(@Request() req) {
    return this.authService.getOrgUserList(req.org);
  }

  @Get("users/:username")
  public async getUser(@Request() req, @Param("username") username: string) {
    const user = await this.authService.getOrgUser(req.org, username);

    if (!user) {
      throw new NotFoundException('could not find user');
    }

    return user
  }

  @Post("users")
  public async createUser(@Request() req, @Body() user: CreateUserDto) {
    return await this.authService.createOrgUser(req.org, user);
  }

  // @Patch("users/:username")
  // public async updateUser(@Request() req, @Body() user: PatchUserDto) {
  //   return await this.authService.createOrgUser(req.org, user);
  // }

  @Delete("users/:username")
  public async deleteUser(@Request() req, @Param('username') username: string) {
    return await this.authService.deleteUser(req.org, username);
  }
}
