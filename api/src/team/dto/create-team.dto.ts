import { Organization } from "src/typeorm";

export class CreateTeamDto {
  name: string;
  organization: Organization;
}
