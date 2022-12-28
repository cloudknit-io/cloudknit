import { Organization, Team } from "src/typeorm";

export class CreateEnvironmentDto {
  name: string;
  team: Team;
  org: Organization;
}
