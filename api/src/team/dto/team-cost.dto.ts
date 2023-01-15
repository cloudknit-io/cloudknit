import { Team } from 'src/typeorm';

export class TeamWrapDto extends Team {
  estimatedCost?: number;
}
