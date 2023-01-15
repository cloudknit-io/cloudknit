import { EnvSpecComponentDto } from './env-spec.dto';

export class CreateEnvironmentDto {
  name: string;
  dag: EnvSpecComponentDto[];
}
