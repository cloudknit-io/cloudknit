import { Organization, Team } from "src/typeorm";
import { EnvSpecComponentDto } from "./env-spec.dto";

export class CreateEnvironmentDto {
  name: string;
  team: Team;
  organization: Organization;
  dag: EnvSpecComponentDto[];
}
