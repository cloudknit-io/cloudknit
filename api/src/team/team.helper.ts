import { Team } from "src/typeorm";
import { NoEnvironmentsError } from "src/types";

export function calculateTeamCost(team: Team): number {
  if (!team.environments) {
    throw new NoEnvironmentsError();
  }

  let total = 0.0;
      
  for (const env of team.environments) {
    total += env.estimatedCost;
  }

  return Math.round(total*100) / 100;
}
